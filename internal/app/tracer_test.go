package app

import (
	"path/filepath"
	"testing"
)

func TestTracerDurableSanitizedProjection(t *testing.T) {
	p := filepath.Join(t.TempDir(), "tracer.db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	e = s.Insert([]Evidence{{Name: "ordinary.event", Signal: "logs", Source: "browser", Attrs: map[string]string{}}})
	if e != nil {
		t.Fatal(e)
	}
	s.DB.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.DB.Close()
	es, e := s.Evidence()
	if e != nil {
		t.Fatal(e)
	}
	if len(es) != 1 || es[0].Name != "ordinary.event" {
		t.Fatalf("committed standard event must survive restart; got %v", es)
	}
}
