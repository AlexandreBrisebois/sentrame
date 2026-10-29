package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func registryFixture(t *testing.T) (*Registry, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "legacy.db")
	m, e := OpenRegistry(path)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(m.Close)
	return m, path
}
func mustRegister(t *testing.T, m *Registry, name string) AppEntry {
	t.Helper()
	a, e := m.Register(name)
	if e != nil {
		t.Fatal(e)
	}
	return a
}
func TestRegistryValidationPrivacyAndStores(t *testing.T) {
	m, path := registryFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Bêta")
	if a.ID == b.ID || a.Generation == b.Generation || a.Browser == a.Server || a.Server == b.Server {
		t.Fatal("authority collision")
	}
	defs, e := m.stores[a.ID].Definitions()
	if e != nil || len(defs) != 1 || defs[0].Goal != "" || len(defs[0].Mappings) != 0 || validateDefinition(defs[0]) != "" {
		t.Fatal("new app definition is not empty and valid")
	}
	for _, name := range []string{" ", " ALPHA ", strings.Repeat("a", 81), "bad\nname", "bad\tname"} {
		if _, e = m.Register(name); e == nil {
			t.Fatal("invalid name accepted")
		}
	}
	mustRegister(t, m, "Straße")
	if _, e = m.Register("STRASSE"); e == nil {
		t.Fatal("Unicode case folding duplicate accepted")
	}
	original := m.current
	for _, mutate := range []func(*appRegistry){
		func(v *appRegistry) { v.Version = 2 }, func(v *appRegistry) { v.Apps[1].Browser = a.Server }, func(v *appRegistry) { v.Apps[1].Name = "ALPHA" }, func(v *appRegistry) { v.Apps[1].ID = a.ID }, func(v *appRegistry) { v.Apps[0].Storage = "../escape" }, func(v *appRegistry) { v.Apps[0].Storage = "/absolute" }, func(v *appRegistry) { v.Apps[0].Storage = "apps/other/telemetry.db" },
	} {
		v := original
		v.Apps = append([]AppEntry{}, original.Apps...)
		mutate(&v)
		if m.validate(v) == nil {
			t.Fatal("invalid registry accepted")
		}
	}
	for _, p := range []string{m.root, filepath.Join(m.root, "apps.json"), filepath.Join(m.root, ".apps-initialized")} {
		st, e := os.Stat(p)
		if e != nil || st.Mode().Perm()&0077 != 0 {
			t.Fatal("private permissions missing")
		}
	}
	if _, e = OpenRegistry(path); e == nil {
		t.Fatal("second process lease accepted")
	}
	// Enumerate every SQL table and every cell, including schema-internal tables.
	for _, store := range m.stores {
		rows, e := store.DB.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
		if e != nil {
			t.Fatal(e)
		}
		names := []string{}
		for rows.Next() {
			var n string
			rows.Scan(&n)
			names = append(names, n)
		}
		rows.Close()
		for _, name := range names {
			r, e := store.DB.Query(`SELECT * FROM "` + name + `"`)
			if e != nil {
				t.Fatal(e)
			}
			cols, _ := r.Columns()
			for r.Next() {
				vals := make([]any, len(cols))
				ptrs := make([]any, len(cols))
				for i := range vals {
					ptrs[i] = &vals[i]
				}
				if e = r.Scan(ptrs...); e != nil {
					t.Fatal(e)
				}
				for _, val := range vals {
					value := fmt.Sprint(val)
					if raw, ok := val.([]byte); ok {
						value = string(raw)
					}
					for _, secret := range []string{a.Browser, a.Server, b.Browser, b.Server} {
						if strings.Contains(value, secret) {
							t.Fatal("credential copied into SQLite")
						}
					}
				}
			}
			r.Close()
		}
	}
	m.Close()
	reopened, e := OpenRegistry(path)
	if e != nil {
		t.Fatal(e)
	}
	defer reopened.Close()
	if reopened.current.Apps[0].Server != a.Server {
		t.Fatal("restart overwrote authority")
	}
}
func TestRegistryEntropyLimitsAndSymlinks(t *testing.T) {
	m, _ := registryFixture(t)
	m.entropy = bytes.NewReader(nil)
	if _, e := m.Register("No entropy"); e == nil || len(m.current.Apps) != 0 {
		t.Fatal("entropy failure activated app")
	}
	m.entropy = nil
	// Restore the production source using a separately initialized registry's source.
	other, _ := registryFixture(t)
	m.entropy = other.entropy
	a := mustRegister(t, m, "One")
	link := filepath.Join(m.root, "apps", a.ID, "linked")
	if e := os.Symlink(t.TempDir(), link); e != nil {
		t.Fatal(e)
	}
	if safePath(filepath.Join(link, "telemetry.db")) == nil {
		t.Fatal("symlink escape accepted")
	}
	v := m.current
	v.Apps = append([]AppEntry{}, v.Apps...)
	for len(v.Apps) <= MaxApps {
		v.Apps = append(v.Apps, a)
	}
	if m.validate(v) == nil {
		t.Fatal("limit accepted")
	}
}
func TestRegistryAtomicFailureAndRestart(t *testing.T) {
	for _, stage := range []string{"write", "file-sync", "rename", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			m, path := registryFixture(t)
			a := mustRegister(t, m, "Alpha")
			old, _ := os.ReadFile(filepath.Join(m.root, "apps.json"))
			m.fault = func(s string) error {
				if s == stage {
					return errors.New("fault")
				}
				return nil
			}
			if e := m.Update(a.ID, "", "server"); e == nil {
				t.Fatal("fault reported success")
			}
			_, _, authorized := m.credential("Bearer " + a.Server)
			if stage == "directory-sync" {
				if authorized || !m.uncertain {
					t.Fatal("uncertain commit kept authority")
				}
			} else {
				if !authorized {
					t.Fatal("precommit failure changed authority")
				}
				now, _ := os.ReadFile(filepath.Join(m.root, "apps.json"))
				if !bytes.Equal(old, now) {
					t.Fatal("precommit failure modified file")
				}
			}
			m.Close()
			next, e := OpenRegistry(path)
			if e != nil {
				t.Fatal(e)
			}
			defer next.Close()
			_, _, authorized = next.credential("Bearer " + a.Server)
			if authorized == (stage == "directory-sync") {
				t.Fatal("restart did not reconcile committed file")
			}
		})
	}
}
func TestRegistryConcurrentUpdatesAndRemovalRecovery(t *testing.T) {
	m, path := registryFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if e := m.Update(a.ID, fmt.Sprintf("Alpha %d", i), ""); e != nil {
				t.Error(e)
			}
		}(i)
	}
	wg.Wait()
	m.fault = func(s string) error {
		if s == "purge" {
			return errors.New("fault")
		}
		return nil
	}
	current, _ := m.entry(a.ID)
	if e := m.Delete(a.ID, current.Name); e == nil {
		t.Fatal("purge failure absent")
	}
	if _, _, ok := m.credential("Bearer " + a.Server); ok {
		t.Fatal("pending deletion active")
	}
	if _, _, ok := m.credential("Bearer " + b.Server); !ok {
		t.Fatal("other app disabled")
	}
	m.Close()
	next, e := OpenRegistry(path)
	if e != nil {
		t.Fatal(e)
	}
	defer next.Close()
	if len(next.current.Apps) != 1 {
		t.Fatal("deletion did not resume")
	}
	if e = next.Delete(b.ID, b.Name); e != nil {
		t.Fatal(e)
	}
	if len(next.current.Apps) != 0 || next.current.Apps == nil {
		t.Fatal("empty registry invalid")
	}
	c := mustRegister(t, next, "Alpha")
	if c.ID == a.ID || c.Generation == a.Generation || c.Server == a.Server {
		t.Fatal("removed authority reused")
	}
}
func TestRegistryOrphansAndMissingStateFailClosed(t *testing.T) {
	t.Run("orphan", func(t *testing.T) {
		m, path := registryFixture(t)
		m.fault = func(s string) error {
			if s == "rename" {
				return errors.New("fault")
			}
			return nil
		}
		if _, e := m.Register("Orphan"); e == nil {
			t.Fatal("failed commit accepted")
		}
		if len(m.current.Apps) != 0 || len(m.stores) != 0 {
			t.Fatal("orphan authorized or handle leaked")
		}
		m.Close()
		if _, e := OpenRegistry(path); e == nil {
			t.Fatal("orphan silently adopted")
		}
	})
	for _, file := range []string{"apps.json", ".apps-initialized"} {
		t.Run(file, func(t *testing.T) {
			m, path := registryFixture(t)
			m.Close()
			os.Remove(filepath.Join(filepath.Dir(path), file))
			if _, e := OpenRegistry(path); e == nil {
				t.Fatal("missing state accepted")
			}
		})
	}
}
func TestRegistryGreenfieldRejectsLegacyWithoutChangingFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	old, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	old.DB.Close()
	before, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = OpenRegistry(path); e == nil {
		t.Fatal("legacy data silently imported")
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("legacy data changed")
	}
}

