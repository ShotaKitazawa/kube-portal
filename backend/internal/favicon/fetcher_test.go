package favicon_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShotaKitazawa/kube-portal/internal/favicon"
)

func TestResolve_linkRelIcon(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<html><head><link rel="icon" href="/my-icon.png"></head></html>`)
		case "/my-icon.png":
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	got := favicon.New().Resolve(context.Background(), srv.URL)
	want := srv.URL + "/my-icon.png"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_linkRelShortcutIcon(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<html><head><link rel="shortcut icon" href="/favicon.png"></head></html>`)
		}
	}))
	defer srv.Close()

	got := favicon.New().Resolve(context.Background(), srv.URL)
	want := srv.URL + "/favicon.png"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_absoluteHref(t *testing.T) {
	cdn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer cdn.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><head><link rel="icon" href="%s/cdn-icon.png"></head></html>`, cdn.URL)
		}
	}))
	defer srv.Close()

	got := favicon.New().Resolve(context.Background(), srv.URL)
	want := cdn.URL + "/cdn-icon.png"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_fallbackFaviconIco(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<html><head><title>No icon link</title></head></html>`)
		case "/favicon.ico":
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	got := favicon.New().Resolve(context.Background(), srv.URL)
	want := srv.URL + "/favicon.ico"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_noFavicon(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<html><head><title>No icon</title></head></html>`)
		case "/favicon.ico":
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	got := favicon.New().Resolve(context.Background(), srv.URL)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestResolve_caches(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<html><head><link rel="icon" href="/icon.png"></head></html>`)
		}
	}))
	defer srv.Close()

	f := favicon.New()
	f.Resolve(context.Background(), srv.URL)
	f.Resolve(context.Background(), srv.URL)

	if calls != 1 {
		t.Errorf("expected 1 HTTP call due to caching, got %d", calls)
	}
}

func TestResolve_serverDown(t *testing.T) {
	got := favicon.New().Resolve(context.Background(), "http://localhost:1")
	if got != "" {
		t.Errorf("expected empty string for unreachable server, got %q", got)
	}
}
