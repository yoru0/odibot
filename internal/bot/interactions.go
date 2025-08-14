package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/internal/store"
)

func (b *Bot) onInteractionCreate(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	if ic.Type != discordgo.InteractionMessageComponent {
		return
	}
	data := ic.MessageComponentData()
	custom := data.CustomID

	parts := strings.SplitN(custom, ":", 3)
	if len(parts) != 3 {
		return
	}
	kind, lobbyID, actingID := parts[0], parts[1], parts[2]
	sess := b.manager.Get(lobbyID)

	if sess == nil || !sess.Started {
		_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Session is no longer active."},
		})
		return
	}

	switch kind {
	case customSelect:
		b.handleCustomSelect(s, ic, sess, actingID, data.Values)

	case customClear:
		b.handleCustomClear(s, ic, sess, actingID)

	case customSkip:
		b.handleCustomSkip(s, ic, sess, actingID)

	case customPlay:
		b.handleCustomPlay(s, ic, sess, actingID)
	}
}

func (b *Bot) handleCustomSelect(s *discordgo.Session, ic *discordgo.InteractionCreate, sess *store.Session, actingID string, values []string) {
	if sess.Selected == nil {
		sess.Selected = make(map[string][]string)
	}
	sess.Selected[actingID] = append([]string(nil), values...)
	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (b *Bot) handleCustomClear(s *discordgo.Session, ic *discordgo.InteractionCreate, sess *store.Session, actingID string) {
	if sess.Selected != nil {
		delete(sess.Selected, actingID)
	}
	b.respondWithMessage(s, ic, "Selection cleared.")
}

func (b *Bot) handleCustomSkip(s *discordgo.Session, ic *discordgo.InteractionCreate, sess *store.Session, actingID string) {
	clicker := b.interactionUserID(ic)
	if clicker == "" {
		b.respondWithError(s, ic, "Unknown clicker.")
		return
	}
	currID, _, currIsDummy := sess.Game.CurrentPlayerInfo()
	if clicker == b.ownerID && currIsDummy {
		actingID = currID
	}
	if clicker != actingID && clicker != b.ownerID {
		b.respondWithError(s, ic, "Not your turn.")
		return
	}

	b.deferResponse(s, ic)

	msg, err := sess.Game.Skip(actingID)
	if err != nil {
		s.ChannelMessageSend(ic.ChannelID, err.Error())
		return
	}

	// Disable the components.
	b.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    ic.ChannelID,
		ID:         ic.Message.ID,
		Components: &[]discordgo.MessageComponent{},
	})

	b.broadcast(sess, msg)
	b.sendTurnUI(sess)
}

func (b *Bot) handleCustomPlay(s *discordgo.Session, ic *discordgo.InteractionCreate, sess *store.Session, actingID string) {
	clicker := b.interactionUserID(ic)
	if clicker == "" {
		b.respondWithError(s, ic, "Unknown clicker.")
		return
	}
	currID, _, currIsDummy := sess.Game.CurrentPlayerInfo()
	if clicker == b.ownerID && currIsDummy {
		actingID = currID
	}
	if clicker != actingID && clicker != b.ownerID {
		b.respondWithError(s, ic, "Not your turn.")
		return
	}

	values := sess.Selected[actingID]
	if len(values) == 0 {
		b.respondWithError(s, ic, "Select 1-5 cards first.")
		return
	}
	b.deferResponse(s, ic)

	msg, err := sess.Game.Play(actingID, values)
	if err != nil {
		s.ChannelMessageSend(ic.ChannelID, err.Error())
		return
	}

	// Clear selection after a successful play.
	delete(sess.Selected, actingID)

	// Disable the components.
	b.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    ic.ChannelID,
		ID:         ic.Message.ID,
		Components: &[]discordgo.MessageComponent{},
	})

	b.broadcast(sess, msg)
	if sess.Game.IsOver() {
		b.broadcast(sess, "Game over. Standings:\n"+sess.Game.ResultsString())
		b.manager.Delete(sess.LobbyChannelID)
		return
	}

	b.sendTurnUI(sess)

}

// * Helper
func (b *Bot) respondWithError(s *discordgo.Session, ic *discordgo.InteractionCreate, message string) {
	b.respondWithMessage(s, ic, message)
}

func (b *Bot) respondWithMessage(s *discordgo.Session, ic *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: message},
	})
}

func (b *Bot) deferResponse(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (b *Bot) interactionUserID(ic *discordgo.InteractionCreate) string {
	if ic.Member != nil && ic.Member.User != nil {
		return ic.Member.User.ID
	}
	if ic.User != nil {
		return ic.User.ID
	}
	return ""
}
