package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otellog "go.opentelemetry.io/otel/log"
	apimetric "go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func protocolBody(t *testing.T, signal, encoding, id string) []byte {
	t.Helper()
	now := pcommon.NewTimestampFromTime(time.Now())
	var body []byte
	var err error
	switch signal {
	case "logs":
		ls := plog.NewLogs()
		r := ls.ResourceLogs().AppendEmpty()
		r.Resource().Attributes().PutStr("private.password", "PROTOCOL-PRIVATE-SENTINEL")
		l := r.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
		l.SetEventName(id)
		l.SetTimestamp(now)
		l.Body().SetStr("PROTOCOL-PRIVATE-SENTINEL raw body")
		l.Attributes().PutStr("event.id", id)
		l.Attributes().PutStr("private.user", "PROTOCOL-PRIVATE-SENTINEL")
		req := plogotlp.NewExportRequestFromLogs(ls)
		if encoding == "json" {
			body, err = req.MarshalJSON()
		} else {
			body, err = req.MarshalProto()
		}
	case "traces":
		ts := ptrace.NewTraces()
		r := ts.ResourceSpans().AppendEmpty()
		r.Resource().Attributes().PutStr("private.password", "PROTOCOL-PRIVATE-SENTINEL")
		span := r.ScopeSpans().AppendEmpty().Spans().AppendEmpty()
		span.SetName(id)
		span.SetStartTimestamp(now)
		span.SetEndTimestamp(pcommon.Timestamp(uint64(now) + 1000000))
		h := sha256.Sum256([]byte(id))
		var tr pcommon.TraceID
		copy(tr[:], h[:16])
		var sp pcommon.SpanID
		copy(sp[:], h[16:24])
		span.SetTraceID(tr)
		span.SetSpanID(sp)
		span.Attributes().PutStr("private.user", "PROTOCOL-PRIVATE-SENTINEL")
		ev := span.Events().AppendEmpty()
		ev.SetName(id + ".event")
		ev.SetTimestamp(now)
		ev.Attributes().PutStr("exception.type", "SyntheticProtocolError")
		ev.Attributes().PutStr("exception.message", "Synthetic protocol failure password=PROTOCOL-PRIVATE-SENTINEL")
		req := ptraceotlp.NewExportRequestFromTraces(ts)
		if encoding == "json" {
			body, err = req.MarshalJSON()
		} else {
			body, err = req.MarshalProto()
		}
	case "metrics":
		ms := pmetric.NewMetrics()
		r := ms.ResourceMetrics().AppendEmpty()
		r.Resource().Attributes().PutStr("private.password", "PROTOCOL-PRIVATE-SENTINEL")
		m := r.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
		m.SetName(id)
		p := m.SetEmptyGauge().DataPoints().AppendEmpty()
		p.SetTimestamp(now)
		p.SetIntValue(7)
		p.Attributes().PutStr("private.user", "PROTOCOL-PRIVATE-SENTINEL")
		req := pmetricotlp.NewExportRequestFromMetrics(ms)
		if encoding == "json" {
			body, err = req.MarshalJSON()
		} else {
			body, err = req.MarshalProto()
		}
	default:
		t.Fatal("unknown signal")
	}
	if err != nil {
		t.Fatal(err)
	}
	return body
}
func protocolRequest(t *testing.T, c *http.Client, base, signal, contentType, encoding, token string, b []byte) (int, []byte) {
	t.Helper()
	r, e := http.NewRequest("POST", base+"/v1/"+signal, bytes.NewReader(b))
	if e != nil {
		t.Fatal(e)
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	r.Header.Set("Content-Type", contentType)
	if encoding != "" {
		r.Header.Set("Content-Encoding", encoding)
	}
	resp, e := c.Do(r)
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if e != nil {
		t.Fatal(e)
	}
	return resp.StatusCode, raw
}
func compressed(t *testing.T, b []byte) []byte {
	t.Helper()
	var out bytes.Buffer
	w := gzip.NewWriter(&out)
	if _, e := w.Write(b); e != nil {
		t.Fatal(e)
	}
	if e := w.Close(); e != nil {
		t.Fatal(e)
	}
	return out.Bytes()
}

func TestOTLPRejectedResponseUsesStandardStatusAndSanitizedMessage(t *testing.T) {
	_, server, c := fixtureService(t)
	code, raw := protocolRequest(t, c, server.URL, "logs", "application/json", "", "browser-test", []byte(`{"resourceLogs": "PROTOCOL-PRIVATE-SENTINEL"}`))
	if code != 400 {
		t.Fatalf("malformed request status %d", code)
	}
	var st statuspb.Status
	if e := protojson.Unmarshal(raw, &st); e != nil {
		t.Fatalf("error must be standard google.rpc.Status: %v", e)
	}
	if st.Code != 3 {
		t.Fatalf("permanent rejection must preserve InvalidArgument code 3, got %d", st.Code)
	}
	if bytes.Contains(raw, []byte("PROTOCOL-PRIVATE-SENTINEL")) {
		t.Fatal("malformed raw body leaked in rejection")
	}
}

func TestOTLPHTTPWireMatrixAndPrivacyProjection(t *testing.T) {
	_, server, c := fixtureService(t)
	expected := map[string]string{}
	for _, signal := range []string{"logs", "traces", "metrics"} {
		for _, encoding := range []string{"protobuf", "json"} {
			for _, zip := range []bool{false, true} {
				compression := "none"
				if zip {
					compression = "gzip"
				}
				name := "wire." + signal + "." + encoding + "." + compression
				t.Run(signal+"/"+encoding+"/"+compression, func(t *testing.T) {
					contentType := "application/x-protobuf"
					body := protocolBody(t, signal, encoding, name)
					if encoding == "json" {
						contentType = "application/json"
						body = append(body[:len(body)-1], []byte(`,"unknownProtocolField":"PROTOCOL-PRIVATE-SENTINEL"}`)...)
					} else {
						body = append(body, 0xa0, 0x06, 0x01)
					}
					ce := ""
					if zip {
						body = compressed(t, body)
						ce = "gzip"
					}
					code, raw := protocolRequest(t, c, server.URL, signal, contentType, ce, "server-test", body)
					if code != 200 {
						t.Fatalf("valid standard request rejected: %d %s", code, raw)
					}
					expected[name] = signal
					if signal == "traces" {
						expected[name+".event"] = "span-events"
					}
				})
			}
		}
	}
	loginTest(t, c, server.URL)
	resp, e := c.Get(server.URL + "/api/snapshot")
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	raw, e := io.ReadAll(resp.Body)
	if e != nil {
		t.Fatal(e)
	}
	if bytes.Contains(raw, []byte("PROTOCOL-PRIVATE-SENTINEL")) {
		t.Fatal("unsafe bodies/attributes/unknown fields reached protected output")
	}
	var snap struct {
		Evidence []Evidence `json:"evidence"`
	}
	if e = json.Unmarshal(raw, &snap); e != nil {
		t.Fatal(e)
	}
	if len(snap.Evidence) != 16 {
		t.Fatalf("12 requests must retain 16 sanitized records incl4 span-events, got%d", len(snap.Evidence))
	}
	for _, x := range snap.Evidence {
		if expected[x.Name] != x.Signal {
			t.Errorf("missing/wrong signal for %s: %s", x.Name, x.Signal)
		}
		if x.Scope != "" || x.Tab != "" {
			t.Fatal("ordinary telemetry acquired fabricated journey context")
		}
		delete(expected, x.Name)
	}
	if len(expected) != 0 {
		t.Fatalf("protected read omitted qualified records %v", expected)
	}
}

func TestOTLPHTTPPermanentErrorsAuthAndBoundedDecompression(t *testing.T) {
	svc, server, c := fixtureService(t)
	for _, signal := range []string{"logs", "traces", "metrics"} {
		t.Run(signal, func(t *testing.T) {
			valid := protocolBody(t, signal, "protobuf", "invalid."+signal)
			tests := []struct {
				name, ct, ce, token string
				body                []byte
				want                int
			}{
				{"no-producer", "application/x-protobuf", "", "", valid, 401},
				{"wrong-producer", "application/x-protobuf", "", "forged", valid, 401},
				{"malformed-binary", "application/x-protobuf", "", "server-test", []byte{0xff, 0xff, 0xff}, 400},
				{"malformed-json", "application/json", "", "server-test", []byte(`{"invalid":"PROTOCOL-PRIVATE-SENTINEL"`), 400},
				{"wrong-content-type", "text/plain", "", "server-test", valid, 415},
				{"unknown-encoding", "application/x-protobuf", "br", "server-test", valid, 415},
				{"invalid-gzip", "application/x-protobuf", "gzip", "server-test", []byte("PROTOCOL-PRIVATE-SENTINEL not gzip"), 400},
				{"request-body-limit", "application/x-protobuf", "", "server-test", bytes.Repeat([]byte{0}, (1<<20)+1), 400},
				{"decompression-limit", "application/json", "gzip", "server-test", compressed(t, bytes.Repeat([]byte(" "), (4<<20)+1)), 400},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					code, raw := protocolRequest(t, c, server.URL, signal, tt.ct, tt.ce, tt.token, tt.body)
					if code != tt.want {
						t.Fatalf("status%d want%d", code, tt.want)
					}
					if bytes.Contains(raw, []byte("PROTOCOL-PRIVATE-SENTINEL")) {
						t.Fatal("raw rejection content leaked")
					}
					var rpc statuspb.Status
					var err error
					if tt.ct == "application/x-protobuf" {
						err = proto.Unmarshal(raw, &rpc)
					} else {
						err = protojson.Unmarshal(raw, &rpc)
					}
					if err != nil {
						t.Fatalf("rejection is not standard OTLP Status: %v", err)
					}
					wantCode := int32(3)
					if tt.want == 401 {
						wantCode = 16
					}
					if rpc.Code != wantCode || rpc.Message == "" || len(rpc.Details) != 0 {
						t.Fatalf("unexpected sanitized status code/message/details: %v", &rpc)
					}
				})
			}
		})
	}
	es, e := svc.Store.Evidence()
	if e != nil {
		t.Fatal(e)
	}
	if len(es) != 0 {
		t.Fatalf("rejected requests committed%d records", len(es))
	}
	for i := 0; i < cap(svc.slots); i++ {
		svc.slots <- struct{}{}
	}
	defer func() {
		for i := 0; i < cap(svc.slots); i++ {
			<-svc.slots
		}
	}()
	code, _ := protocolRequest(t, c, server.URL, "logs", "application/x-protobuf", "", "server-test", protocolBody(t, "logs", "protobuf", "bounded.overload"))
	if code != 503 {
		t.Fatalf("bounded overload must be retryable503, got%d", code)
	}
}

