package dc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/capsa/player"
)

func HandleCapsaCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if m.GuildID == "" {
	// 	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
	// 		Content: "Can't start game in a DM",
	// 	})
	// 	return
	// }

	args := strings.Fields(m.Content)
	if len(args) != 2 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Usage: `capsa <number_of_players>`",
		})
		return
	}

	numPlayers, err := strconv.Atoi(args[1])
	if err != nil || numPlayers < 2 || numPlayers > 4 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Enter a valid number (2-4)",
		})
		return
	}

	if _, exists := Lobbies[m.GuildID]; exists {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "A Capsa game is already waiting in this server",
		})
		return
	}

	Lobbies[m.GuildID] = &Lobby{
		GuildID:    m.GuildID,
		ChannelID:  m.ChannelID,
		HostID:     m.Author.ID,
		NumPlayers: numPlayers,
		JoinedUsers: map[string]*player.Player{
			m.Author.ID: {
				UserID:   m.Author.ID,
				Username: m.Author.Username,
				Joined:   true,
			},
		},
	}

	left := numPlayers - 1
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("%s started a Capsa game. Type `join` to participate. Waiting for %d more players.", m.Author.Username, left),
	})
}

