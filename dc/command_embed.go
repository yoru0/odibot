package dc

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func EmbedJoin() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Round [x]",
		Description: "[x]'s turn",
		URL:         "https://open.spotify.com/track/4xF4ZBGPZKxECeDFrqSAG4?si=2ab61dfe13ec4d56",
		Color:       0xCEDEBD,

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Inline Field 1.....................................................................",
				Value: "Value A",
			},
		},

		Footer: &discordgo.MessageEmbedFooter{
			Text: "...",
		},

		Timestamp: time.Now().Format(time.RFC3339),
	}
}
