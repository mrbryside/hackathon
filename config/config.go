package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Cfg *Config

// Config represents the application configuration
type Config struct {
	GitHubToken string
}

func LoadConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	Cfg = &Config{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
	}
}
