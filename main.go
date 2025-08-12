package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/yoru0/odibot/internal/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN not set")
	}
	
	ownerID := os.Getenv("DISCORD_USER_ID")

	b, err := bot.New(token, ownerID)
	if err != nil {
		log.Fatalf("init bot: %v", err)
	}

	if err := b.Start(); err != nil {
		log.Fatalf("start bot: %v", err)
	}
	log.Println("Odi is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	b.Stop()
	log.Println("Stopped")
}
