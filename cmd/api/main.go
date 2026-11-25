package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/dployar/v2/internal/handler"
	"github.com/dployar/v2/internal/middlewares"
	"github.com/dployar/v2/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	srv := server.New(":8080")

	healthHandler := handler.NewHealthHandler()

	srv.Use(middlewares.Logger(logger))

	// Registrar las rutas sobre el router del servidor
	srv.Router().Get("/health", healthHandler.ServeHTTP)

	logger.Info("Servidor iniciando", slog.String("addr", ":8080"))
	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}
}
