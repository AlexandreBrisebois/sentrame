package main

import "testing"

func TestOriginAcceptsOnlyAnExplicitHTTPAuthority(t *testing.T) {
	for _, origin := range []string{"http://127.0.0.1:8096", "https://sentrame.example.com", "http://[::1]:8096"} {
		if !validOrigin(origin) {
			t.Fatalf("valid origin rejected %q", origin)
		}
	}
	for _, origin := range []string{"", "/", "https://", "http://:8096", "ftp://host", "https://user:password@host", "https://host/", "https://host/path", "https://host?", "https://host?q=1", "https://host#", "https://host#fragment"} {
		if validOrigin(origin) {
			t.Fatalf("unsafe origin accepted %q", origin)
		}
	}
}
