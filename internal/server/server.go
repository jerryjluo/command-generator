package server

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const DefaultPort = 8765

// Server manages the HTTP server for the log viewer
type Server struct {
	port       int
	assets     fs.FS
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer(port int, assets fs.FS) *Server {
	return &Server{
		port:   port,
		assets: assets,
	}
}

// Start begins serving the web viewer
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/v1/logs", withMiddleware(handleListLogs))
	mux.HandleFunc("GET /api/v1/logs/{id}", withMiddleware(handleGetLog))

	// File server for static assets (from embedded filesystem)
	fileServer := http.FileServerFS(s.assets)

	// SPA fallback handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		f, err := s.assets.Open(path)
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		index, err := s.assets.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer index.Close()
		content, _ := io.ReadAll(index)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
	})

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// URL returns the server's URL
func (s *Server) URL() string {
	return fmt.Sprintf("http://localhost:%d", s.port)
}

// OpenBrowser opens the default browser to the given URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
