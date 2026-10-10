package config

import "os"

type Config struct{ Host, Port string }

func Load() Config {
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8080"
	}
	return Config{Host: host, Port: port}
}

func (c Config) Address() string { return c.Host + ":" + c.Port }
