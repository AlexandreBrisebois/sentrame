package app

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"net/http/httptest"
	"testing"
	"time"
)

func TestManagedStandardGoExporterNeedsOnlyEndpointAndHeader(t *testing.T) {
	m, _, server, _ := managedFixture(t)
	a := mustRegister(t, m, "Alpha")
	b := mustRegister(t, m, "Beta")
	for _, app := range []AppEntry{a, b} {
		ctx := context.Background()
		exp, e := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(server.URL+"/v1/logs"), otlploghttp.WithHeaders(map[string]string{"Authorization": "Bearer " + app.Key}), otlploghttp.WithInsecure())
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
		if e != nil || len(es) != 1 || es[0].Name != "ordinary.standard.log" || es[0].Source != "client" || es[0].Scope != "" {
			t.Fatal("standard provider attribution/receipt failed")
		}
	}
}

func TestMigratedAppStandardExporterReplayAcrossRegistryRestart(t *testing.T) {
	path, app, retired, _ := legacyClientMigrationFixture(t)
	if _, err := MigrateClientRegistry(path, true, true); err != nil {
		t.Fatal(err)
	}
	registry, err := OpenRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewManaged(registry, "synthetic owner passphrase", "http://localhost")
	if err != nil {
		registry.Close()
		t.Fatal(err)
	}
	if err = service.StartReceiver(); err != nil {
		registry.Close()
		t.Fatal(err)
	}
	server := httptest.NewServer(service.Handler(""))
	service.Origin = server.URL
	defer server.Close()
	defer service.Close()
	defer func() { registry.Close() }()
	ctx := context.Background()
	exporter, err := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(server.URL+"/v1/logs"), otlploghttp.WithHeaders(map[string]string{"Authorization": "Bearer " + app.Key}), otlploghttp.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exporter)))
	defer provider.Shutdown(ctx)
	var record otellog.Record
	stamp := time.Now().Add(-time.Second)
	record.SetEventName("client.screen.shown")
	record.SetTimestamp(stamp)
	record.SetObservedTimestamp(stamp)
	record.AddAttributes(attribute.String("event.id", "new-client-record"), attribute.String("session.id", "new-visit"), attribute.String("tab.id", "new-tab"))
	emit := func() {
		provider.Logger("existing-client-provider").Emit(ctx, record)
		if err := provider.ForceFlush(ctx); err != nil {
			t.Fatal(err)
		}
	}
	emit()
	if code := requestLogs(t, server.Client(), server.URL, retired, logFixture("retired")); code != 401 {
		t.Fatal("retired key was accepted", code)
	}
	service.Close()
	registry.Close()
	registry, err = OpenRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	service.Registry = registry
	if err = service.StartReceiver(); err != nil {
		t.Fatal(err)
	}
	emit()
	es, err := registry.stores[app.ID].Evidence()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range es {
		if e.Name == "client.screen.shown" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("standard exporter replay duplicated migrated evidence: %d", count)
	}
	defs, err := registry.stores[app.ID].Definitions()
	if err != nil {
		t.Fatal(err)
	}
	model := Derive(es, defs, time.Now().UnixNano())
	found := false
	for _, j := range model.Journeys {
		if j.Scope == "new-visit" {
			found = true
		}
	}
	if !found {
		t.Fatal("new client journey unavailable after migration and replay")
	}
}
