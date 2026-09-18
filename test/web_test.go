package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"osmonitor/internal/web"
)

// The embedded dist directory may hold either the real Nuxt build or only the
// .gitkeep placeholder, so these tests assert behaviour that holds in both cases.

func TestWebServesHTMLAtRoot(t *testing.T) {
	ts := httptest.NewServer(web.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	// The Nuxt build sets its <title> client-side, so check for a document, not a title.
	if !strings.Contains(strings.ToLower(string(body)), "<!doctype html") {
		t.Error("root page is not an HTML document")
	}
}

func TestWebFallsBackToIndexForClientRoutes(t *testing.T) {
	ts := httptest.NewServer(web.Handler())
	defer ts.Close()

	root, _ := http.Get(ts.URL + "/")
	rootBody, _ := io.ReadAll(root.Body)
	root.Body.Close()

	deep, err := http.Get(ts.URL + "/some/client/route")
	if err != nil {
		t.Fatal(err)
	}
	defer deep.Body.Close()
	deepBody, _ := io.ReadAll(deep.Body)

	if deep.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", deep.StatusCode)
	}
	if string(deepBody) != string(rootBody) {
		t.Error("deep route did not return the same document as /")
	}
	if cc := deep.Header.Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("Cache-Control = %q, want no-cache for the HTML shell", cc)
	}
}

func TestWebMissingAssetIs404NotFallback(t *testing.T) {
	ts := httptest.NewServer(web.Handler())
	defer ts.Close()

	for _, p := range []string{"/_nuxt/does-not-exist.js", "/_fonts/nope.woff2"} {
		resp, err := http.Get(ts.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s: status = %d, want 404", p, resp.StatusCode)
		}
	}
}

func TestWebRejectsNonGET(t *testing.T) {
	ts := httptest.NewServer(web.Handler())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", resp.StatusCode)
	}
}
