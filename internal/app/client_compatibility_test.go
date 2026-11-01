package app

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestClientHistoryPreservesMixedRecordsWithoutInventedTransitions(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "history.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.DB.Close()
	now := time.Now()
	s.Now = func() time.Time { return now }
	d := Definition{Version: 1, Mappings: []Mapping{}}
	raw, _ := json.Marshal(d)
	if _, err = s.DB.Exec(`UPDATE definitions SET data=? WHERE version=1`, string(raw)); err != nil {
		t.Fatal(err)
	}
	es := []Evidence{
		{Source: "browser", Signal: "logs", Name: "A", Scope: "mixed", Tab: "tab", Record: "a", Time: now.Add(-5 * time.Minute).UnixNano(), Projection: "structural-v2"},
		{Source: "server", Signal: "logs", Name: "B", Scope: "mixed", Record: "b", Time: now.Add(-4 * time.Minute).UnixNano(), Projection: "structural-v2"},
		{Source: "browser", Signal: "logs", Name: "C", Scope: "mixed", Tab: "tab", Record: "c", Time: now.Add(-3 * time.Minute).UnixNano(), Projection: "structural-v2"},
		{Source: "browser", Signal: "logs", Name: "Home", Scope: "pure", Tab: "tab", Record: "pure", Time: now.Add(-2 * time.Minute).UnixNano(), Projection: "structural-v2"},
	}
	if err = s.Insert(es); err != nil {
		t.Fatal(err)
	}
	original, err := s.rawEvidence()
	if err != nil {
		t.Fatal(err)
	}
	if err = stageClientHistory(s.DB); err != nil {
		t.Fatal(err)
	}
	s.ClientCutover = now.UnixNano()
	got, err := s.EvidenceLimit(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Scope != "pure" || got[0].Source != "client" {
		t.Fatalf("eligible client evidence: %+v", got)
	}
	all, err := s.rawEvidence()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 4 || all[1].Source != "server" || all[0].Source != "browser" {
		t.Fatal("migration rewrote legacy evidence")
	}
	defs, _ := s.Definitions()
	m := Derive(got, defs, now.UnixNano())
	if len(m.Edges) != 0 {
		t.Fatal("invented edge across removed legacy records")
	}
	if s.ClientCompatibility()["excludedRecords"] != 3 {
		t.Fatal("missing compatibility explanation")
	}
	pure := SavedView{Scope: "pure", From: now.Add(-time.Hour).Format(time.RFC3339Nano), To: now.Add(time.Hour).Format(time.RFC3339Nano)}
	if ok, err := s.ClientReferenceCompatible(pure); err != nil || !ok {
		t.Fatalf("equivalent client reference unavailable: %v", err)
	}
	pure.Scope = "mixed"
	if ok, err := s.ClientReferenceCompatible(pure); err != nil || ok {
		t.Fatalf("mixed reference incorrectly available: %v", err)
	}
	// The retained key's client namespace aliases only old browser identities.
	replay := original[3]
	replay.ID = 0
	replay.Source = "client"
	replay.Projection = projectionPolicy
	if err = s.Insert([]Evidence{replay}); err != nil {
		t.Fatalf("client replay failed: %v", err)
	}
	all, _ = s.rawEvidence()
	if len(all) != 4 {
		t.Fatal("old client replay duplicated history")
	}
	// A server record expiring cannot remove the persisted mixed-history barrier.
	if _, err = s.DB.Exec(`DELETE FROM evidence WHERE id=?`, original[1].ID); err != nil {
		t.Fatal(err)
	}
	got, err = s.Evidence()
	if err != nil || len(got) != 1 {
		t.Fatalf("expiry fabricated eligible mixed history: %v %+v", err, got)
	}
	// A delayed old mixed observation is retained but remains excluded.
	late := Evidence{Source: "client", Signal: "logs", Name: "Late", Scope: "mixed", Tab: "tab", Record: "late", Time: now.Add(-time.Minute).UnixNano(), Projection: projectionPolicy}
	if err = s.Insert([]Evidence{late}); err != nil {
		t.Fatal(err)
	}
	got, err = s.Evidence()
	if err != nil || len(got) != 1 {
		t.Fatalf("late evidence crossed legacy boundary: %v", err)
	}
	// A continuation after cutover is eligible only with its missing beginning explicit.
	late.Record = "new-current"
	late.Time = now.Add(time.Minute).UnixNano()
	if err = s.Insert([]Evidence{late}); err != nil {
		t.Fatal(err)
	}
	got, err = s.Evidence()
	if err != nil || len(got) != 2 || got[1].Attrs["_partial_start"] != "true" {
		t.Fatalf("cutover continuation lost partial marker: %v %+v", err, got)
	}
}
