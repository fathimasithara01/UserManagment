package initalizeres

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEncVariable() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("errorr loading .env file")
	}
}
