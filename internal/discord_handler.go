package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func StartDiscord() {
	dg, err := discordgo.New("Bot " + getDiscordToken())
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	channelID := "1396742940612886618"
	dg.ChannelMessageSend(channelID, "Odi is online")

	fmt.Println("Odi is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.ChannelMessageSend(channelID, "Bye bye")
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")

	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")

	case "ping dm":
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println("error creating channel:", err)
			s.ChannelMessageSend(m.ChannelID, "Something went wrong while sending the DM!")
			return
		}

		_, err = s.ChannelMessageSend(channel.ID, "Pong!")
		if err != nil {
			fmt.Println("error sending DM message:", err)
			s.ChannelMessageSend(m.ChannelID, "Failed to send you a DM. "+"Did you disable DM in your privacy settings?")
			return
		}

	case "embed":
		s.ChannelMessageSendEmbed(m.ChannelID, EmbedJoin())
	}
}

func getDiscordToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	return os.Getenv("DISCORD_TOKEN")
}
