package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yoru0/goodi/internal/adapters/discord"
	"github.com/yoru0/goodi/internal/adapters/postgres"
	"github.com/yoru0/goodi/internal/core/calculator"
)

func main() {
	godotenv.Load()

	store, err := postgres.NewSupabaseStore()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("Connected to DB")

	handler := &calculator.CalculatorBot{Store: store}
	token := os.Getenv("BOT_TOKEN")

	err = discord.StartDiscord(handler, token)
	if err != nil {
		log.Fatal("Discord bot failed:", err)
	}
}
