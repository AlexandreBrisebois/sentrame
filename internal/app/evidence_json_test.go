package app

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEvidenceReadPreservesNanosecondsBeyondJavaScriptIntegerPrecision(t *testing.T) {
	const ns int64 = 1793217517832529123
	e := Evidence{Time: ns, Receipt: ns + 1, Observed: ns - 1}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"time":"1793217517832529123"`) {
		t.Fatalf("API rounded or exposed unsafe numeric nanos: %s", b)
	}
	var got Evidence
	if err = json.Unmarshal(b, &got); err != nil || got.Time != ns {
		t.Fatalf("decimal roundtrip %d %v", got.Time, err)
	}
	if err = json.Unmarshal([]byte(`{"time":1793217517832529123}`), &got); err == nil {
		t.Fatal("old numeric projection must not be silently decoded")
	}
}
