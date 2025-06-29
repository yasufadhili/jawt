package build

import (
	"fmt"
	"github.com/yasufadhili/jawt/internal/project"
	"net/http"
	"os"
)

type DevServer struct {
	project *project.Structure
	port    int
	server  *http.Server
}

func NewDevServer(p *project.Structure) *DevServer {
	port := p.Config.Server.Port
	if port == 0 {
		port = 6500
	}
	return &DevServer{
		project: p,
		port:    port,
	}
}

// Start starts the development server
func (ds *DevServer) Start() error {

	if ds.project.TempDir == "" {
		return fmt.Errorf("TempDir is not set")
	}

	if stat, err := os.Stat(ds.project.TempDir); os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Errorf("directory %q does not exist or is not a directory", ds.project.TempDir)
	}

	fileServer := http.FileServer(http.Dir(ds.project.TempDir))

	http.Handle("/", fileServer)

	addr := fmt.Sprintf(":%d", ds.port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop stops the development server
func (ds *DevServer) Stop() error {
	if ds.server != nil {
		fmt.Println("🛑 Stopping development server...")
		return ds.server.Close()
	}
	return nil
}

func (ds *DevServer) GetAddress() string {
	return ds.server.Addr
}

// handleRequest handles HTTP requests (placeholder)
func (ds *DevServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>JAWT Development Server</h1>")
	fmt.Fprintf(w, "<p>Project: %s</p>", ds.project.Config.Name)
	fmt.Fprintf(w, "<p>Pages: %d</p>", len(ds.project.Pages))
	fmt.Fprintf(w, "<p>Components: %d</p>", len(ds.project.Components))
	fmt.Fprintf(w, "<p>Assets: %d</p>", len(ds.project.Assets))

	fmt.Fprintf(w, "<h2>Routes:</h2><ul>")
	for _, page := range ds.project.Pages {
		fmt.Fprintf(w, "<li><a href=\"%s\">%s</a> - %s</li>", page.Route, page.Route, page.Title)
	}
	fmt.Fprintf(w, "</ul>")
}
