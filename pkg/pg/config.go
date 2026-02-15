package pg

import (
    "fmt"
    "os"
)

type Config struct {
    TelegramToken string
    DBHost        string
    DBPort        string  
    DBUser        string
    DBPassword    string
    DBName        string
}

func Load() (*Config, error) {
    cfg := &Config{
        TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
        DBHost:        os.Getenv("DB_HOST"),
        DBPort:        os.Getenv("DB_PORT"),
        DBUser:        os.Getenv("DB_USER"),
        DBPassword:    os.Getenv("DB_PASSWORD"),
        DBName:        os.Getenv("DB_NAME"),
    }
    
    required := map[string]string{
        "TELEGRAM_BOT_TOKEN": cfg.TelegramToken,
        "DB_HOST":           cfg.DBHost,
        "DB_PORT":           cfg.DBPort,
        "DB_USER":           cfg.DBUser,
        "DB_PASSWORD":       cfg.DBPassword,
        "DB_NAME":           cfg.DBName,
    }
    
    for key, value := range required {
        if value == "" {
            return nil, fmt.Errorf("%s is required", key)
        }
    }
    
    return cfg, nil
}

func (c *Config) DatabaseURL() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
    )
}