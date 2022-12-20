package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql" // The database driver in use.
	"github.com/jmoiron/sqlx"
	"github.com/mathif92/auth-service/internal/config"
	"github.com/mathif92/auth-service/internal/handlers"
	"github.com/mathif92/auth-service/internal/logger"
	"github.com/mathif92/auth-service/internal/services"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	zapLogger := logger.NewZapLogger(zap.InfoLevel)
	defer zapLogger.Sync()

	health := handlers.NewHealth()

	cfg := NewConfig(zapLogger)
	db, err := OpenDB(cfg.DB)
	if err != nil {
		zapLogger.Error("Unexpected error", zap.String("error", errors.Wrap(err, "connecting to db").Error()))
		os.Exit(1)
	}

	tokenService := services.NewToken(cfg.SecretKey)
	authService := services.NewAuthentication(db, tokenService)
	authHandler := handlers.NewAuthenticationHandler(authService)

	router := chi.NewRouter()

	// All the middlewares used in the app
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", health.Ping)
	router.Post("/credentials", authHandler.CreateCredentials)
	router.Post("/auth", authHandler.ValidateCredentials)

	http.ListenAndServe(":8080", router)
}

func NewConfig(logger *zap.Logger) config.Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		logger.Fatal("CONFIG_FILE_PATH environment variable is not set")
	}

	config, err := config.NewConfig(configFilePath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error reading config file %s: %s", configFilePath, err.Error()))
	}

	logger.Info(fmt.Sprintf("Using config: %+v", config))

	return config
}

// OpenDB knows how to open a database connection based on the configuration.
func OpenDB(cfg config.DBConfig) (*sqlx.DB, error) {
	sslMode := "true"
	if cfg.DisableTLS {
		sslMode = "false"
	}

	q := make(url.Values)
	q.Set("tls", sslMode)
	q.Set("loc", "UTC")
	q.Set("parseTime", "true")
	q.Set("timeout", cfg.Timeout.String())
	q.Set("readTimeout", cfg.ReadTimeout.String())
	q.Set("writeTimeout", cfg.WriteTimeout.String())

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.DBName, q.Encode())

	db, err := sqlx.Open(cfg.Driver, connStr)
	if err != nil {
		fmt.Println("DB error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("DB PING error", err)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