// These cells use unmodified maintained Go SDK exporters configured with ordinary
// endpoint/headers/compression options; wire fixtures above do not establish them.
func TestOTLPHTTPStandardSDKExportersBinaryNoneAndGzip(t *testing.T) {
	_, server, c := fixtureService(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	expected := map[string]string{}
	for _, zip := range []bool{false, true} {
		compression := "none"
		if zip {
			compression = "gzip"
		}
		t.Run("logs/"+compression, func(t *testing.T) {
			comp := otlploghttp.NoCompression
			if zip {
				comp = otlploghttp.GzipCompression
			}
			exp, e := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(server.URL+"/v1/logs"), otlploghttp.WithHeaders(map[string]string{"Authorization": "Bearer server-test"}), otlploghttp.WithInsecure(), otlploghttp.WithCompression(comp))
			if e != nil {
				t.Fatal(e)
			}
			provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
			var rec otellog.Record
			name := "sdk.logs." + compression
			rec.SetEventName(name)
			rec.SetTimestamp(time.Now())
			rec.SetBody(attribute.StringValue("PROTOCOL-PRIVATE-SENTINEL SDK body"))
			rec.AddAttributes(attribute.String("private.user", "PROTOCOL-PRIVATE-SENTINEL"))
			provider.Logger("standard-http-sdk").Emit(ctx, rec)
			if e = provider.Shutdown(ctx); e != nil {
				t.Fatal(e)
			}
			expected[name] = "logs"
		})
		t.Run("traces/"+compression, func(t *testing.T) {
			comp := otlptracehttp.NoCompression
			if zip {
				comp = otlptracehttp.GzipCompression
			}
			exp, e := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(server.URL+"/v1/traces"), otlptracehttp.WithHeaders(map[string]string{"Authorization": "Bearer server-test"}), otlptracehttp.WithInsecure(), otlptracehttp.WithCompression(comp))
			if e != nil {
				t.Fatal(e)
			}
			provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exp))
			name := "sdk.traces." + compression
			_, span := provider.Tracer("standard-http-sdk").Start(ctx, name)
			span.SetAttributes(attribute.String("private.user", "PROTOCOL-PRIVATE-SENTINEL"))
			span.End()
			if e = provider.Shutdown(ctx); e != nil {
				t.Fatal(e)
			}
			expected[name] = "traces"
		})
		t.Run("metrics/"+compression, func(t *testing.T) {
			comp := otlpmetrichttp.NoCompression
			if zip {
				comp = otlpmetrichttp.GzipCompression
			}
			exp, e := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(server.URL+"/v1/metrics"), otlpmetrichttp.WithHeaders(map[string]string{"Authorization": "Bearer server-test"}), otlpmetrichttp.WithInsecure(), otlpmetrichttp.WithCompression(comp))
			if e != nil {
				t.Fatal(e)
			}
			reader := sdkmetric.NewManualReader()
			provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			name := "sdk.metrics." + compression
			counter, e := provider.Meter("standard-http-sdk").Int64Counter(name)
			if e != nil {
				t.Fatal(e)
			}
			counter.Add(ctx, 7, apimetric.WithAttributes(attribute.String("private.user", "PROTOCOL-PRIVATE-SENTINEL")))
			var rm metricdata.ResourceMetrics
			if e = reader.Collect(ctx, &rm); e != nil {
				t.Fatal(e)
			}
			if e = exp.Export(ctx, &rm); e != nil {
				t.Fatal(e)
			}
			if e = provider.Shutdown(ctx); e != nil {
				t.Fatal(e)
			}
			if e = exp.Shutdown(ctx); e != nil {
				t.Fatal(e)
			}
			expected[name] = "metrics"
		})
	}
	loginTest(t, c, server.URL)
	resp, e := c.Get(server.URL + "/api/snapshot")
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	raw, e := io.ReadAll(resp.Body)
	if e != nil {
		t.Fatal(e)
	}
	if bytes.Contains(raw, []byte("PROTOCOL-PRIVATE-SENTINEL")) {
		t.Fatal("SDK private attributes/body reached owner read")
	}
	var snap struct {
		Evidence []Evidence `json:"evidence"`
	}
	if e = json.Unmarshal(raw, &snap); e != nil {
		t.Fatal(e)
	}
	if len(snap.Evidence) != 6 {
		t.Fatalf("six exporter cells must retain six sanitized records/points, got%d", len(snap.Evidence))
	}
	for _, x := range snap.Evidence {
		if expected[x.Name] != x.Signal {
			t.Fatalf("missing or wrong SDK receipt %s/%s", x.Name, x.Signal)
		}
		if x.Scope != "" {
			t.Fatal("SDK technical telemetry acquired fabricated journey")
		}
		delete(expected, x.Name)
	}
	if len(expected) != 0 {
		t.Fatalf("missing protected SDK evidence %v", expected)
	}
}

