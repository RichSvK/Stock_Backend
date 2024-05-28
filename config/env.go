package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func InitEnvironment() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env")
		return
	}
}
