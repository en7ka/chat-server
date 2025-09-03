package config

import (
	"errors"
	"net"
	"os"
)

const (
	swaggerHostEnvName = "SWAGGER_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() (SwaggerConfig, error) {
	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("swagger host not found")
	}

	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}

	return &swaggerConfig{
		port: port,
		host: host,
	}, nil

}

func (c *swaggerConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
