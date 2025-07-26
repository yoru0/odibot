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

func GetDiscordUserID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return os.Getenv("DISCORD_USER_ID")
}

func GetDiscordChannelID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return os.Getenv("DISCORD_CHANNEL_ID")
}
