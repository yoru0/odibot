package bot

import (
	"fmt"
	"log"
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
		b.session.ChannelMessageSend(m.ChannelID, "Usage: `!odi <capsa|join|status|leave|dummy|shutdown>`")
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
		b.session.ChannelMessageSend(m.ChannelID, "Unknown command")
	}
}

func (b *Bot) handleShutdown(m *discordgo.MessageCreate) {
	if b.ownerID == "" || m.Author.ID != b.ownerID {
		return
	}
	b.Stop()
	log.Println("Stopped")
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
	sess := b.manager.Get(m.ChannelID)
	if sess == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here. Create one with `!odi capsa <numPlayers>`.")
		return
	}
	name := m.Author.Username
	if m.Member != nil && m.Member.Nick != "" {
		name = m.Member.Nick
	}
	if err := sess.Game.AddPlayer(m.Author.ID, name, m.Author.Username); err != nil {
		b.session.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	if dmID, err := b.dmChannelID(m.Author.ID); err == nil {
		sess.SetDMChannel(m.Author.ID, dmID)
	}

	left := sess.Desired - sess.Game.NumPlayers()
	if left > 0 {
		b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s joined. Waiting for %d more.", name, left))
		return
	}

	// Start game.
	if err := sess.Game.Start(); err != nil {
		b.session.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	b.manager.MarkStarted(m.ChannelID, b.ownerID)

	for _, player := range sess.Game.PlayersSnapshot() {
		b.dm(player.UserID, "["+player.Name+"] Your hand:\n"+player.HandString()+
			"\nUse: `play <cards>`, `skip`, `hand`, `table`, `quit`.")
	}

	b.session.ChannelMessageSend(m.ChannelID, "Game started in DMs. All further actions happen in private messages.")
	b.broadcastEmbed(sess, "Threes", "```\n"+sess.Game.FormatThreesReport()+"\n```", colorInfo)
	b.broadcastEmbed(sess, "Table", sess.Game.TableStateString(), colorInfo)
	b.sendTurnUI(sess)
}

func (b *Bot) handleStatus(m *discordgo.MessageCreate) {
	sess := b.manager.Get(m.ChannelID)
	if sess == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here")
		return
	}
	b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Players: %s (%d/%d).",
		sess.Game.PlayerList(), sess.Game.NumPlayers(), sess.Desired))
}

func (b *Bot) handleLeave(m *discordgo.MessageCreate) {
	sess := b.manager.Get(m.ChannelID)
	if sess == nil {
		b.session.ChannelMessageSend(m.ChannelID, "No lobby here")
		return
	}
	sess.Game.RemovePlayer(m.Author.ID)
	b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Left. Players: %s (%d/%d).",
		sess.Game.PlayerList(), sess.Game.NumPlayers(), sess.Desired))

	if sess.Game.NumPlayers() == 0 {
		b.manager.Delete(m.ChannelID)
		b.session.ChannelMessageSend(m.ChannelID, "Lobby removed.")
	}
}

func (b *Bot) handleDummy(m *discordgo.MessageCreate, args []string) {
	if m.Author.ID != b.ownerID {
		return
	}
	if len(args) < 3 {
		b.session.ChannelMessageSend(m.ChannelID, "Usage: `!odi dummy <1-4>`")
		return
	}
	n, err := strconv.Atoi(args[2])
	if err != nil || n < 1 || n > 4 {
		b.session.ChannelMessageSend(m.ChannelID, "Number must be 1-4.")
		return
	}
	sess := b.manager.Get(m.ChannelID)
	if sess == nil {
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
		if sess.Game.NumPlayers() >= sess.Desired {
			break
		}
		dummyID := fmt.Sprintf("dummy:%s:%d", m.ChannelID, sess.Game.NumPlayers())
		dummyName := fmt.Sprintf("Dummy%d", i+1)
		if err := sess.Game.AddDummy(dummyID, dummyName); err != nil {
			b.session.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		sess.SetDMChannel(dummyID, ownerDM)
		sess.HasDummy = true
		added++
	}

	left := sess.Desired - sess.Game.NumPlayers()
	if added == 0 {
		b.session.ChannelMessageSend(m.ChannelID, "No seats available for dummies.")
		return
	}
	if left > 0 {
		b.session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("Added %d dummy player(s). Waiting for %d more.", added, left))
		return
	}
	if err := sess.Game.Start(); err != nil {
		b.session.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	b.manager.MarkStarted(m.ChannelID, b.ownerID)
	for _, player := range sess.Game.PlayersSnapshot() {
		b.dm(player.UserID, "["+player.Name+"] Your hand:\n"+player.HandString()+
			"\nUse: `play <cards>`, `skip`, `hand`, `table`, `quit`.")
	}
	b.session.ChannelMessageSend(m.ChannelID, "Game started in DMs with dummies.")
	b.broadcastEmbed(sess, "Threes", "```\n"+sess.Game.FormatThreesReport()+"\n```", colorInfo)
	b.broadcastEmbed(sess, "Table", sess.Game.TableStateString(), colorInfo)
	b.sendTurnUI(sess)
}
