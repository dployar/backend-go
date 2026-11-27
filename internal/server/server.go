package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

// Encapusula al server http y sus dependencias
type Server struct {
	httpServer      *http.Server
	router          *chi.Mux
	shutdownTimeout time.Duration
	logger          *slog.Logger
}

// Func para poder crear un servidor nuevo
func New(addr string, logger *slog.Logger) *Server {
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
		logger:     logger,
	}
}

// Router expone el enrutador de Chi para registrar las rutas
// sobre el mismo handler que usa el servidor HTTP.
func (server *Server) Router() *chi.Mux {
	return server.router
}

func (server *Server) Start() error {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(quit)

	serveErr := make(chan error, 1)
	go func() {
		server.logger.Info("servidor iniciado", slog.String("addr", server.httpServer.Addr))
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		return err
	case sig := <-quit:
		server.logger.Info("señal recibida, iniciando shutdown",
			slog.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("error en shutdown", slog.Any("error", err))
		return err
	}

	server.logger.Info("Servidor detenido correctamente")

	return nil
}

func (server *Server) Use(middleware func(http.Handler) http.Handler) {
	server.httpServer.Handler = middleware(server.httpServer.Handler)
}
