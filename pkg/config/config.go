package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		User     string
		Password string
		Name     string
		Port     string
	}
	Secret struct {
		JwtSecretKey string
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Load Server Config
	config.Server.Port = os.Getenv("SERVER_PORT")

	// Load Database Config
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Name = os.Getenv("DB_NAME")

	// Secret
	config.Secret.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")

	return config, nil
}
