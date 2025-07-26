package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/pkg"
)

func StartDiscord() {
	dg, err := discordgo.New("Bot " + pkg.GetDiscordToken())
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

	}
}
