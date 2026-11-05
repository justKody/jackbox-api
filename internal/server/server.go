package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	v1 "github.com/justKody/jackbox-api/internal/server/v1"
)

type Server struct {
	ListenAddr string
	Router     *mux.Router
	httpServer *http.Server
}

// New builds the server with every route in place, so nothing is registered
// while it is already taking requests
func NewServer(listenAddr, assetDir string, handleWs http.HandlerFunc) *Server {
	router := mux.NewRouter()
	v1.HandleRoutes(router.PathPrefix("/v1").Subrouter())

	// the app builds these two from the socket url, so they stay outside /v1
	router.HandleFunc("/ws", handleWs).Methods(http.MethodGet)
	// question voice-overs, until they move to a bucket (see question.VoiceURL)
	router.PathPrefix("/assets/v2/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(assetDir))))

	return &Server{
		ListenAddr: listenAddr,
		Router:     router,
		httpServer: &http.Server{
			Addr:              listenAddr,
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second, // to prevent ddos slowloris attack
		},
	}
}

// Start blocks until the server stops, it always returns an error
// (http.ErrServerClosed after Shutdown)
func (s *Server) Start() error {
	fmt.Println("starting HTTP server on", s.ListenAddr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
