package runtime

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yasufadhili/jawt/internal/core"
)

// DevServer provides a WebSocket server for live reloading
type DevServer struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger core.Logger

	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

// NewDevServer creates a new DevServer
func NewDevServer(ctx context.Context, logger core.Logger) *DevServer {
	serverCtx, cancel := context.WithCancel(ctx)
	return &DevServer{
		ctx:    serverCtx,
		cancel: cancel,
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start starts the dev server
func (s *DevServer) Start(addr string) error {
	http.HandleFunc("/ws", s.handleWebSocket)
	s.logger.Info("Starting dev server", core.StringField("address", addr))
	return http.ListenAndServe(addr, nil)
}

// Stop stops the dev server
func (s *DevServer) Stop() {
	s.cancel()
}

// Broadcast sends a message to all connected clients
func (s *DevServer) Broadcast(message []byte) {
	for client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			s.logger.Error("Failed to write message to client", core.ErrorField(err))
			delete(s.clients, client)
			client.Close()
		}
	}
}

func (s *DevServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("Failed to upgrade WebSocket connection", core.ErrorField(err))
		return
	}
	defer conn.Close()

	s.clients[conn] = true
	s.logger.Info("Client connected to dev server")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			s.logger.Error("Failed to read message from client", core.ErrorField(err))
			delete(s.clients, conn)
			break
		}
	}
}
