package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) routeDM(m *discordgo.MessageCreate) {
	args := strings.TrimSpace(strings.ToLower(m.Content))

	switch {
	case strings.HasPrefix(args, "hand"):
		b.dmHandleHand(m)

	case strings.HasPrefix(args, "play "):
		b.dmHandlePlay(m, strings.TrimSpace(args[5:]))

	case strings.HasPrefix(args, "skip"):
		b.dmHandleSkip(m)

	case strings.HasPrefix(args, "table"):
		b.dmHandleTable(m)

	case strings.HasPrefix(args, "quit"):
		b.dmHandleQuit(m)
	}
}

func (b *Bot) dmHandleHand(m *discordgo.MessageCreate) {
	sess := b.manager.GetByUser(m.Author.ID)
	if sess == nil || !sess.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}

	currID, currName, currIsDummy := sess.Game.CurrentPlayerInfo()
	targetID := m.Author.ID
	title := "Your hand:\n"
	if m.Author.ID == b.ownerID && currIsDummy {
		targetID = currID
		title = "[" + currName + "] hand:\n"
	}

	player := sess.Game.FindPlayer(targetID)
	if player == nil {
		b.dm(m.Author.ID, "You are not in this game")
		return
	}

	b.dm(m.Author.ID, title+player.HandString())
}

func (b *Bot) dmHandlePlay(m *discordgo.MessageCreate, tail string) {
	sess := b.manager.GetByUser(m.Author.ID)
	if sess == nil || !sess.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}
	codes := strings.Fields(tail)

	actingID := m.Author.ID
	currID, _, currIsDummy := sess.Game.CurrentPlayerInfo()
	if m.Author.ID == b.ownerID && currIsDummy {
		actingID = currID
	}

	msg, err := sess.Game.Play(actingID, codes)
	if err != nil {
		b.dm(m.Author.ID, err.Error())
		return
	}
	b.broadcast(sess, msg)
	b.sendTurnUI(sess)
	if sess.Game.IsOver() {
		b.broadcast(sess, "Game over. Standings:\n"+sess.Game.ResultsString())
		b.manager.Delete(sess.LobbyChannelID)
	}
}

func (b *Bot) dmHandleSkip(m *discordgo.MessageCreate) {
	sess := b.manager.GetByUser(m.Author.ID)
	if sess == nil || !sess.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}
	actingID := m.Author.ID
	currID, _, currIsDummy := sess.Game.CurrentPlayerInfo()
	if m.Author.ID == b.ownerID && currIsDummy {
		actingID = currID
	}
	msg, err := sess.Game.Skip(actingID)
	if err != nil {
		b.dm(m.Author.ID, err.Error())
		return
	}
	b.broadcast(sess, msg)
	b.sendTurnUI(sess)
}

func (b *Bot) dmHandleTable(m *discordgo.MessageCreate) {
	sess := b.manager.GetByUser(m.Author.ID)
	if sess == nil || !sess.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}
	b.dm(m.Author.ID, sess.Game.TableStateString())
}

func (b *Bot) dmHandleQuit(m *discordgo.MessageCreate) {
	sess := b.manager.GetByUser(m.Author.ID)
	if sess == nil {
		b.dm(m.Author.ID, "No active session")
		return
	}
	sess.Game.RemovePlayer(m.Author.ID)
	b.broadcast(sess, m.Author.Username+" quit the game.")
	if sess.Game.NumPlayers() == 0 || sess.Game.IsOver() {
		b.manager.Delete(sess.LobbyChannelID)
	}
}
