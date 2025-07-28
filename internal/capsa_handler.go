package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleCapsaCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) != 2 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Usage: `!capsa <number_of_players>`",
		})
		return
	}

	numPlayers, err := strconv.Atoi(args[1])
	if err != nil || numPlayers < 3 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Enter a valid number (minimal 3)",
		})
		return
	}

	if _, exists := Lobbies[m.GuildID]; exists {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "A Capsa game is already waiting in this server",
		})
		return
	}

	Lobbies[m.ChannelID] = &Lobby{
		GuildID:     m.GuildID,
		ChannelID:   m.ChannelID,
		HostID:      m.Author.ID,
		NumPlayers:  numPlayers,
		JoinedUsers: map[string]bool{m.Author.ID: true},
	}

	left := numPlayers - 1
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("%s started a Capsa game. Type `join` to participate. Waiting for %d more players.", m.Author.Username, left),
	})
}

func HandleJoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	lobby, exists := Lobbies[m.GuildID]
	if !exists {
		return
	}

	// if lobby.JoinedUsers[m.Author.ID] {
	// 	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
	// 		Content: "You already joined the game",
	// 	})
	// 	return
	// }

	// lobby.JoinedUsers[m.Author.ID] = true
	left := lobby.NumPlayers - len(lobby.JoinedUsers)
	if left > 0 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("%s joined! Waiting for %d more...", m.Author.Username, left),
		})
		return
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "All players joined! Starting game...",
	})

	for userID := range lobby.JoinedUsers {
		dm, _ := s.UserChannelCreate(userID)
		s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
			Content: "Capsa game started",
		})

		// TODO Assign Card
	}

	delete(Lobbies, m.GuildID)

}
