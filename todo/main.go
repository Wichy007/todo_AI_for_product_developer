package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hello/todo/internal/handler"
	"hello/todo/internal/repository"
	"hello/todo/internal/service"
)

//go:embed web/static
var staticFiles embed.FS

func main() {
	port := envOrDefault("PORT", "8080")

	// Wire dependencies (Dependency-Inversion: each layer depends on interfaces)
	repo := repository.NewMemoryTaskRepository()
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	mux := http.NewServeMux()
	h.Register(mux)

	// Serve embedded static files for everything not matched by /api/
	static, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		log.Fatalf("static fs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(static)))

	srv := &http.Server{
		Addr:         "0.0.0.0:" + port, // bind all interfaces → reachable from LAN
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown (12-factor: disposability)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down…")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("Done.")
}

// envOrDefault reads an env var and falls back to a default value (12-factor config).
func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
