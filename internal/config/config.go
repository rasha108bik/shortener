package config

type Config struct {
	BaseURL       string `env:"BASE_URL" envDefault:"/"`
	ServerAddress string `env:"PORT" envDefault:"8080"`
}
