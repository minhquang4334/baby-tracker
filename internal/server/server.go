package server

import (
	"io/fs"
	"net/http"

	"baby-care/internal/handler"
	"baby-care/internal/middleware"
	"baby-care/internal/store"
)

func New(st *store.Store, staticFS fs.FS) http.Handler {
	h := &handler.Handler{Store: st}
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Child API
	mux.HandleFunc("GET /api/v1/child", h.GetChild)
	mux.HandleFunc("POST /api/v1/child", h.CreateChild)
	mux.HandleFunc("PUT /api/v1/child", h.UpdateChild)

	// Sleep API
	mux.HandleFunc("GET /api/v1/sleep", h.ListSleep)
	mux.HandleFunc("POST /api/v1/sleep", h.CreateSleep)
	mux.HandleFunc("GET /api/v1/sleep/active", h.GetActiveSleep)
	mux.HandleFunc("PUT /api/v1/sleep/{logId}", h.UpdateSleep)
	mux.HandleFunc("DELETE /api/v1/sleep/{logId}", h.DeleteSleep)

	// Feeding API
	mux.HandleFunc("GET /api/v1/feeding", h.ListFeeding)
	mux.HandleFunc("POST /api/v1/feeding", h.CreateFeeding)
	mux.HandleFunc("GET /api/v1/feeding/active", h.GetActiveFeeding)
	mux.HandleFunc("PUT /api/v1/feeding/{logId}", h.UpdateFeeding)
	mux.HandleFunc("DELETE /api/v1/feeding/{logId}", h.DeleteFeeding)

	// Diaper API
	mux.HandleFunc("GET /api/v1/diaper", h.ListDiaper)
	mux.HandleFunc("POST /api/v1/diaper", h.CreateDiaper)
	mux.HandleFunc("PUT /api/v1/diaper/{logId}", h.UpdateDiaper)
	mux.HandleFunc("DELETE /api/v1/diaper/{logId}", h.DeleteDiaper)

	// Growth API
	mux.HandleFunc("GET /api/v1/growth", h.ListGrowth)
	mux.HandleFunc("POST /api/v1/growth", h.CreateGrowth)
	mux.HandleFunc("PUT /api/v1/growth/{logId}", h.UpdateGrowth)
	mux.HandleFunc("DELETE /api/v1/growth/{logId}", h.DeleteGrowth)

	// Summary API
	mux.HandleFunc("GET /api/v1/summary", h.GetSummary)

	// Static file server with SPA fallback
	static := http.FileServer(http.FS(staticFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the static file; fall back to index.html for SPA routing
		_, err := staticFS.Open(r.URL.Path[1:])
		if err != nil {
			// Serve index.html for SPA routes
			f, err2 := staticFS.Open("index.html")
			if err2 != nil {
				http.NotFound(w, r)
				return
			}
			f.Close()
			http.ServeFileFS(w, r, staticFS, "index.html")
			return
		}
		static.ServeHTTP(w, r)
	})

	return middleware.Chain(mux, middleware.Logger, middleware.CORS)
}
