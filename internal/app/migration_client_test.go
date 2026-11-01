package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Only synthetic TempDir state is used by these tests. The fixture writes the
// previous wire formats explicitly so the new startup never has to import them.
func legacyClientMigrationFixture(t *testing.T) (string, AppEntry, string, string) {
	t.Helper()
	m, path := registryFixture(t)
	a := mustRegister(t, m, "Migration fixture")
	old := m.stores[a.ID]
	oldDef := Definition{Version: 1, Mappings: []Mapping{{Event: "screen.opened", Step: "Screen", Kind: "step", Source: "browser"}}}
	d, _ := json.Marshal(oldDef)
	if _, e := old.DB.Exec(`UPDATE definitions SET data=? WHERE version=1`, string(d)); e != nil {
		t.Fatal(e)
	}
	evidence := Evidence{Source: "browser", Signal: "logs", Name: "screen.opened", Scope: "synthetic-visit", Tab: "synthetic-tab", Record: "synthetic-record", Time: time.Now().UnixNano(), Receipt: time.Now().UnixNano(), Definition: 1, CaptureDefinition: 1, Attrs: map[string]string{}}
	raw, _ := json.Marshal(evidence)
	if _, e := old.DB.Exec(`INSERT INTO evidence(data,identity,at,scope) VALUES(?,?,?,?)`, string(raw), "browser/logs/synthetic-record", evidence.Time, evidence.Scope); e != nil {
		t.Fatal(e)
	}
	if _, e := old.DB.Exec(`INSERT INTO summaries VALUES(?,?)`, time.Now().AddDate(0, -2, 0).UTC().Format("2006-01-02"), `{"groups":{}}`); e != nil {
		t.Fatal(e)
	}
	if _, e := old.DB.Exec(`PRAGMA user_version=3`); e != nil {
		t.Fatal(e)
	}
	m.Close()
	server := strings.Repeat("f", 64)
	legacy := legacyAppRegistry{Version: 1, Revision: 8, Apps: []legacyAppEntry{{AppEntry: AppEntry{ID: a.ID, Generation: a.Generation, Name: a.Name, State: a.State, Storage: a.Storage}, Browser: a.Key, Server: server}}}
	b, _ := json.Marshal(legacy)
	if e := os.WriteFile(filepath.Join(filepath.Dir(path), "apps.json"), b, 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(filepath.Dir(path), ".apps-initialized"), []byte("1\n"), 0600); e != nil {
		t.Fatal(e)
	}
	return path, a, server, string(raw)
}
func TestClientMigrationPreflightPreservesStateAndNeedsBackup(t *testing.T) {
	path, a, _, _ := legacyClientMigrationFixture(t)
	root := filepath.Dir(path)
	registry, _ := os.ReadFile(filepath.Join(root, "apps.json"))
	store, _ := os.ReadFile(filepath.Join(root, a.Storage))
	report, e := MigrateClientRegistry(path, false, false)
	if e != nil || report.State != "ready" || report.Apps != 1 || report.Records != 1 || report.DefinitionsRequiringReview != 1 {
		t.Fatalf("preflight: %#v %v", report, e)
	}
	if _, e = MigrateClientRegistry(path, true, false); e == nil {
		t.Fatal("apply without independent backup confirmation")
	}
	afterRegistry, _ := os.ReadFile(filepath.Join(root, "apps.json"))
	afterStore, _ := os.ReadFile(filepath.Join(root, a.Storage))
	if !bytes.Equal(registry, afterRegistry) || !bytes.Equal(store, afterStore) {
		t.Fatal("preflight changed existing state")
	}
	if _, e = OpenRegistry(path); e == nil || !strings.Contains(e.Error(), "offline client migration") {
		t.Fatal("startup adopted old credentials", e)
	}
}
func TestClientMigrationRetainsOneKeyAndOriginalHistory(t *testing.T) {
	path, a, server, original := legacyClientMigrationFixture(t)
	root := filepath.Dir(path)
	report, e := MigrateClientRegistry(path, true, true)
	if e != nil || report.State != "complete" {
		t.Fatal(report, e)
	}
	m, e := OpenRegistry(path)
	if e != nil {
		t.Fatal(e)
	}
	defer m.Close()
	got, source, ok := m.credential("Bearer " + a.Key)
	if !ok || source != "client" || got.ID != a.ID || got.Generation != a.Generation {
		t.Fatal("client authority changed")
	}
	if _, _, ok = m.credential("Bearer " + server); ok {
		t.Fatal("retired key authorizes intake")
	}
	b, e := os.ReadFile(filepath.Join(root, "apps.json"))
	if e != nil || bytes.Contains(b, []byte(server)) || bytes.Contains(b, []byte("browserCredential")) || bytes.Contains(b, []byte("serverCredential")) {
		t.Fatal("registry retains retired roles")
	}
	var raw string
	if e = m.stores[a.ID].DB.QueryRow(`SELECT data FROM evidence WHERE identity='browser/logs/synthetic-record'`).Scan(&raw); e != nil || raw != original {
		t.Fatal("migration rewrote original evidence", e)
	}
	var count int
	if e = m.stores[a.ID].DB.QueryRow(`SELECT count(*) FROM summaries`).Scan(&count); e != nil || count != 1 {
		t.Fatal("migration removed summary", e)
	}
	defs, e := m.stores[a.ID].Definitions()
	if e != nil || len(defs) != 2 || len(defs[1].Mappings) != 0 {
		t.Fatal("prospective definition review boundary missing", e)
	}
	for _, p := range []string{filepath.Join(root, ".client-migration", "backup", "apps.json"), filepath.Join(root, ".client-migration", "backup", a.ID+".db")} {
		st, e := os.Stat(p)
		if e != nil || st.Mode().Perm()&0077 != 0 {
			t.Fatal("private recovery copy missing")
		}
	}
	if e = m.Update(a.ID, "", "server"); e == nil {
		t.Fatal("retired role rotation accepted")
	}
	if _, _, ok = m.credential("Bearer " + a.Key); !ok {
		t.Fatal("retired operation changed key")
	}
	if e = m.Update(a.ID, "", "key"); e != nil {
		t.Fatal(e)
	}
	if _, _, ok = m.credential("Bearer " + a.Key); ok {
		t.Fatal("rotation retained old key")
	}
}
func TestClientMigrationInterruptedPublicationResumesFailClosed(t *testing.T) {
	for _, phase := range []string{"staged", "authority-disabled", "store-published", "registry-published"} {
		t.Run(phase, func(t *testing.T) {
			path, a, server, _ := legacyClientMigrationFixture(t)
			_, e := migrateClientRegistry(path, true, true, func(s string) error {
				if s == phase {
					return errors.New("synthetic interruption")
				}
				return nil
			})
			if e == nil {
				t.Fatal("fault returned success")
			}
			if m, e := OpenRegistry(path); e == nil {
				m.Close()
				t.Fatal("partial migration started intake")
			}
			report, e := MigrateClientRegistry(path, false, false)
			if e != nil || report.State != "interrupted" {
				t.Fatal(report, e)
			}
			if _, e = MigrateClientRegistry(path, true, true); e != nil {
				t.Fatal("resume", e)
			}
			m, e := OpenRegistry(path)
			if e != nil {
				t.Fatal(e)
			}
			defer m.Close()
			if _, _, ok := m.credential("Bearer " + a.Key); !ok {
				t.Fatal("resume lost retained key")
			}
			if _, _, ok := m.credential("Bearer " + server); ok {
				t.Fatal("resume revived retired key")
			}
		})
	}
}
func TestClientMigrationRejectsTamperingAndLiveWriter(t *testing.T) {
	t.Run("live writer", func(t *testing.T) {
		m, path := registryFixture(t)
		mustRegister(t, m, "Active")
		if _, e := MigrateClientRegistry(path, false, false); e == nil {
			t.Fatal("preflight acquired live lease")
		}
	})
	t.Run("identity mismatch", func(t *testing.T) {
		path, a, _, _ := legacyClientMigrationFixture(t)
		db, e := migrationDB(filepath.Join(filepath.Dir(path), a.Storage))
		if e != nil {
			t.Fatal(e)
		}
		db.Close()
		// Replace the registry identity binding, retaining valid opaque syntax.
		p := filepath.Join(filepath.Dir(path), "apps.json")
		b, _ := os.ReadFile(p)
		b = bytes.ReplaceAll(b, []byte(a.Generation), []byte(strings.Repeat("a", 64)))
		os.WriteFile(p, b, 0600)
		if _, e = MigrateClientRegistry(path, false, false); e == nil {
			t.Fatal("mismatched store accepted")
		}
	})
	t.Run("newer deny marker", func(t *testing.T) {
		path, _, _, _ := legacyClientMigrationFixture(t)
		os.WriteFile(filepath.Join(filepath.Dir(path), ".apps-initialized"), []byte("2\n"), 0600)
		if _, e := MigrateClientRegistry(path, false, false); e == nil {
			t.Fatal("old registry accepted beside newer authority marker")
		}
	})
	t.Run("staged corruption", func(t *testing.T) {
		path, a, _, _ := legacyClientMigrationFixture(t)
		_, _ = migrateClientRegistry(path, true, true, func(s string) error {
			if s == "staged" {
				return errors.New("stop")
			}
			return nil
		})
		os.WriteFile(filepath.Join(filepath.Dir(path), ".client-migration", a.ID+".db"), []byte("corrupt"), 0600)
		if _, e := MigrateClientRegistry(path, true, true); e == nil {
			t.Fatal("corrupt stage published")
		}
		if m, e := OpenRegistry(path); e == nil {
			m.Close()
			t.Fatal("corrupt migration started")
		}
	})
}
