package ynab_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
	"thde.io/ynab"
)

func TestWithBearerToken(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, err := ynab.NewClientWithResponses(srv.URL, ynab.WithBearerToken("my-token"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.GetUserWithResponse(t.Context()); err != nil {
		t.Fatal(err)
	}

	if want := "Bearer my-token"; gotHeader != want {
		t.Errorf("Authorization = %q, want %q", gotHeader, want)
	}
}

func TestWithUserAgent(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, err := ynab.NewClientWithResponses(srv.URL, ynab.WithUserAgent("my-agent/1.0"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.GetUserWithResponse(t.Context()); err != nil {
		t.Fatal(err)
	}

	if want := "my-agent/1.0"; gotHeader != want {
		t.Errorf("User-Agent = %q, want %q", gotHeader, want)
	}
}

func TestWithRateLimiter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, err := ynab.NewClientWithResponses(srv.URL, ynab.WithRateLimiter(nil))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.GetUserWithResponse(t.Context()); err != nil {
		t.Errorf("unexpected error with token available: %v", err)
	}
}

func TestWithRateLimiter_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Consume the one available token so the next request must wait.
	limiter := rate.NewLimiter(rate.Every(time.Hour), 1)
	limiter.Reserve()

	client, err := ynab.NewClientWithResponses(srv.URL, ynab.WithRateLimiter(limiter))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, err := client.GetUserWithResponse(ctx); err == nil {
		t.Error("expected error for cancelled context, got nil")
	}
}
