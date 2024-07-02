package config

import (
	"fmt"
	"time"
)

type DSNConfig struct {
	UserName     string         // Username
	Password     string         // Password (requires UserName)
	DBName       string         // Database name
	Network      string         // Network type
	Address      string         // Network address (requires Net)
	ParseTime    bool           // Parse time values to time.Time
	TimeLocation *time.Location // Location for time.Time values
}

func NewDSNConfig(env Env) DSNConfig {
	net := "tcp"
	address := fmt.Sprintf("%s:%s", env.DBHost, env.DBPort)
	if env.Environment == "development" || env.Environment == "production" {
		net = "unix"
		address = fmt.Sprintf("/cloudsql/%s", env.DBHost)
	}

	location, _ := time.LoadLocation(env.TimeZone)

	return DSNConfig{
		UserName:     env.DBUsername,
		Password:     env.DBPassword,
		Network:      net,
		Address:      address,
		DBName:       env.DBName,
		ParseTime:    true,
		TimeLocation: location,
	}
}
