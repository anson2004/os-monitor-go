// Package web serves the compiled dashboard (ui/) embedded into the binary.
//
// The dist directory is populated by `make ui-build` (or the Dockerfile) from
// ui/.output/public. When it only contains the .gitkeep placeholder, a short
// notice page is served instead so the API still works in a Go-only build.
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var dist embed.FS

const notBuilt = `<!doctype html><meta charset="utf-8"><title>OS Monitor</title>
<body style="font-family:system-ui;margin:3rem;color:#334">
<h1>Dashboard not built</h1>
<p>The API is running. Build the UI with <code>make ui-build</code> and rebuild,
or query <a href="/api/metrics">/api/metrics</a> directly.</p>`

// Handler serves the embedded single-page app.
//
//   - Existing files are served as-is. Hashed assets under /_nuxt/ get a
//     one-year immutable cache header; missing ones are a 404, not a fallback.
//   - Any other path falls back to index.html so client-side routes work.
func Handler() http.Handler {
	files, err := fs.Sub(dist, "dist")
	if err != nil {
		panic("web: embedded dist directory missing: " + err.Error())
	}
	fileServer := http.FileServerFS(files)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "index.html"
		}

		if isFile(files, name) {
			if strings.HasPrefix(r.URL.Path, "/_nuxt/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(w, r)
			return
		}

		// Missing build asset: report it rather than masking it with index.html.
		if strings.HasPrefix(r.URL.Path, "/_nuxt/") || strings.HasPrefix(r.URL.Path, "/_fonts/") {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Cache-Control", "no-cache")
		if !isFile(files, "index.html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(notBuilt))
			return
		}
		http.ServeFileFS(w, r, files, "index.html")
	})
}

func isFile(fsys fs.FS, name string) bool {
	info, err := fs.Stat(fsys, name)
	return err == nil && !info.IsDir()
}
