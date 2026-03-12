package config

import (
	"github.com/joho/godotenv"
)

func InitEnvironment(fileName string) {
	_ = godotenv.Load(fileName)
}