func TestRegistrySupportedSchemaIdentityAndExternalEdits(t *testing.T) {
	for _, kind := range []string{"identity", "schema", "json", "version", "external"} {
		t.Run(kind, func(t *testing.T) {
			m, path := registryFixture(t)
			a := mustRegister(t, m, "Alpha")
			file := filepath.Join(m.root, "apps.json")
			switch kind {
			case "identity":
				m.stores[a.ID].DB.Exec(`UPDATE settings SET value=? WHERE key='generation'`, strings.Repeat("b", 64))
			case "schema":
				m.stores[a.ID].DB.Exec(`PRAGMA user_version=2`)
			case "json":
				os.WriteFile(file, []byte("{"), 0600)
			case "version":
				v := m.current
				v.Version = 2
				b, _ := json.Marshal(v)
				os.WriteFile(file, b, 0600)
			case "external":
				b, _ := os.ReadFile(file)
				os.WriteFile(file, append(b, ' '), 0600)
				if e := m.Update(a.ID, "Changed", ""); e == nil || !m.uncertain {
					t.Fatal("external edit was overwritten")
				}
			}
			m.Close()
			next, e := OpenRegistry(path)
			if kind == "external" {
				if e != nil {
					t.Fatal(e)
				}
				next.Close()
			} else if e == nil {
				next.Close()
				t.Fatal("invalid state accepted")
			}
		})
	}
}
func TestRegistryActualRegistrationLimitAndDuplicateEntropy(t *testing.T) {
	m, _ := registryFixture(t)
	for i := 0; i < MaxApps; i++ {
		mustRegister(t, m, fmt.Sprintf("App %d", i))
	}
	if _, e := m.Register("Too many"); e == nil || len(m.stores) != MaxApps {
		t.Fatal("actual registration limit missing")
	}
	other, _ := registryFixture(t)
	other.entropy = bytes.NewReader(bytes.Repeat([]byte{1}, 128))
	if _, e := other.Register("Duplicate credentials"); e == nil || len(other.current.Apps) != 0 {
		t.Fatal("duplicate generated tokens accepted")
	}
}

