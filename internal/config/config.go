package config

type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"/"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
}
