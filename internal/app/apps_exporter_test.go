package app

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"testing"
	"time"
)

func TestManagedStandardGoExporterNeedsOnlyEndpointAndHeader(t *testing.T) {
	m, _, server, _ := managedFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	for _, app := range []AppEntry{a, b} {
		ctx := context.Background()
		exp, e := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(server.URL+"/v1/logs"), otlploghttp.WithHeaders(map[string]string{"Authorization": "Bearer " + app.Server}), otlploghttp.WithInsecure())
		if e != nil {
			t.Fatal(e)
		}
		provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
		var record otellog.Record
		record.SetEventName("ordinary.standard.log")
		record.SetTimestamp(time.Now())
		provider.Logger("existing-provider").Emit(ctx, record)
		if e = provider.Shutdown(ctx); e != nil {
			t.Fatal(e)
		}
	}
	for _, app := range []AppEntry{a, b} {
		es, e := m.stores[app.ID].Evidence()
		if e != nil || len(es) != 1 || es[0].Name != "ordinary.standard.log" || es[0].Source != "server" || es[0].Scope != "" {
			t.Fatal("standard provider attribution/receipt failed")
		}
	}
}
