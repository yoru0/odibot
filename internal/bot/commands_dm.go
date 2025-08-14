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
	session := b.manager.GetByUser(m.Author.ID)
	if session == nil || !session.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}

	currID, currName, currIsDummy := session.Game.CurrentPlayerInfo()
	targetID := m.Author.ID
	title := "Your hand:\n"
	if m.Author.ID == b.ownerID && currIsDummy {
		targetID = currID
		title = "[" + currName + "] hand:\n"
	}

	player := session.Game.FindPlayer(targetID)
	if player == nil {
		b.dm(m.Author.ID, "You are not in this game")
		return
	}

	b.dm(m.Author.ID, title+player.HandString())
}

func (b *Bot) dmHandlePlay(m *discordgo.MessageCreate, tail string) {
	session := b.manager.GetByUser(m.Author.ID)
	if session == nil || !session.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}

	codes := strings.Fields(tail)

	actingID := m.Author.ID
	currID, _, currIsDummy := session.Game.CurrentPlayerInfo()
	if m.Author.ID == b.ownerID && currIsDummy {
		actingID = currID
	}

	msg, err := session.Game.Play(actingID, codes)
	if err != nil {
		b.dm(m.Author.ID, err.Error())
		return
	}

	b.broadcast(session, msg)
	b.sendTurnUI(session)
	if session.Game.IsOver() {
		b.broadcast(session, "Game over. Standings:\n"+session.Game.ResultsString())
		b.manager.Delete(session.LobbyChannelID)
	}
}

func (b *Bot) dmHandleSkip(m *discordgo.MessageCreate) {
	session := b.manager.GetByUser(m.Author.ID)
	if session == nil || !session.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}
	actingID := m.Author.ID
	currID, _, currIsDummy := session.Game.CurrentPlayerInfo()
	if m.Author.ID == b.ownerID && currIsDummy {
		actingID = currID
	}
	msg, err := session.Game.Skip(actingID)
	if err != nil {
		b.dm(m.Author.ID, err.Error())
		return
	}
	b.broadcast(session, msg)
	b.sendTurnUI(session)
}

func (b *Bot) dmHandleTable(m *discordgo.MessageCreate) {
	session := b.manager.GetByUser(m.Author.ID)
	if session == nil || !session.Started {
		b.dm(m.Author.ID, "You are not in an active game.")
		return
	}
	b.dm(m.Author.ID, session.Game.TableStateString())
}

func (b *Bot) dmHandleQuit(m *discordgo.MessageCreate) {
	session := b.manager.GetByUser(m.Author.ID)
	if session == nil {
		b.dm(m.Author.ID, "No active session")
		return
	}
	session.Game.RemovePlayer(m.Author.ID)
	b.broadcast(session, m.Author.Username+" quit the game.")
	if session.Game.NumPlayers() == 0 {
		b.manager.Delete(session.LobbyChannelID)
	} else if session.Game.IsOver() {
		b.manager.Delete(session.LobbyChannelID)
	}
}
