package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Encapusula al server http y sus dependencias
type Server struct {
	httpServer      *http.Server
	router          *chi.Mux
	shutdownTimeout time.Duration
}

// Func para poder crear un servidor nuevo
func New(addr string) *Server {
	// Enrutador http
	r := chi.NewRouter()

	// Configurar servidor HTTP
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	return &Server{
		httpServer: httpServer,
		router:     r,
	}
}

// Router expone el enrutador de Chi para registrar las rutas
// sobre el mismo handler que usa el servidor HTTP.
func (server *Server) Router() *chi.Mux {
	return server.router
}

func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

func (server *Server) Use(middleware func(http.Handler) http.Handler) {
	server.httpServer.Handler = middleware(server.httpServer.Handler)
}
