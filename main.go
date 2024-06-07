package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"backup/config"
	"backup/log"
	"backup/routes"
)

func main() {
	cfg, err := config.ReadConfig("config.json")
	logger := slog.New(
		log.NewFormatterHandler(
			log.TimezoneConverter(time.UTC),
			log.TimeFormatter(time.RFC3339, nil),
		)(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		),
	)
	logger = logger.With("env", "production")
	slog.SetDefault(logger)

	// httpmuxgo121 := godebug.New("httpmuxgo121")
	// slog.Info("", "v", httpmuxgo121.Value())

	if err != nil {
		port := flag.String("port", cfg.Server.Port, "IP address")
		flag.Parse()

		//User is expected to give :8080 like input, if they give 8080
		//we'll append the required ':'
		if !strings.HasPrefix(*port, ":") {
			cfg.Server.Port = ":" + *port
		} else {
			cfg.Server.Port = *port
		}
	}

	// Start server
	srv := &http.Server{
		Handler: routes.NewRouter(logger),
		Addr:    cfg.Server.Port,
		// Good practice: enforce timeouts for servers you create!
		ReadTimeout:  5 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
		IdleTimeout:  60 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	logger.Info("Server started", slog.String("port", cfg.Server.Port))

	// Block the rest of the code until a signal is received.
	sig := <-c
	slog.Info("Shutting everything down gracefully", "sig", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
	}
	slog.Info("Server shutdown successfully")
}
