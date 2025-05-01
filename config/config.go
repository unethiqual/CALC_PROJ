package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL       string
    JWTSecret         string
    TimeAdditionMs    int
    TimeSubtractionMs int
    TimeMultiplicationMs int
    TimeDivisionMs    int
    ComputingPower    int
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    return &Config{
        DatabaseURL:       getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/calc_proj?sslmode=disable"),
        JWTSecret:         getEnv("JWT_SECRET", "your_jwt_secret"),
        TimeAdditionMs:    getEnvAsInt("TIME_ADDITION_MS", 100),
        TimeSubtractionMs: getEnvAsInt("TIME_SUBTRACTION_MS", 100),
        TimeMultiplicationMs: getEnvAsInt("TIME_MULTIPLICATIONS_MS", 200),
        TimeDivisionMs:    getEnvAsInt("TIME_DIVISIONS_MS", 200),
        ComputingPower:    getEnvAsInt("COMPUTING_POWER", 4),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}