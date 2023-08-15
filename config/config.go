package config

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config include all parametrs for run server
// SERVER_ADDRESS, BASE_URL, FILE_STORAGE_PATH, DATABASE_DSN.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	EnableHTTPS     string `env:"ENABLE_HTTPS" envDefault:"" json:"enable_https"`

	FileName string `env:"CONFIG" envDefault:""`
}

var (
	serverAddress     string
	baseURL           string
	fileStoragePath   string
	databaseDSN       string
	listenAndServeTLS string

	fileName string
)

// NewConfig returns a newly initialized Config objects.
func NewConfig() *Config {
	flag.StringVar(&serverAddress, "a", "", "server address")
	flag.StringVar(&baseURL, "b", "", "base URL")
	flag.StringVar(&fileStoragePath, "f", "", "file storage path")
	flag.StringVar(&databaseDSN, "d", "", "databaseDSN path")
	flag.StringVar(&listenAndServeTLS, "s", "", "start ListenAndServeTLS")

	flag.StringVar(&fileName, "c", "", "config file name")
	flag.StringVar(&fileName, "config", "", "config file name")

	flag.Parse()
	log.Printf("server address: %s, base URL: %s, file storagePath: %s databaseDSN: %s\n", serverAddress, baseURL, fileStoragePath, databaseDSN)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if serverAddress != "" {
		cfg.ServerAddress = serverAddress
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}

	cfg.FileStoragePath = fileStoragePath
	cfg.DatabaseDSN = databaseDSN

	if listenAndServeTLS != "" {
		cfg.EnableHTTPS = listenAndServeTLS
	}

	if fileName != "" {
		cfg.FileName = fileName
	}

	if cfg.FileName != "" {
		conf, err := ReadJSONFile(cfg.FileName)
		if err != nil {
			log.Println("ReadJSONFile failed")
		}

		patchData(&cfg, conf)
	}

	return &cfg
}

// ReadJSONFile read json file by name file
func ReadJSONFile(dir string) (Config, error) {
	var cfg Config

	jsonFile, err := os.Open(dir)
	if err != nil {
		return cfg, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	result := &Config{}
	err = json.Unmarshal([]byte(byteValue), result)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func patchData(cfg *Config, cfgFromFileJSON Config) {
	if cfg.ServerAddress != "" {
		cfg.ServerAddress = cfgFromFileJSON.ServerAddress
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = cfgFromFileJSON.BaseURL
	}

	if cfg.DatabaseDSN != "" {
		cfg.DatabaseDSN = cfgFromFileJSON.DatabaseDSN
	}

	if cfg.FileStoragePath != "" {
		cfg.FileStoragePath = cfgFromFileJSON.FileStoragePath
	}

	if cfg.EnableHTTPS != "" {
		cfg.EnableHTTPS = cfgFromFileJSON.EnableHTTPS
	}
}
