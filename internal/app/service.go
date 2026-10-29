package app

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"sync"
	"time"
)

type session struct {
	CSRF              string
	Created, Last     time.Time
	AppID, Generation string
	InviteID          int64
	Expires           time.Time
}
type Service struct {
	Store                             *Store
	BrowserToken, ServerToken, Origin string
	password                          []byte
	Mu                                *sync.Mutex
	Registry                          *Registry
	selected                          *AppEntry
	redemptionReserved                bool
	credentials                       []string
	queries                           chan struct{}
	sessions                          map[string]session
	failures                          []time.Time
	slots                             chan struct{}
	streams                           chan struct{}
	logins                            chan struct{}
	ReceiverURL                       string
	closeReceiver                     func()
	static                            http.Handler
}
type host struct{}

func (host) GetExtensions() map[component.ID]component.Component { return nil }
func New(s *Store, pass, browser, server, origin string) (*Service, error) {
	if len(pass) < 12 || browser == "" || server == "" || browser == server {
		return nil, fmt.Errorf("passphrase must have 12 characters; producer tokens must be distinct and nonempty")
	}
	hash, e := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if e != nil {
		return nil, e
	}
	return &Service{Mu: &sync.Mutex{}, queries: make(chan struct{}, 8), Store: s, BrowserToken: browser, ServerToken: server, Origin: origin, password: hash, sessions: map[string]session{}, slots: make(chan struct{}, 8), streams: make(chan struct{}, 8), logins: make(chan struct{}, 2)}, nil
}
func (s *Service) StartReceiver() error {
	ln, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		return e
	}
	addr := ln.Addr().String()
	ln.Close()
	f := otlpreceiver.NewFactory()
	cfg := f.CreateDefaultConfig().(*otlpreceiver.Config)
	cfg.Protocols.GRPC = configoptional.None[configgrpc.ServerConfig]()
	hc := cfg.Protocols.HTTP.GetOrInsertDefault()
	hc.ServerConfig.NetAddr.Endpoint = addr
	hc.ServerConfig.IncludeMetadata = true
	hc.ServerConfig.MaxRequestBodySize = 4 << 20
	hc.ServerConfig.ReadTimeout = 10 * time.Second
	hc.ServerConfig.ReadHeaderTimeout = 5 * time.Second
	cfg.Protocols.HTTP = configoptional.Some(*hc)
	set := receiver.Settings{ID: component.NewID(f.Type()), TelemetrySettings: component.TelemetrySettings{Logger: zap.NewNop(), MeterProvider: noop.NewMeterProvider(), TracerProvider: tracenoop.NewTracerProvider(), Resource: pcommon.NewResource()}}
	lc, _ := consumer.NewLogs(s.consumeLogs)
	tc, _ := consumer.NewTraces(s.consumeTraces)
	mc, _ := consumer.NewMetrics(s.consumeMetrics)
	l, e := f.CreateLogs(context.Background(), set, cfg, lc)
	if e != nil {
		return e
	}
	t, e := f.CreateTraces(context.Background(), set, cfg, tc)
	if e != nil {
		return e
	}
	m, e := f.CreateMetrics(context.Background(), set, cfg, mc)
	if e != nil {
		return e
	}
	for _, c := range []component.Component{l, t, m} {
		if e = c.Start(context.Background(), host{}); e != nil {
			return e
		}
	}
	s.ReceiverURL = "http://" + addr
	s.closeReceiver = func() {
		for _, c := range []component.Component{l, t, m} {
			c.Shutdown(context.Background())
		}
	}
	return nil
}
func (s *Service) Close() {
	if s.closeReceiver != nil {
		s.closeReceiver()
		s.closeReceiver = nil
	}
}

// NewManaged shares one instance authority and receiver across bounded app stores.
func NewManaged(registry *Registry, pass, origin string) (*Service, error) {
	s, e := New(nil, pass, "unused-browser", "unused-server", origin)
	if e != nil {
		return nil, e
	}
	s.BrowserToken = ""
	s.ServerToken = ""
	s.Registry = registry
	return s, nil
}
func (s *Service) forApp(a AppEntry) *Service {
	secrets := []string{}
	for _, entry := range s.Registry.current.Apps {
		secrets = append(secrets, entry.Browser, entry.Server)
	}
	return &Service{credentials: secrets, Store: s.Registry.stores[a.ID], BrowserToken: a.Browser, ServerToken: a.Server, Origin: s.Origin, password: s.password, Mu: s.Mu, sessions: s.sessions, slots: s.slots, streams: s.streams, logins: s.logins, queries: s.queries, Registry: s.Registry, selected: &a}
}
