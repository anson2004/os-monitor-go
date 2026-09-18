package test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"osmonitor/internal/collector"
	"osmonitor/internal/server"
)

func newTestServer(t *testing.T) (*server.Server, *httptest.Server) {
	t.Helper()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := server.New(collector.New(), time.Hour, log)
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return srv, ts
}

func getJSON(t *testing.T, url string, into any) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		t.Fatalf("decode %s: %v", url, err)
	}
	return resp
}

func TestHealth(t *testing.T) {
	_, ts := newTestServer(t)

	var body map[string]string
	resp := getJSON(t, ts.URL+"/health", &body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if body["status"] != "ok" {
		t.Errorf(`body = %v, want {"status":"ok"}`, body)
	}
}

func TestLiveMetricsCollectsOnDemand(t *testing.T) {
	_, ts := newTestServer(t)

	var snap collector.Snapshot
	resp := getJSON(t, ts.URL+"/api/metrics/live", &snap)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if snap.Timestamp.IsZero() {
		t.Error("live snapshot has no timestamp")
	}
	if snap.CPU.LogicalCores == 0 {
		t.Error("live snapshot has no CPU data")
	}
}

func TestCachedMetricsIsEmptyUntilRunRefreshes(t *testing.T) {
	srv, ts := newTestServer(t)

	// Before Run has executed, the cache holds the zero Snapshot.
	var before collector.Snapshot
	getJSON(t, ts.URL+"/api/metrics", &before)
	if !before.Timestamp.IsZero() {
		t.Fatalf("expected empty cache before Run, got timestamp %v", before.Timestamp)
	}

	// Run performs one refresh immediately, then waits on its ticker.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go srv.Run(ctx)

	deadline := time.Now().Add(10 * time.Second)
	for {
		var after collector.Snapshot
		getJSON(t, ts.URL+"/api/metrics", &after)
		if !after.Timestamp.IsZero() {
			if after.CPU.LogicalCores == 0 {
				t.Error("cached snapshot has no CPU data")
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("cache was never populated by Run")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestRunStopsWhenContextCancelled(t *testing.T) {
	srv, _ := newTestServer(t)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		srv.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestUnknownRouteAndMethod(t *testing.T) {
	_, ts := newTestServer(t)

	cases := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/nope", http.StatusNotFound},
		{http.MethodPost, "/health", http.StatusMethodNotAllowed},
		{http.MethodDelete, "/api/metrics", http.StatusMethodNotAllowed},
	}
	for _, tc := range cases {
		req, _ := http.NewRequest(tc.method, ts.URL+tc.path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", tc.method, tc.path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != tc.want {
			t.Errorf("%s %s: status = %d, want %d", tc.method, tc.path, resp.StatusCode, tc.want)
		}
	}
}
