package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once         sync.Once
	instance     Configuration
	instantiated bool
)

func readConfig() Configuration {
	once.Do(func() {
		err := godotenv.Load("dev.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		instance = CreateDefaultConfiguration()
		instantiated = true
	})

	return instance
}

//Get configurations
func Get() Configuration {
	if !instantiated {
		return readConfig()
	}
	return instance
}
