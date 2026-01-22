package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divyanshujswl-zs/students-api/internal/config"
	"github.com/divyanshujswl-zs/students-api/internal/http/handlers/student"
	"github.com/divyanshujswl-zs/students-api/internal/storage"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// db setup
	storage, err := storage.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info(
		"storage initialised",
		slog.String("env", cfg.Env),
		slog.String("version", "1.0.0"),
		slog.String("connected_db", cfg.DB.Driver+":"+cfg.DB.Name),
	)

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students API!"))
	})

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetUserById(storage))
	router.HandleFunc("GET /api/students", student.GetUserList(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started ", slog.String("config address:", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Server failed to start %s", err.Error())
		}
	}()

	<-done

	// gracefully shutdown
	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server ", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")
}
