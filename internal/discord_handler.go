package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/pkg"
)

var (
	botToken  = pkg.GetDiscordToken()
	channelID = pkg.GetDiscordChannelID()
	ownerID   = pkg.GetDiscordUserID()
	quit      = make(chan struct{})
)

func StartDiscord() {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	if err := dg.Open(); err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	defer dg.Close()

	fmt.Println("Odi is running")
	dg.ChannelMessageSend(channelID, "Odi is now online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sc:
		fmt.Println("Shutting down via Ctrl+C")
	case <-quit:
		fmt.Println("Shutting down via Discord")
	}

	dg.ChannelMessageSend(channelID, "Odi is shutting down. Bye bye")
	fmt.Println("Bot shut down cleanly")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "capsa":
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Capsa on production",
			Reference: &discordgo.MessageReference{
				MessageID: m.ID,
				GuildID:   m.GuildID,
				ChannelID: m.ChannelID,
			},
		})

	case "dm":
		dmChannel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println("error creating DM channel:", err)
		}
		s.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
			Content: "Still on production",
		})

	case "embed":
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Here's a reply",
			Embed:   EmbedJoin(),
			Reference: &discordgo.MessageReference{
				MessageID: m.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			},
		})

	case "exit":
		if m.Author.ID == ownerID {
			go func() {
				time.Sleep(1 * time.Second)
				quit <- struct{}{}
			}()
		}
	}
}
