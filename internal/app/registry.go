package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/text/cases"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"unicode"
	"unicode/utf8"
)

const MaxApps = 32

type AppEntry struct {
	ID         string `json:"id"`
	Generation string `json:"generation"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Storage    string `json:"storage"`
	Browser    string `json:"browserCredential,omitempty"`
	Server     string `json:"serverCredential,omitempty"`
}
type appRegistry struct {
	Version  int        `json:"version"`
	Revision uint64     `json:"revision"`
	Apps     []AppEntry `json:"apps"`
}

// Registry owns the authority lock. Intake holds a read lock through SQLite commit;
// rotation/removal holds the write lock through durable registry publication.
type Registry struct {
	storeLeases  map[string]*os.File
	mu           sync.RWMutex
	root         string
	current      appRegistry
	stores       map[string]*Store
	leases       []*os.File
	uncertain    bool
	entropy      io.Reader
	fault        func(string) error
	cleanupIndex int
	diskDigest   string
}

var opaqueID = regexp.MustCompile(`^[a-f0-9]{64}$`)

func checkedRandom(reader io.Reader) (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(reader, b); err != nil {
		return "", errors.New("secure randomness unavailable")
	}
	return hex.EncodeToString(b), nil
}
func registryError() error {
	return errors.New("app registry unavailable; restore a coordinated private backup and restart")
}
func (m *Registry) stage(name string) error {
	if m.fault != nil {
		return m.fault(name)
	}
	return nil
}
func privateLease(path string) (*os.File, error) {
	if safePath(path) != nil {
		return nil, registryError()
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, registryError()
	}
	if err = f.Chmod(0600); err == nil {
		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	}
	if err != nil {
		f.Close()
		return nil, errors.New("persistent state already owned or private lease unavailable")
	}
	return f, nil
}

// safePath rejects symlinks in every existing path component, including the root.
func safePath(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return registryError()
	}
	p := string(filepath.Separator)
	for _, part := range strings.Split(strings.TrimPrefix(abs, p), p) {
		p = filepath.Join(p, part)
		st, e := os.Lstat(p)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil || st.Mode()&os.ModeSymlink != 0 {
			return registryError()
		}
	}
	return nil
}
func (m *Registry) validate(v appRegistry) error {
	if v.Version != 1 || v.Apps == nil || len(v.Apps) > MaxApps {
		return registryError()
	}
	ids, names, tokens, paths := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, a := range v.Apps {
		name := strings.TrimSpace(a.Name)
		fold := cases.Fold().String(name)
		if !opaqueID.MatchString(a.ID) || !opaqueID.MatchString(a.Generation) || ids[a.ID] || name == "" || name != a.Name || !utf8.ValidString(name) || utf8.RuneCountInString(name) > 80 || names[fold] || (a.State != "active" && a.State != "deleting") {
			return registryError()
		}
		if strings.IndexFunc(name, unicode.IsControl) >= 0 {
			return registryError()
		}
		if a.Storage != filepath.Join("apps", a.ID, "telemetry.db") {
			return registryError()
		}
		if filepath.IsAbs(a.Storage) || filepath.Clean(a.Storage) != a.Storage || a.Storage == ".." || strings.HasPrefix(a.Storage, ".."+string(filepath.Separator)) || paths[a.Storage] {
			return registryError()
		}
		if err := safePath(filepath.Join(m.root, a.Storage)); err != nil {
			return err
		}
		for _, token := range []string{a.Browser, a.Server} {
			if !opaqueID.MatchString(token) || tokens[token] {
				return registryError()
			}
			tokens[token] = true
		}
		ids[a.ID] = true
		names[fold] = true
		paths[a.Storage] = true
	}
	return nil
}
func syncDirectory(path string) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	return f.Sync()
}
func (m *Registry) commit(v appRegistry) error {
	if m.uncertain {
		return registryError()
	}
	disk, readErr := os.ReadFile(filepath.Join(m.root, "apps.json"))
	if m.diskDigest != "" && (readErr != nil || digest(string(disk)) != m.diskDigest) || m.diskDigest == "" && !os.IsNotExist(readErr) {
		m.uncertain = true
		return registryError()
	}
	if err := m.validate(v); err != nil {
		return err
	}
	v.Revision = m.current.Revision + 1
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return registryError()
	}
	if err = m.stage("write"); err != nil {
		return registryError()
	}
	f, err := os.CreateTemp(m.root, ".apps-*")
	if err != nil {
		return registryError()
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if err = f.Chmod(0600); err == nil {
		_, err = f.Write(b)
	}
	if err == nil {
		err = m.stage("file-sync")
	}
	if err == nil {
		err = f.Sync()
	}
	if err == nil {
		err = f.Close()
	}
	if err == nil {
		err = m.stage("rename")
	}
	if err == nil {
		err = os.Rename(f.Name(), filepath.Join(m.root, "apps.json"))
	}
	if err != nil {
		return registryError()
	}
	// The replacement may be durable even if directory synchronization fails. Do
	// not publish speculative credentials or keep trusting the previous ones.
	if err = m.stage("directory-sync"); err == nil {
		err = syncDirectory(m.root)
	}
	if err != nil {
		m.uncertain = true
		return registryError()
	}
	m.current = v
	m.diskDigest = digest(string(b))
	return nil
}
func (m *Registry) openStore(a AppEntry, fresh bool) error {
	path := filepath.Join(m.root, a.Storage)
	for _, suffix := range []string{"", "-wal", "-shm"} {
		if err := safePath(path + suffix); err != nil {
			return err
		}
	}
	if fresh {
		if err := os.MkdirAll(filepath.Dir(filepath.Dir(path)), 0700); err != nil {
			return registryError()
		}
		if err := os.Mkdir(filepath.Dir(path), 0700); err != nil {
			return registryError()
		}
	} else if _, err := os.Stat(path); err != nil {
		return registryError()
	}
	lease, err := privateLease(path + ".lock")
	if err != nil {
		return err
	}
	s, err := Open(path)
	if err != nil {
		lease.Close()
		return err
	}
	ok := false
	defer func() {
		if !ok {
			s.DB.Close()
			lease.Close()
		}
	}()
	if fresh {
		d := Definition{Version: 1, Mappings: []Mapping{}}
		b, _ := json.Marshal(d)
		tx, e := s.DB.Begin()
		if e != nil {
			return registryError()
		}
		defer tx.Rollback()
		if _, e = tx.Exec(`UPDATE settings SET value=? WHERE key='app_id';`, a.ID); e != nil {
			return registryError()
		}
		if _, e = tx.Exec(`UPDATE settings SET value=? WHERE key='generation'`, a.Generation); e != nil {
			return registryError()
		}
		if _, e = tx.Exec(`UPDATE definitions SET data=? WHERE version=1`, string(b)); e != nil {
			return registryError()
		}
		if e = tx.Commit(); e != nil {
			return registryError()
		}
	}
	var id, generation string
	if s.DB.QueryRow(`SELECT value FROM settings WHERE key='app_id'`).Scan(&id) != nil || s.DB.QueryRow(`SELECT value FROM settings WHERE key='generation'`).Scan(&generation) != nil || id != a.ID || generation != a.Generation {
		return registryError()
	}
	if err = s.Cleanup(); err != nil {
		return registryError()
	}
	if fresh {
		if err = m.stage("store-sync"); err != nil {
			return registryError()
		}
		if _, err = s.DB.Exec(`PRAGMA wal_checkpoint(FULL)`); err != nil {
			return registryError()
		}
		for _, directory := range []string{filepath.Dir(path), filepath.Dir(filepath.Dir(path)), m.root} {
			if err = syncDirectory(directory); err != nil {
				return registryError()
			}
		}
	}
	m.stores[a.ID] = s
	m.storeLeases[a.ID] = lease
	ok = true
	return nil
}

// OpenRegistry initializes a greenfield registry; legacy databases are never imported.
func OpenRegistry(dbPath string) (*Registry, error) {
	abs, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, registryError()
	}
	root := filepath.Dir(abs)
	// Canonicalize system ancestors (macOS /var is a symlink), while refusing a
	// symlink at the configured data root or anywhere inside the private tree.
	parent, e := filepath.EvalSymlinks(filepath.Dir(root))
	if e != nil {
		return nil, registryError()
	}
	root = filepath.Join(parent, filepath.Base(root))
	abs = filepath.Join(root, filepath.Base(abs))
	m := &Registry{root: root, stores: map[string]*Store{}, storeLeases: map[string]*os.File{}, entropy: rand.Reader}
	if err = safePath(abs); err != nil {
		return nil, err
	}
	if err = os.MkdirAll(m.root, 0700); err != nil {
		return nil, registryError()
	}
	if err = os.Chmod(m.root, 0700); err != nil {
		return nil, registryError()
	}
	lease, err := privateLease(filepath.Join(m.root, ".apps.lock"))
	if err != nil {
		return nil, err
	}
	m.leases = append(m.leases, lease)
	success := false
	defer func() {
		if !success {
			m.Close()
		}
	}()
	registryPath := filepath.Join(m.root, "apps.json")
	marker := filepath.Join(m.root, ".apps-initialized")
	if safePath(registryPath) != nil || safePath(marker) != nil {
		return nil, registryError()
	}
	b, err := os.ReadFile(registryPath)
	if err == nil {
		if len(b) > 1<<20 || json.Unmarshal(b, &m.current) != nil || m.validate(m.current) != nil {
			return nil, registryError()
		}
		if _, e := os.Stat(marker); e != nil {
			return nil, registryError()
		}
		m.diskDigest = digest(string(b))
		if os.Chmod(registryPath, 0600) != nil {
			return nil, registryError()
		}
	} else if os.IsNotExist(err) {
		if _, e := os.Stat(marker); !os.IsNotExist(e) {
			return nil, registryError()
		}
		if _, e := os.Stat(filepath.Join(m.root, "apps")); !os.IsNotExist(e) {
			return nil, registryError()
		}
		v := appRegistry{Version: 1, Apps: []AppEntry{}}
		if _, e := os.Stat(abs); e == nil {
			return nil, errors.New("legacy single-app database found; select a fresh data directory and register apps in the owner UI (existing files were preserved)")
		} else if !os.IsNotExist(e) {
			return nil, registryError()
		}
		if m.validate(v) != nil {
			return nil, registryError()
		}
		f, e := os.OpenFile(marker, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return nil, registryError()
		}
		_, e = f.WriteString("1\n")
		if e == nil {
			e = f.Sync()
		}
		f.Close()
		if e == nil {
			e = syncDirectory(m.root)
		}
		if e != nil {
			return nil, registryError()
		}
		if err = m.commit(v); err != nil {
			return nil, err
		}
	} else {
		return nil, registryError()
	}
	for _, a := range m.current.Apps {
		if a.State == "active" {
			if err = m.openStore(a, false); err != nil {
				return nil, err
			}
		}
	}
	if err = m.resumeDeleting(); err != nil {
		return nil, err
	}
	// Orphans never become authority. Require explicit offline cleanup instead of
	// silently hiding an interrupted registration.
	dirs, e := os.ReadDir(filepath.Join(m.root, "apps"))
	if e != nil && !os.IsNotExist(e) {
		return nil, registryError()
	}
	for _, d := range dirs {
		found := false
		for _, a := range m.current.Apps {
			if a.Storage == filepath.Join("apps", d.Name(), "telemetry.db") {
				found = true
			}
		}
		if !found {
			return nil, errors.New("interrupted app registration; quarantine unregistered app directory offline and restart")
		}
	}
	success = true
	return m, nil
}
func (m *Registry) Close() {
	for _, s := range m.stores {
		s.DB.Close()
	}
	for _, l := range m.storeLeases {
		l.Close()
	}
	for _, l := range m.leases {
		l.Close()
	}
}
func (m *Registry) entry(id string) (AppEntry, bool) {
	if m.uncertain {
		return AppEntry{}, false
	}
	for _, a := range m.current.Apps {
		if a.ID == id && a.State == "active" {
			return a, true
		}
	}
	return AppEntry{}, false
}
func (m *Registry) credential(h string) (AppEntry, string, bool) {
	if m.uncertain {
		return AppEntry{}, "", false
	}
	for _, a := range m.current.Apps {
		if a.State == "active" {
			if h == "Bearer "+a.Browser {
				return a, "browser", true
			}
			if h == "Bearer "+a.Server {
				return a, "server", true
			}
		}
	}
	return AppEntry{}, "", false
}
func (m *Registry) Register(name string) (AppEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.uncertain {
		return AppEntry{}, registryError()
	}
	a := AppEntry{Name: strings.TrimSpace(name), State: "active"}
	var e error
	for _, p := range []*string{&a.ID, &a.Generation, &a.Browser, &a.Server} {
		*p, e = checkedRandom(m.entropy)
		if e != nil {
			return AppEntry{}, e
		}
	}
	a.Storage = filepath.Join("apps", a.ID, "telemetry.db")
	v := m.current
	v.Apps = append(append([]AppEntry{}, v.Apps...), a)
	if e = m.validate(v); e != nil {
		return AppEntry{}, errors.New("use a unique app name of 1–80 characters; maximum 32 apps")
	}
	if e = m.openStore(a, true); e != nil {
		return AppEntry{}, e
	}
	if e = m.commit(v); e != nil {
		m.stores[a.ID].DB.Close()
		delete(m.stores, a.ID)
		m.storeLeases[a.ID].Close()
		delete(m.storeLeases, a.ID)
		return AppEntry{}, e
	}
	return a, nil
}
func (m *Registry) Update(id, name, role string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.current
	v.Apps = append([]AppEntry{}, v.Apps...)
	found := false
	for i, a := range v.Apps {
		if a.ID != id || a.State != "active" {
			continue
		}
		found = true
		if role != "" {
			token, e := checkedRandom(m.entropy)
			if e != nil {
				return e
			}
			switch role {
			case "browser":
				a.Browser = token
			case "server":
				a.Server = token
			default:
				return errors.New("choose browser or server")
			}
		} else {
			a.Name = strings.TrimSpace(name)
		}
		v.Apps[i] = a
	}
	if !found {
		return errors.New("app unavailable")
	}
	if m.validate(v) != nil && role == "" {
		return errors.New("use a unique app name of 1–80 characters; maximum 32 apps")
	}
	return m.commit(v)
}
func (m *Registry) Delete(id, confirmation string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var a AppEntry
	ok := false
	for _, candidate := range m.current.Apps {
		if candidate.ID == id {
			a = candidate
			ok = true
		}
	}
	if m.uncertain {
		return registryError()
	}
	if !ok || confirmation != a.Name {
		return errors.New("confirm the exact app name")
	}
	if a.State == "deleting" {
		return m.resumeDeleting()
	}
	v := m.current
	v.Apps = append([]AppEntry{}, v.Apps...)
	for i := range v.Apps {
		if v.Apps[i].ID == id {
			v.Apps[i].State = "deleting"
		}
	}
	if e := m.commit(v); e != nil {
		return e
	}
	return m.resumeDeleting()
}
func (m *Registry) resumeDeleting() error {
	for _, a := range m.current.Apps {
		if a.State != "deleting" {
			continue
		}
		if s := m.stores[a.ID]; s != nil {
			s.DB.Close()
			delete(m.stores, a.ID)
		}
		if lease := m.storeLeases[a.ID]; lease != nil {
			lease.Close()
			delete(m.storeLeases, a.ID)
		}
		path := filepath.Join(m.root, a.Storage)
		if safePath(path) != nil {
			return registryError()
		}
		if e := m.stage("purge"); e != nil {
			return registryError()
		}
		for _, suffix := range []string{"", "-wal", "-shm"} {
			if safePath(path+suffix) != nil {
				return registryError()
			}
			if e := os.Remove(path + suffix); e != nil && !os.IsNotExist(e) {
				return registryError()
			}
		}
		{
			if e := os.Remove(path + ".lock"); e != nil && !os.IsNotExist(e) {
				return registryError()
			}
			if e := os.Remove(filepath.Dir(path)); e != nil && !os.IsNotExist(e) {
				return registryError()
			}
			if e := syncDirectory(filepath.Dir(filepath.Dir(path))); e != nil {
				return registryError()
			}
		}
		if e := syncDirectory(m.root); e != nil {
			return registryError()
		}
		v := m.current
		v.Apps = []AppEntry{}
		for _, other := range m.current.Apps {
			if other.ID != a.ID {
				v.Apps = append(v.Apps, other)
			}
		}
		if e := m.commit(v); e != nil {
			return e
		}
	}
	return nil
}
func (m *Registry) Cleanup() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.uncertain {
		return registryError()
	}
	n := len(m.current.Apps)
	if n == 0 {
		return nil
	}
	m.cleanupIndex = (m.cleanupIndex + 1) % n
	a := m.current.Apps[m.cleanupIndex]
	if s := m.stores[a.ID]; a.State == "active" && s != nil {
		return s.Cleanup()
	}
	return nil
}
func (a AppEntry) public() map[string]string {
	return map[string]string{"id": a.ID, "generation": a.Generation, "name": a.Name, "state": a.State}
}
