package app

import (
	"crypto/sha256"
	"encoding/hex"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"testing"
)

func TestTraceContextUsesCanonicalSourceHexIdentity(t *testing.T) {
	svc, server, c := fixtureService(t)
	b := protocolBody(t, "traces", "proto", "canonical.trace")
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	es, _ := svc.Store.Evidence()
	h := sha256.Sum256([]byte("canonical.trace"))
	if len(es) != 2 || es[0].Trace != hex.EncodeToString(h[:16]) || es[0].Span != hex.EncodeToString(h[16:24]) {
		t.Fatalf("source trace/span IDs were transformed: %+v", es)
	}
}
func TestInvalidStandardSpanIdentityRejectsWholeDirectRequest(t *testing.T) {
	svc, server, c := fixtureService(t)
	tr := ptrace.NewTraces()
	span := tr.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()
	span.SetName("missing.standard.ids")
	b, _ := ptraceotlp.NewExportRequestFromTraces(tr).MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 400 {
		t.Fatalf("invalid standard trace/span IDs accepted: %d", code)
	}
	es, _ := svc.Store.Evidence()
	if len(es) != 0 {
		t.Fatal("invalid standard span committed")
	}
}
func TestSpanEventInheritsOnlyVerifiedParentLogicalContext(t *testing.T) {
	svc, server, c := fixtureService(t)
	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(protocolBody(t, "traces", "proto", "discovery.entered")); err != nil {
		t.Fatal(err)
	}
	span := req.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)
	span.Attributes().PutStr("session.id", "verified-logical-scope")
	span.Attributes().PutStr("tab.id", "verified-tab")
	b, _ := req.MarshalProto()
	code, _ := protocolRequest(t, c, server.URL, "traces", "application/x-protobuf", "", "server-test", b)
	if code != 200 {
		t.Fatal(code)
	}
	es, _ := svc.Store.Evidence()
	if len(es) != 2 || es[1].Scope != "verified-logical-scope" || es[1].Tab != "verified-tab" {
		t.Fatalf("span event lost verified parent context: %+v", es)
	}
	defs, _ := svc.Store.Definitions()
	m := Derive(es, defs, svc.Store.Now().UnixNano())
	if len(m.Diagnostics) != 1 || m.Diagnostics[0].JourneyID == "" {
		t.Fatal("qualified span event remained unlinked despite verified parent relation")
	}
}
