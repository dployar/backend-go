package main

import (
	"log"
	"log/slog"

	"github.com/dployar/v2/internal/handler"
	"github.com/dployar/v2/internal/logger"
	"github.com/dployar/v2/internal/middlewares"
	"github.com/dployar/v2/internal/server"
)

func main() {
	logger := logger.New(logger.DefaultConfig())

	// Aqui arranca el servidor
	srv := server.New(":8080", logger)

	healthHandler := handler.NewHealthHandler()

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	// Registrar las rutas sobre el router del servidor
	srv.Router().Get("/health", healthHandler.ServeHTTP)

	logger.Info("Servidor iniciando", slog.String("addr", ":8080"))
	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}
}
