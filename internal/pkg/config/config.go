package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type RequestUUIDKey string
type LoggerKey string

const (
	OutputLogPath                        = "stdout logs.json"
	ErrorOutputLogPath                   = "stderr err_logs.json"
	RequestUUIDContextKey RequestUUIDKey = "requestUUID"
	LoggerContextKey      LoggerKey      = "logger"
)

type ServerConfig struct {
	Host    string   `yaml:"host"`
	Port    string   `yaml:"port"`
	Timeout int      `yaml:"timeout"`
	Origins []string `yaml:"origins"`
	Headers []string `yaml:"headers"`
	Methods []string `yaml:"methods"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func ReadConfig(cfgPath string) *Config {
	cfg := &Config{}

	file, err := os.Open(cfgPath)
	if err != nil {
		log.Println("Something went wrong while opening config file:", err)

		return nil
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		log.Println("Something went wrong while reading config file:", err)

		return nil
	}

	log.Println("Successfully opened config")

	return cfg
}
