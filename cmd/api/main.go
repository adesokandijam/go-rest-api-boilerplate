package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adesokandijam/go-rest-api-boilerplate/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn                        string
		maxOpenConns, maxIdleConns int
		maxIdleTime                time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "environment")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "database max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "database max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "database max idle time")
	flag.Parse()
	cfg.db.dsn = os.Getenv("DATABASE_URL")
	if cfg.db.dsn == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, err := openDB(cfg)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	app := &application{config: cfg, logger: logger, models: data.NewModels(db)}
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.port), Handler: app.routes(), IdleTimeout: time.Minute, ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError)}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownError <- server.Shutdown(ctx)
	}()
	logger.Info("starting server", "address", server.Addr, "environment", cfg.env)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
	if err := <-shutdownError; err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeExec
	db := stdlib.OpenDB(*connConfig)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
