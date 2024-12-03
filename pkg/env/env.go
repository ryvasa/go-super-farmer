package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		User     string
		Password string
		Name     string
		Port     string
		Timezone string
	}
	Secret struct {
		JwtSecretKey string
	}
}

func LoadEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	env := &Env{}

	// Load Server Config
	env.Server.Port = os.Getenv("SERVER_PORT")

	// Load Database Config
	env.Database.Host = os.Getenv("DB_HOST")
	env.Database.User = os.Getenv("DB_USER")
	env.Database.Port = os.Getenv("DB_PORT")
	env.Database.Password = os.Getenv("DB_PASSWORD")
	env.Database.Name = os.Getenv("DB_NAME")
	env.Database.Timezone = os.Getenv("DB_TIMEZONE")

	// Secret
	env.Secret.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")

	return env, nil
}
