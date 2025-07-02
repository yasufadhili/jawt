package server

import (
	"errors"
	"fmt"
	"github.com/yasufadhili/jawt/internal/project"
	"net/http"
	"os"
)

type DevServer struct {
	project *project.Project
	port    int
	host    string
	server  *http.Server
}

func NewDevServer(project *project.Project) *DevServer {
	port := project.Config.Server.Port
	if port == 0 {
		port = 6500
	}
	host := project.Config.Server.Host
	if host == "" {
		host = "localhost"
	}
	return &DevServer{
		project: project,
		port:    int(port),
		host:    host,
	}
}

func (s *DevServer) Start() error {
	fmt.Printf("📡 Development server running on http://%s\n", s.GetAddress())
	fmt.Println("   Press Ctrl+C to stop")

	if s.project.OutputDir == "" {
		return fmt.Errorf("no output directory specified")
	}

	if stat, err := os.Stat(s.project.OutputDir); os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Errorf("output directory %s does not exist", s.project.OutputDir)
	}

	fileServer := http.FileServer(http.Dir(s.project.OutputDir))

	http.Handle("/", fileServer)

	addr := fmt.Sprintf(":%d", s.port)

	s.server = &http.Server{Addr: addr, Handler: fileServer}
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), fileServer)
}

// Stop stops the development server
func (s *DevServer) Stop() error {
	if s.server != nil {
		fmt.Println("🛑 Stopping development server...")
		return s.server.Close()
	}
	return nil
}

func (s *DevServer) GetAddress() string {
	if s.server == nil {
		return fmt.Sprintf("%s:%d", s.host, s.port)
	}
	return s.server.Addr
}
