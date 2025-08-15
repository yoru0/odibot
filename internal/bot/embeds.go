package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/internal/store"
)

const (
	colorInfo    = 0x2F81F7 // blue
	colorSuccess = 0x31C48D // green
	colorWarn    = 0xF7B924 // yellow
	colorError   = 0xE02424 // red
)

func newEmbed(title, desc string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       color,
	}
}

func (b *Bot) sendEmbed(channelID, title, desc string, color int) {
	b.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{newEmbed(title, desc, color)},
	})
}

func (b *Bot) respondEmbed(s *discordgo.Session, ic *discordgo.InteractionCreate, title, desc string, color int) {
	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{newEmbed(title, desc, color)},
		},
	})
}

func (b *Bot) broadcastEmbed(sess *store.Session, title, desc string, color int) {
	seen := map[string]bool{}
	for _, chID := range sess.DMChannel {
		if seen[chID] {
			continue
		}
		seen[chID] = true
		b.sendEmbed(chID, title, desc, color)
	}
}
