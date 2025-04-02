package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig get variables from .env and load
func LoadConfig() (*Config, error) {
	// use default values as setup
	setDefaults()

	// read .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("cannot read .env file: %w", err)
	}

	// check and use env variables
	viper.AutomaticEnv()

	return &Config{
		Server: ServerConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetInt("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		}}, nil

}

// setDefaults set default env values
func setDefaults() {
	// server default setup
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("SERVER_PORT", 8080)

	// database default setup
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 26260)
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "checkin")
	viper.SetDefault("DB_SSL_MODE", "disable")
}

// GetServerAddress get server host address
func (c *Config) GetServerAddress() string {
	serverAddress := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Host)
	return serverAddress
}

// GetDatabaseConnectionString get database connection string with or without password
func (c *Config) GetDatabaseConnectionString() string {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.DBName)

	if c.Database.Password != "" {
		connectionString += fmt.Sprintf(" password=%s", c.Database.Password)
	}

	connectionString += fmt.Sprintf(" sslmode=%s", c.Database.SSLMode)

	return connectionString
}
