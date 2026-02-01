package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const DefaultPort = 8765

// Server manages the HTTP server for the log viewer
type Server struct {
	port       int
	webDir     string
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer(port int, webDir string) *Server {
	return &Server{
		port:   port,
		webDir: webDir,
	}
}

// Start begins serving the web viewer
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/v1/logs", withMiddleware(handleListLogs))
	mux.HandleFunc("GET /api/v1/logs/{id}", withMiddleware(handleGetLog))

	// Serve static files from web/dist
	distDir := filepath.Join(s.webDir, "dist")
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		return fmt.Errorf("web/dist directory not found. Run 'npm run build' in the web directory first")
	}

	// File server for static assets
	fs := http.FileServer(http.Dir(distDir))

	// SPA fallback handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := filepath.Join(distDir, r.URL.Path)
		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
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