func TestRegistryEmptyRestartAndStoreInitializationFailure(t *testing.T) {
	m, path := registryFixture(t)
	a := mustRegister(t, m, "Alpha")
	if e := m.Delete(a.ID, a.Name); e != nil {
		t.Fatal(e)
	}
	m.Close()
	next, e := OpenRegistry(path)
	if e != nil {
		t.Fatal(e)
	}
	defer next.Close()
	if len(next.current.Apps) != 0 {
		t.Fatal("empty registry recreated an app")
	}
	next.fault = func(stage string) error {
		if stage == "store-sync" {
			return errors.New("fault")
		}
		return nil
	}
	if _, e := next.Register("Failed initialization"); e == nil || len(next.current.Apps) != 0 || len(next.stores) != 0 {
		t.Fatal("unready store activated")
	}
}

func TestRegistryNeverReusesAnOrphanDirectory(t *testing.T) {
	m, _ := registryFixture(t)
	id := strings.Repeat("a", 64)
	directory := filepath.Join(m.root, "apps", id)
	if e := os.MkdirAll(directory, 0700); e != nil {
		t.Fatal(e)
	}
	sentinel := filepath.Join(directory, "telemetry.db")
	if e := os.WriteFile(sentinel, []byte("orphan preserved"), 0600); e != nil {
		t.Fatal(e)
	}
	bytesOf := func(value byte) []byte { return bytes.Repeat([]byte{value}, 32) }
	entropy := append(bytesOf(0xaa), bytesOf(0xbb)...)
	entropy = append(entropy, bytesOf(0xcc)...)
	entropy = append(entropy, bytesOf(0xdd)...)
	m.entropy = bytes.NewReader(entropy)
	if _, e := m.Register("Do not adopt orphan"); e == nil {
		t.Fatal("orphan reused")
	}
	raw, _ := os.ReadFile(sentinel)
	if string(raw) != "orphan preserved" || len(m.current.Apps) != 0 || len(m.stores) != 0 {
		t.Fatal("orphan changed or activated")
	}
}
