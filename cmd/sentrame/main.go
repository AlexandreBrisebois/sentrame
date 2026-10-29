package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sentrame/internal/app"
	"syscall"
	"time"
)

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func validOrigin(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Hostname() != "" && parsed.User == nil && value == parsed.Scheme+"://"+parsed.Host
}
func main() {
	dbPath := env("SENTRAME_DB", ".sentrame/local.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
		log.Fatal("Cannot create database directory")
	}
	origin := env("SENTRAME_ORIGIN", "http://127.0.0.1:8096")
	if !validOrigin(origin) {
		log.Fatal("SENTRAME_ORIGIN must be an absolute HTTP(S) origin without a path, credentials, query or fragment")
	}

	lease, e := os.OpenFile(dbPath+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		log.Fatal("Cannot open local database lease")
	}
	defer lease.Close()
	if e = syscall.Flock(int(lease.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		log.Fatal("Database already owned by a running local service or reset")
	}
	s, e := app.Open(dbPath)
	if e != nil {
		log.Fatal(e)
	}
	defer s.DB.Close()
	if e = s.Cleanup(); e != nil {
		log.Fatal(e)
	}
	svc, e := app.New(s, os.Getenv("SENTRAME_PASSPHRASE"), os.Getenv("SENTRAME_BROWSER_TOKEN"), os.Getenv("SENTRAME_SERVER_TOKEN"), origin)
	if e != nil {
		log.Fatal(e)
	}
	if e = svc.StartReceiver(); e != nil {
		log.Fatal(e)
	}
	defer svc.Close()
	srv := &http.Server{Addr: env("SENTRAME_ADDR", "127.0.0.1:8096"), Handler: svc.Handler("web/dist"), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, IdleTimeout: 30 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if err := s.Cleanup(); err != nil {
					log.Print("Lifecycle cleanup unavailable")
				}
			}
		}
	}()
	go func() { <-ctx.Done(); srv.Shutdown(context.Background()) }()
	log.Print("Sentrame started; viewer origin: ", origin)
	if e = srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		log.Fatal(e)
	}
}
