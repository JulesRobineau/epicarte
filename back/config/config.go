package config

import (
	"encoding/json"
	"gin-template/logging"
	"os"
)

type Environment string

const (
	Development Environment = "dev"
	Production  Environment = "prod"
)

type JwtConfig struct {
	// Secret is the secret key used to sign the JWT token
	Secret string `json:"secret" example:"secret"`
	// Expiration is the duration of the token in hours
	Expiration int `json:"expiration" example:"24" default:"24"`
}

type DatabaseConfig struct {
	// Host is the host of the database
	Host string `json:"host"`
	// Port is the port of the database
	Port int `json:"port"`
	// User is the user of the database
	Username string `json:"username"`
	// Password is the password of the database
	Password string `json:"password"`
	// Database is the name of the database
	DatabaseName string `json:"name"`
}

type ServerConfig struct {
	// Host is the host of the server
	Host string `json:"host"`
	// Port is the port of the server
	Port int `json:"port"`
}

type Config struct {
	// Env is the environment of the application
	Env Environment `json:"env"`
	// Version is the version of the application
	Version string `json:"version"`
	// Jwt is the configuration of the JWT
	Jwt JwtConfig `json:"jwt"`
	// Db is the configuration of the database
	Db DatabaseConfig `json:"database"`
	// Server is the configuration of the server
	Server ServerConfig `json:"server"`
}

// NewConfig returns a new configuration from a file path that is by default config.json
// If the file path is not provided, it will use the default config.json
// Config path is relative to the root of the project
// The file must be a json file
// The file path can be given by the environment variable CONFIG_PATH
func NewConfig() Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.json"
	}

	open, err := os.Open(path)
	if err != nil {
		return Config{
			Env:     Development,
			Version: "1.0.0",
			Jwt: JwtConfig{
				Secret:     "secret",
				Expiration: 24,
			},
			Db: DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				Username:     "postgres",
				Password:     "postgres",
				DatabaseName: "postgres",
			},
			Server: ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
		}
	}
	defer open.Close()

	decoder := json.NewDecoder(open)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		logging.Error.Fatal(err)
	}

	return configuration
}
