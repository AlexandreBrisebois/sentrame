package app

// This converter is deliberately separate from startup. It never runs merely
// because a new image finds old state, and never generates replacement keys.
import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type ClientMigrationReport struct {
	Apps                       int    `json:"apps"`
	Records                    int    `json:"records"`
	DefinitionsRequiringReview int    `json:"definitionsRequiringReview"`
	Cutover                    string `json:"cutover,omitempty"`
	State                      string `json:"state"`
}
type legacyAppEntry struct {
	AppEntry
	Browser string `json:"browserCredential"`
	Server  string `json:"serverCredential"`
}
type legacyAppRegistry struct {
	Version  int              `json:"version"`
	Revision uint64           `json:"revision"`
	Apps     []legacyAppEntry `json:"apps"`
}
type clientMigration struct {
	RegistryDigest string                `json:"registryDigest"`
	TargetDigest   string                `json:"targetDigest"`
	Cutover        int64                 `json:"cutover"`
	Stores         map[string]string     `json:"stores"`
	Report         ClientMigrationReport `json:"report"`
}

var errClientMigration = errors.New("client migration unavailable; preserve coordinated state, keep intake stopped and inspect the offline migration files")

// MigrateClientRegistry preflights an existing managed v1 registry. Applying is
// explicit and requires acknowledgement of an independent coordinated backup.
// Its report contains counts only, never keys. All filesystem work is confined
// to that data root; interrupted publication resumes from verified staged files.
func MigrateClientRegistry(dbPath string, apply, backupConfirmed bool) (ClientMigrationReport, error) {
	return migrateClientRegistry(dbPath, apply, backupConfirmed, nil)
}
func migrateClientRegistry(dbPath string, apply, backupConfirmed bool, fault func(string) error) (ClientMigrationReport, error) {
	report := ClientMigrationReport{State: "preflight"}
	abs, err := filepath.Abs(dbPath)
	if err != nil {
		return report, errClientMigration
	}
	root := filepath.Dir(abs)
	parent, err := filepath.EvalSymlinks(filepath.Dir(root))
	if err != nil {
		return report, errClientMigration
	}
	root = filepath.Join(parent, filepath.Base(root))
	if safePath(root) != nil {
		return report, errClientMigration
	}
	if st, e := os.Stat(root); e != nil || !st.IsDir() {
		return report, errClientMigration
	}
	lease, err := privateLease(filepath.Join(root, ".apps.lock"))
	if err != nil {
		return report, err
	}
	defer lease.Close()
	stage := filepath.Join(root, ".client-migration")
	marker := filepath.Join(root, ".client-migration.json")
	for _, p := range []string{stage, marker, filepath.Join(root, "apps.json"), filepath.Join(root, ".apps-initialized")} {
		if safePath(p) != nil {
			return report, errClientMigration
		}
	}
	step := func(name string) error {
		if fault != nil {
			return fault(name)
		}
		return nil
	}
	var state clientMigration
	markerRaw, markerErr := os.ReadFile(marker)
	if markerErr == nil {
		if len(markerRaw) > 1<<20 || json.Unmarshal(markerRaw, &state) != nil || state.Cutover <= 0 || state.Stores == nil {
			return report, errClientMigration
		}
		report = state.Report
		report.State = "interrupted"
		if !apply {
			return report, nil
		}
	} else if !os.IsNotExist(markerErr) {
		return report, errClientMigration
	}
	if apply && !backupConfirmed {
		return report, errors.New("applying requires confirmation of a current independent private coordinated backup")
	}
	raw, err := os.ReadFile(filepath.Join(root, "apps.json"))
	if err != nil || len(raw) > 1<<20 {
		return report, errClientMigration
	}
	var legacy legacyAppRegistry
	if markerErr != nil {
		markerVersion, e := os.ReadFile(filepath.Join(root, ".apps-initialized"))
		if e != nil || string(markerVersion) != "1\n" {
			return report, errClientMigration
		}
		if json.Unmarshal(raw, &legacy) != nil || legacy.Version != 1 || legacy.Apps == nil || len(legacy.Apps) > MaxApps {
			return report, errClientMigration
		}
	} else {
		backup, err := os.ReadFile(filepath.Join(stage, "backup", "apps.json"))
		if err != nil || digest(string(backup)) != state.RegistryDigest || json.Unmarshal(backup, &legacy) != nil {
			return report, errClientMigration
		}
		if digest(string(raw)) != state.RegistryDigest && digest(string(raw)) != state.TargetDigest && string(raw) != "{\"version\":99,\"apps\":[]}" {
			return report, errClientMigration
		}
	}
	target := appRegistry{Version: 2, Revision: legacy.Revision + 1, Apps: []AppEntry{}}
	tokens := map[string]bool{}
	for _, a := range legacy.Apps {
		// Incomplete deletion is never silently completed by a credential migration.
		if a.State != "active" || !opaqueID.MatchString(a.Browser) || !opaqueID.MatchString(a.Server) || a.Browser == a.Server || tokens[a.Browser] || tokens[a.Server] {
			return report, errClientMigration
		}
		tokens[a.Browser] = true
		tokens[a.Server] = true
		a.AppEntry.Key = a.Browser
		target.Apps = append(target.Apps, a.AppEntry)
	}
	validator := &Registry{root: root}
	if validator.validate(target) != nil {
		return report, errClientMigration
	}
	targetRaw, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return report, errClientMigration
	}
	leases := []*os.File{}
	defer func() {
		for _, f := range leases {
			f.Close()
		}
	}()
	for _, a := range target.Apps {
		path := filepath.Join(root, a.Storage)
		for _, suffix := range []string{"", "-wal", "-shm", ".lock"} {
			if safePath(path+suffix) != nil {
				return report, errClientMigration
			}
		}
		l, e := privateLease(path + ".lock")
		if e != nil {
			return report, e
		}
		leases = append(leases, l)
	}
	if markerErr != nil {
		report.Apps = len(target.Apps)
		var bytes int64
		for _, a := range target.Apps {
			path := filepath.Join(root, a.Storage)
			db, e := migrationDB(path)
			if e != nil {
				return report, errClientMigration
			}
			count, review, e := inspectMigrationStore(db, a)
			db.Close()
			if e != nil {
				return report, e
			}
			report.Records += count
			if review {
				report.DefinitionsRequiringReview++
			}
			for _, suffix := range []string{"", "-wal"} {
				if st, e := os.Stat(path + suffix); e == nil {
					bytes += st.Size()
				} else if !os.IsNotExist(e) {
					return report, errClientMigration
				}
			}
		}
		report.State = "ready"
		if !apply {
			return report, nil
		}
		var fs syscall.Statfs_t
		if syscall.Statfs(root, &fs) != nil || uint64(fs.Bavail)*uint64(fs.Bsize) < uint64(bytes)*2+(1<<20) {
			return report, errors.New("insufficient space for coordinated migration staging; state preserved")
		}
		if _, e := os.Stat(stage); !os.IsNotExist(e) {
			return report, errors.New("existing migration staging requires offline inspection; no state changed")
		}
		if os.Mkdir(stage, 0700) != nil || os.Mkdir(filepath.Join(stage, "backup"), 0700) != nil {
			return report, errClientMigration
		}
		// All copies are complete before the marker or authority publication.
		if writeMigrationFile(filepath.Join(stage, "backup", "apps.json"), raw) != nil {
			return report, errClientMigration
		}
		init, e := os.ReadFile(filepath.Join(root, ".apps-initialized"))
		if e != nil {
			return report, errClientMigration
		}
		if writeMigrationFile(filepath.Join(stage, "backup", ".apps-initialized"), init) != nil {
			return report, errClientMigration
		}
		state = clientMigration{RegistryDigest: digest(string(raw)), TargetDigest: digest(string(targetRaw)), Cutover: time.Now().UnixNano(), Stores: map[string]string{}, Report: report}
		for _, a := range target.Apps {
			path := filepath.Join(root, a.Storage)
			db, e := migrationDB(path)
			if e != nil {
				return report, errClientMigration
			}
			backup := filepath.Join(stage, "backup", a.ID+".db")
			_, e = db.Exec(`VACUUM INTO ?`, backup)
			db.Close()
			if e != nil {
				return report, errClientMigration
			}
			if os.Chmod(backup, 0600) != nil {
				return report, errClientMigration
			}
			staged := filepath.Join(stage, a.ID+".db")
			if copyMigrationFile(backup, staged) != nil {
				return report, errClientMigration
			}
			if e = prepareClientStore(staged, state.Cutover); e != nil {
				return report, e
			}
			hash, e := migrationFileDigest(staged)
			if e != nil {
				return report, errClientMigration
			}
			state.Stores[a.ID] = hash
		}
		state.Report.Cutover = time.Unix(0, state.Cutover).UTC().Format(time.RFC3339Nano)
		b, _ := json.Marshal(state)
		if writeMigrationFile(marker, b) != nil {
			return report, errClientMigration
		}
		if e := step("staged"); e != nil {
			return report, e
		}
	}
	if state.TargetDigest != digest(string(targetRaw)) || len(state.Stores) != len(target.Apps) {
		return report, errClientMigration
	}
	// A deliberately unsupported registry blocks old binaries during publication.
	if writeMigrationFile(filepath.Join(root, "apps.json"), []byte(`{"version":99,"apps":[]}`)) != nil {
		return report, errClientMigration
	}
	if e := step("authority-disabled"); e != nil {
		return report, e
	}
	for _, a := range target.Apps {
		staged := filepath.Join(stage, a.ID+".db")
		hash, e := migrationFileDigest(staged)
		if e != nil || hash != state.Stores[a.ID] {
			return report, errClientMigration
		}
		path := filepath.Join(root, a.Storage)
		if copyMigrationFile(staged, path) != nil {
			return report, errClientMigration
		}
		for _, suffix := range []string{"-wal", "-shm"} {
			if e := os.Remove(path + suffix); e != nil && !os.IsNotExist(e) {
				return report, errClientMigration
			}
		}
		if syncDirectory(filepath.Dir(path)) != nil {
			return report, errClientMigration
		}
		if e := step("store-published"); e != nil {
			return report, e
		}
	}
	if writeMigrationFile(filepath.Join(root, "apps.json"), targetRaw) != nil {
		return report, errClientMigration
	}
	if e := step("registry-published"); e != nil {
		return report, e
	}
	if writeMigrationFile(filepath.Join(root, ".apps-initialized"), []byte("2\n")) != nil {
		return report, errClientMigration
	}
	if os.Remove(marker) != nil || syncDirectory(root) != nil {
		return report, errClientMigration
	}
	report = state.Report
	report.State = "complete"
	return report, nil
}
func migrationDB(path string) (*sql.DB, error) {
	u := url.URL{Scheme: "file", Path: path}
	q := u.Query()
	q.Set("mode", "ro")
	u.RawQuery = q.Encode()
	db, e := sql.Open("sqlite", u.String())
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(1)
	return db, nil
}
func inspectMigrationStore(db *sql.DB, a AppEntry) (int, bool, error) {
	var version, count int
	var id, generation string
	if db.QueryRow(`PRAGMA user_version`).Scan(&version) != nil || version != 3 || db.QueryRow(`SELECT value FROM settings WHERE key='app_id'`).Scan(&id) != nil || db.QueryRow(`SELECT value FROM settings WHERE key='generation'`).Scan(&generation) != nil || id != a.ID || generation != a.Generation || db.QueryRow(`SELECT count(*) FROM evidence`).Scan(&count) != nil || count > 100000 {
		return 0, false, errClientMigration
	}
	var raw string
	if db.QueryRow(`SELECT data FROM definitions ORDER BY version DESC LIMIT 1`).Scan(&raw) != nil {
		return 0, false, errClientMigration
	}
	var d Definition
	if json.Unmarshal([]byte(raw), &d) != nil {
		return 0, false, errClientMigration
	}
	return count, legacyDefinitionNeedsReview(d), nil
}
func legacyDefinitionNeedsReview(d Definition) bool {
	for _, m := range d.Mappings {
		if m.Source != "" {
			return true
		}
	}
	return d.GoalSelector != nil && d.GoalSelector.Source != ""
}
func prepareClientStore(path string, cutover int64) error {
	db, e := sql.Open("sqlite", path)
	if e != nil {
		return errClientMigration
	}
	db.SetMaxOpenConns(1)
	defer db.Close()
	if _, e = db.Exec(`PRAGMA journal_mode=DELETE; PRAGMA synchronous=FULL`); e != nil {
		return errClientMigration
	}
	tx, e := db.Begin()
	if e != nil {
		return errClientMigration
	}
	defer tx.Rollback()
	if _, e = tx.Exec(`INSERT OR REPLACE INTO settings VALUES('client_cutover',?)`, fmt.Sprint(cutover)); e != nil {
		return errClientMigration
	}

	if _, e = tx.Exec(`PRAGMA user_version=4`); e != nil {
		return errClientMigration
	}
	if e = tx.Commit(); e != nil {
		return errClientMigration
	}
	if e = stageClientHistory(db); e != nil {
		return e
	}
	if e = db.Close(); e != nil {
		return errClientMigration
	}
	f, e := os.OpenFile(path, os.O_RDWR, 0600)
	if e != nil {
		return errClientMigration
	}
	defer f.Close()
	if f.Sync() != nil {
		return errClientMigration
	}
	return nil
}
func writeMigrationFile(path string, b []byte) error {
	if safePath(path) != nil {
		return errClientMigration
	}
	f, e := os.CreateTemp(filepath.Dir(path), ".migration-*")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if e = f.Chmod(0600); e == nil {
		_, e = f.Write(b)
	}
	if e == nil {
		e = f.Sync()
	}
	if e == nil {
		e = f.Close()
	}
	if e == nil {
		e = os.Rename(f.Name(), path)
	}
	if e == nil {
		e = syncDirectory(filepath.Dir(path))
	}
	return e
}
func migrationFileDigest(path string) (string, error) {
	if safePath(path) != nil {
		return "", errClientMigration
	}
	f, e := os.Open(path)
	if e != nil {
		return "", e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() > 1<<30 {
		return "", errClientMigration
	}
	h := sha256.New()
	if _, e = io.Copy(h, io.LimitReader(f, 1<<30+1)); e != nil {
		return "", e
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func copyMigrationFile(from, to string) error {
	if safePath(from) != nil || safePath(to) != nil {
		return errClientMigration
	}
	in, e := os.Open(from)
	if e != nil {
		return e
	}
	defer in.Close()
	st, e := in.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() > 1<<30 {
		return errClientMigration
	}
	out, e := os.CreateTemp(filepath.Dir(to), ".migration-*")
	if e != nil {
		return e
	}
	defer os.Remove(out.Name())
	defer out.Close()
	if e = out.Chmod(0600); e == nil {
		_, e = io.Copy(out, io.LimitReader(in, 1<<30+1))
	}
	if e == nil {
		e = out.Sync()
	}
	if e == nil {
		e = out.Close()
	}
	if e == nil {
		e = os.Rename(out.Name(), to)
	}
	if e == nil {
		e = syncDirectory(filepath.Dir(to))
	}
	return e
}