func TestOTLPResourceAttributeLimitWholeRequest(t *testing.T) {
	svc, server, c := fixtureService(t)
	for _, signal := range []string{"logs", "traces", "metrics"} {
		t.Run(signal, func(t *testing.T) {
			b := protocolBody(t, signal, "protobuf", "resource-limit."+signal)
			switch signal {
			case "logs":
				q := plogotlp.NewExportRequest()
				if e := q.UnmarshalProto(b); e != nil {
					t.Fatal(e)
				}
				a := q.Logs().ResourceLogs().At(0).Resource().Attributes()
				for i := 0; i < 65; i++ {
					a.PutInt(fmt.Sprintf("field.%d", i), int64(i))
				}
				b, _ = q.MarshalProto()
			case "traces":
				q := ptraceotlp.NewExportRequest()
				if e := q.UnmarshalProto(b); e != nil {
					t.Fatal(e)
				}
				a := q.Traces().ResourceSpans().At(0).Resource().Attributes()
				for i := 0; i < 65; i++ {
					a.PutInt(fmt.Sprintf("field.%d", i), int64(i))
				}
				b, _ = q.MarshalProto()
			case "metrics":
				q := pmetricotlp.NewExportRequest()
				if e := q.UnmarshalProto(b); e != nil {
					t.Fatal(e)
				}
				a := q.Metrics().ResourceMetrics().At(0).Resource().Attributes()
				for i := 0; i < 65; i++ {
					a.PutInt(fmt.Sprintf("field.%d", i), int64(i))
				}
				b, _ = q.MarshalProto()
			}
			before, e := svc.Store.Evidence()
			if e != nil {
				t.Fatal(e)
			}
			code, _ := protocolRequest(t, c, server.URL, signal, "application/x-protobuf", "", "server-test", b)
			after, e := svc.Store.Evidence()
			if e != nil {
				t.Fatal(e)
			}
			if code != 400 || len(after) != len(before) {
				t.Fatalf("resource cardinality bound must atomically reject (400 and zero new records), got status%d new%d", code, len(after)-len(before))
			}
		})
	}
}

