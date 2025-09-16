package db

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "stocky-backend/internal/config"
)

func NewPostgresConnection(cfg config.DatabaseConfig) (*sqlx.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(0)

    return db, nil
}


