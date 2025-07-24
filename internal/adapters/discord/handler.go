package discord

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/goodi/internal/ports"
)

func StartDiscord(handler ports.CommandHandler, token string) error {
	dc, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	dc.AddHandler(NewMessageHandler(handler))

	err = dc.Open()
	if err != nil {
		return err
	}
	defer dc.Close()

	fmt.Println("Bot running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}

func NewMessageHandler(handler ports.CommandHandler) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "!add") {
			response := handler.Process(m.Content)
			s.ChannelMessageSend(m.ChannelID, response)
		}
	}
}
