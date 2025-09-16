package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Redis    RedisConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

func Load() (*Config, error) {
    dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }

    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     dbPort,
            User:     getEnv("DB_USER", "stocky"),
            Password: getEnv("DB_PASSWORD", "password"),
            DBName:   getEnv("DB_NAME", "stocky_db"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Server: ServerConfig{
            Port:         serverPort,
            ReadTimeout:  30,
            WriteTimeout: 30,
        },
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}


