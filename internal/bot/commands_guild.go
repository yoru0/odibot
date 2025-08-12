package bot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) routeGuild(m *discordgo.MessageCreate) {
	if !strings.HasPrefix(strings.ToLower(m.Content), prefix) {
		return
	}
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		b.session.ChannelMessageSend(m.ChannelID, "valid command pls")
		return
	}

	switch strings.ToLower(args[1]) {
	case "shutdown":
		b.handleShutdown(m)

	case "capsa":
		b.handleCapsa(m, args)

	case "join":
		b.handleJoin(m)

	case "status":
		b.handleStatus(m)

	case "leave":
		b.handleLeave(m)

	case "dummy":
		b.handleDummy(m, args)

	default:
		b.session.ChannelMessageSend(m.ChannelID, "Unknown command.")
	}
}

func (b *Bot) handleShutdown(m *discordgo.MessageCreate) {
	if b.ownerID == "" || m.Author.ID != b.ownerID {
		return
	}
	b.Stop()
	os.Exit(0)
}

func (b *Bot) handleCapsa(m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		b.session.ChannelMessageSend(m.ChannelID, "Usage: `!odi capsa <3|4>`")
		return
	}
	n, err := strconv.Atoi(args[2])
	if err != nil || n < 3 || n > 4 {
		b.session.ChannelMessageSend(m.ChannelID, "Number of players must be 3 or 4")
		return
	}
	if b.manager.Has(m.ChannelID) {
		b.session.ChannelMessageSend(m.ChannelID, "A lobby already exists in this channel.")
		return
	}
	b.manager.NewSession(m.ChannelID, n)
	b.session.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf("Capsa lobby created for %d players. Others join with `!odi join`.", n))
}

func (b *Bot) handleJoin(m *discordgo.MessageCreate) {
	session := b.manager.Get(m.ChannelID)
	if session == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here. Create one with `!odi capsa <numPlayers>`.")
		return
	}
	name := m.Author.Username
	if m.Member != nil && m.Member.Nick != "" {
		name = m.Member.Nick
	}
	if err := session.Game.AddPlayer(m.Author.ID, name, m.Author.Username); err != nil {
		b.session.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	dmID, err := b.dmChannelID(m.Author.ID)
	if err == nil {
		session.DMChannel[m.Author.ID] = dmID
	}
	left := session.Desired - session.Game.NumPlayers()
	if left > 0 {
		b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s joined. Waiting for %d more.", name, left))
		return
	}

	// Start game.
	if err := session.Game.Start(); err != nil {
		b.session.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	b.manager.MarkStarted(m.ChannelID, b.ownerID)
	for _, player := range session.Game.PlayersSnapshot() {
		b.dm(player.UserID, "["+player.Name+"] Your hand:\n"+player.HandString()+
			"\nUse: `play <cards>`, `skip`, `hand`, `table`, `quit`.")
	}
	b.broadcast(session, session.Game.TableStateString())
	b.session.ChannelMessageSend(m.ChannelID, "Game started in DMs. All further actions happen in private messages.")
}

func (b *Bot) handleStatus(m *discordgo.MessageCreate) {
	session := b.manager.Get(m.ChannelID)
	if session == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here")
		return
	}
	b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Players: %s (%d/%d).",
		session.Game.PlayerList(), session.Game.NumPlayers(), session.Desired))
}

func (b *Bot) handleLeave(m *discordgo.MessageCreate) {
	session := b.manager.Get(m.ChannelID)
	if session == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here")
		return
	}
	session.Game.RemovePlayer(m.Author.ID)
	b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Left. Players: %s (%d/%d).",
		session.Game.PlayerList(), session.Game.NumPlayers(), session.Desired))

	if session.Game.NumPlayers() == 0 {
		b.manager.Delete(m.ChannelID)
		b.session.ChannelMessageSend(m.ChannelID, "Lobby removed.")
	}
}

func (b *Bot) handleDummy(m *discordgo.MessageCreate, args []string) {
	if m.Author.ID != b.ownerID {
		return
	}
	if len(args) < 3 {
		b.session.ChannelMessageSend(m.ChannelID, "Usage: `!odi dummy <2|3>`")
		return
	}
	n, err := strconv.Atoi(args[2])
	if err != nil || n < 1 || n > 3 {
		b.session.ChannelMessageSend(m.ChannelID, "Number must be 2-3.")
		return
	}
	session := b.manager.Get(m.ChannelID)
	if session == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here. Create one with `!odi capsa <3|4>`.")
		return
	}
	ownerDM, err := b.dmChannelID(b.ownerID)
	if err != nil {
		b.session.ChannelMessageSend(m.ChannelID, "Cannot create owner DM.")
		return
	}

	added := 0
	for i := 0; i < n; i++ {
		if session.Game.NumPlayers() >= session.Desired {
			break
		}
		dummyID := fmt.Sprintf("dummy:%s:%d", m.ChannelID, session.Game.NumPlayers())
		dummyName := fmt.Sprintf("Dummy%d", i+1)
		if err := session.Game.AddDummy(dummyID, dummyName); err != nil {
			b.session.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		session.DMChannel[dummyID] = ownerDM
		session.HasDummy = true
		added++
	}

	left := session.Desired - session.Game.NumPlayers()
	if added == 0 {
		b.session.ChannelMessageSend(m.ChannelID, "No seats available for dummies.")
		return
	}
	if left > 0 {
		b.session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("Added %d dummy player(s). Waiting for %d more.", added, left))
		return
	}
	if err := session.Game.Start(); err != nil {
		b.session.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	b.manager.MarkStarted(m.ChannelID, b.ownerID)
	for _, player := range session.Game.PlayersSnapshot() {
		b.dm(player.UserID, "["+player.Name+"] Your hand:\n"+player.HandString()+
			"\nUse: `play <cards>`, `skip`, `hand`, `table`, `quit`.")
	}
	b.broadcast(session, session.Game.TableStateString())
	b.session.ChannelMessageSend(m.ChannelID, "Game started in DMs with dummies.")
}
