//go:build local

package server

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Overload("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
