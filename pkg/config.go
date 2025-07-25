package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDiscordToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return os.Getenv("DISCORD_TOKEN")
}
