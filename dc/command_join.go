package dc

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/capsa/player"
)

func HandleJoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	lobby, exists := Lobbies[m.GuildID]
	if !exists {
		fmt.Println(lobby, exists)
		return
	}

	// comment if not used
	if m.Author.ID == ownerID && len(lobby.JoinedUsers) == 1 {
		lobby.JoinedUsers["dummy1"] = &player.Player{UserID: "dummy1", Username: "Dummy1", Joined: true}
		lobby.JoinedUsers["dummy2"] = &player.Player{UserID: "dummy2", Username: "Dummy2", Joined: true}
		lobby.JoinedUsers["dummy3"] = &player.Player{UserID: "dummy3", Username: "Dummy3", Joined: true}
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Dev mode: Added 3 dummy players for simulation",
		})
	}

	// this also
	if _, ok := lobby.JoinedUsers[m.Author.ID]; ok && m.Author.ID != ownerID {
		s.ChannelMessageSend(m.ChannelID, "You already joined the game.")
		return
	}

	lobby.JoinedUsers[m.Author.ID] = &player.Player{
		UserID:   m.Author.ID,
		Username: m.Author.Username,
		Joined:   true,
	}

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

	for _, player := range lobby.JoinedUsers {
		dm, err := s.UserChannelCreate(m.Author.ID) // change to player.ID later
		if err != nil {
			fmt.Println("DM error:", err)
			continue
		}
		
		s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
			Content: fmt.Sprintf("On production, %s", player.Username),
			Embed:   EmbedJoin(),
		})
	}

	delete(Lobbies, m.GuildID)
}
