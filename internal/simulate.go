package internal

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SimulateLobby(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID
	channelID := m.ChannelID
	userID := m.Author.ID
	username := m.Author.Username

	lobby := &Lobby{
		GuildID:     guildID,
		ChannelID:   channelID,
		HostID:      userID,
		NumPlayers:  3,
		JoinedUsers: make(map[string]*Player),
	}

	lobby.JoinedUsers[userID] = &Player{
		ID:       userID,
		Username: username,
		Joined:   true,
	}

	lobby.JoinedUsers["dummy1"] = &Player{
		ID:       userID,
		Username: "Dummy1",
		Joined:   true,
	}
	lobby.JoinedUsers["dummy2"] = &Player{
		ID:       userID,
		Username: "Dummy2",
		Joined:   true,
	}

	Lobbies[guildID] = lobby

	s.ChannelMessageSend(channelID, "Simulated Capsa game starting with 3 players.")

	for fakeID, p := range lobby.JoinedUsers {
		dm, err := s.UserChannelCreate(p.ID)
		if err != nil {
			fmt.Println("DM error:", err)
			continue
		}
		s.ChannelMessageSend(dm.ID, fmt.Sprintf("[%s] Your cards: A♠ 2♥ 3♣", fakeID))
	}

	delete(Lobbies, guildID)
}
