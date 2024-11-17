package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func ReadConfig(cfgPath string, cfg interface{}) {
	file, err := os.Open(cfgPath)
	if err != nil {
		log.Println("Something went wrong while opening config file:", err)

		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		log.Println("Something went wrong while reading config file:", err)

		return
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Println("Something went wrong while reading config from .env file:", err)

		return
	}

	log.Println("Successfully opened config")

	return
}