func TestOTLPScopeSpanEventAndMetricPointAttributeLimits(t *testing.T) {
	svc, server, c := fixtureService(t)
	fill := func(a pcommon.Map) {
		for i := 0; i < 65; i++ {
			a.PutInt(fmt.Sprintf("bounded.%d", i), int64(i))
		}
	}
	assertRejected := func(t *testing.T, signal string, b []byte) {
		t.Helper()
		before, e := svc.Store.Evidence()
		if e != nil {
			t.Fatal(e)
		}
		code, _ := protocolRequest(t, c, server.URL, signal, "application/x-protobuf", "", "server-test", b)
		after, e := svc.Store.Evidence()
		if e != nil {
			t.Fatal(e)
		}
		if code != 400 || len(after) != len(before) {
			t.Fatalf("attribute bound must atomically reject, status%d newrecords%d", code, len(after)-len(before))
		}
	}
	for _, signal := range []string{"logs", "traces", "metrics"} {
		t.Run(signal+"/scope", func(t *testing.T) {
			b := protocolBody(t, signal, "protobuf", "scope-limit."+signal)
			switch signal {
			case "logs":
				q := plogotlp.NewExportRequest()
				if e := q.UnmarshalProto(b); e != nil {
					t.Fatal(e)
				}
				fill(q.Logs().ResourceLogs().At(0).ScopeLogs().At(0).Scope().Attributes())
				b, _ = q.MarshalProto()
			case "traces":
				q := ptraceotlp.NewExportRequest()
				if e := q.UnmarshalProto(b); e != nil {
					t.Fatal(e)
				}
				fill(q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Scope().Attributes())
				b, _ = q.MarshalProto()
			case "metrics":
				q := pmetricotlp.NewExportRequest()
				if e := q.UnmarshalProto(b); e != nil {
					t.Fatal(e)
				}
				fill(q.Metrics().ResourceMetrics().At(0).ScopeMetrics().At(0).Scope().Attributes())
				b, _ = q.MarshalProto()
			}
			assertRejected(t, signal, b)
		})
	}
	t.Run("traces/span-event", func(t *testing.T) {
		q := ptraceotlp.NewExportRequest()
		if e := q.UnmarshalProto(protocolBody(t, "traces", "protobuf", "span-event-limit")); e != nil {
			t.Fatal(e)
		}
		fill(q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Events().At(0).Attributes())
		b, e := q.MarshalProto()
		if e != nil {
			t.Fatal(e)
		}
		assertRejected(t, "traces", b)
	})
	t.Run("traces/span", func(t *testing.T) {
		q := ptraceotlp.NewExportRequest()
		if e := q.UnmarshalProto(protocolBody(t, "traces", "protobuf", "span-limit")); e != nil {
			t.Fatal(e)
		}
		fill(q.Traces().ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes())
		b, e := q.MarshalProto()
		if e != nil {
			t.Fatal(e)
		}
		assertRejected(t, "traces", b)
	})
	for _, kind := range []string{"gauge", "sum", "histogram", "exponential-histogram", "summary"} {
		t.Run("metrics/point/"+kind, func(t *testing.T) {
			ms := pmetric.NewMetrics()
			m := ms.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
			m.SetName("point-limit." + kind)
			now := pcommon.NewTimestampFromTime(time.Now())
			// A valid point preceding the invalid point may not leak through a partial commit.
			switch kind {
			case "gauge":
				ps := m.SetEmptyGauge().DataPoints()
				ps.AppendEmpty().SetTimestamp(now)
				p := ps.AppendEmpty()
				p.SetTimestamp(now)
				fill(p.Attributes())
			case "sum":
				ps := m.SetEmptySum().DataPoints()
				ps.AppendEmpty().SetTimestamp(now)
				p := ps.AppendEmpty()
				p.SetTimestamp(now)
				fill(p.Attributes())
			case "histogram":
				ps := m.SetEmptyHistogram().DataPoints()
				ps.AppendEmpty().SetTimestamp(now)
				p := ps.AppendEmpty()
				p.SetTimestamp(now)
				fill(p.Attributes())
			case "exponential-histogram":
				ps := m.SetEmptyExponentialHistogram().DataPoints()
				ps.AppendEmpty().SetTimestamp(now)
				p := ps.AppendEmpty()
				p.SetTimestamp(now)
				fill(p.Attributes())
			case "summary":
				ps := m.SetEmptySummary().DataPoints()
				ps.AppendEmpty().SetTimestamp(now)
				p := ps.AppendEmpty()
				p.SetTimestamp(now)
				fill(p.Attributes())
			}
			b, e := pmetricotlp.NewExportRequestFromMetrics(ms).MarshalProto()
			if e != nil {
				t.Fatal(e)
			}
			assertRejected(t, "metrics", b)
		})
	}
}
