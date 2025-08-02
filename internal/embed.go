package internal

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func EmbedJoin() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Well hello there",
		Description: "Testing the embed message",
		URL:         "https://open.spotify.com/track/4xF4ZBGPZKxECeDFrqSAG4?si=2ab61dfe13ec4d56",
		Color:       0x7289DA,

		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Discord Bot",
			URL:     "https://github.com/bwmarrin/discordgo",
			IconURL: "https://cdn-icons-png.flaticon.com/512/25/25231.png",
		},

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Inline Field 1",
				Value:  "Value A",
				Inline: true,
			},
			{
				Name:   "Inline Field 2",
				Value:  "Value B",
				Inline: true,
			},
			{
				Name:   "Non-inline Field",
				Value:  "Spans full width of embed",
				Inline: true,
			},
		},

		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Footer text here",
		},

		Timestamp: time.Now().Format(time.RFC3339),
	}
}
