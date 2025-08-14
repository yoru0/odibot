package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/internal/game"
	"github.com/yoru0/odibot/internal/store"
)

const (
	customSelect = "capsa_select"
	customPlay   = "capsa_play"
	customSkip   = "capsa_skip"
	customClear  = "capsa_clear"
)

func (b *Bot) sendTurnUI(sess *store.Session) {
	actingID, actingName, isDummy := sess.Game.CurrentPlayerInfo()
	chID := sess.DMChannel[actingID]
	if isDummy && b.ownerID != "" {
		chID = sess.DMChannel[actingID]
	}

	if chID == "" {
		return
	}

	content := fmt.Sprintf("%s's turn. Select up to 5 cards, then press **Play**. Or press **Skip**.", actingName)

	comps := b.buildCardPickerComponents(sess, actingID, actingName)

	b.session.ChannelMessageSendComplex(chID, &discordgo.MessageSend{
		Content:    content,
		Components: comps,
	})
}

func (b *Bot) buildCardPickerComponents(sess *store.Session, actingUserID, actingName string) []discordgo.MessageComponent {
	hand := sess.Game.HandSnapshot(actingUserID)

	selected := make(map[string]bool, len(sess.Selected[actingUserID]))
	for _, v := range sess.Selected[actingUserID] {
		selected[strings.ToUpper(strings.TrimSpace(v))] = true
	}

	opts := make([]discordgo.SelectMenuOption, 0, len(hand))
	for _, c := range hand {
		val := strings.ToUpper(game.Card{Rank: c.Rank, Suit: c.Suit}.String())
		opts = append(opts, discordgo.SelectMenuOption{
			Label:   prettyLabel(c),
			Value:   val,
			Default: selected[val],
		})
	}

	min := 0
	menu := discordgo.SelectMenu{
		CustomID:    fmt.Sprintf("%s:%s:%s", customSelect, sess.LobbyChannelID, actingUserID),
		Placeholder: "Choose up to 5 cards",
		MinValues:   &min,
		MaxValues:   5,
		Options:     opts,
		MenuType:    discordgo.StringSelectMenu,
	}

	row1 := discordgo.ActionsRow{Components: []discordgo.MessageComponent{menu}}
	row2 := discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.Button{
			Style:    discordgo.PrimaryButton,
			Label:    "Play Selected",
			CustomID: fmt.Sprintf("%s:%s:%s", customPlay, sess.LobbyChannelID, actingUserID),
		},
		discordgo.Button{
			Style:    discordgo.SecondaryButton,
			Label:    "Skip",
			CustomID: fmt.Sprintf("%s:%s:%s", customSkip, sess.LobbyChannelID, actingUserID),
		},
		discordgo.Button{
			Style:    discordgo.SecondaryButton,
			Label:    "Clear",
			CustomID: fmt.Sprintf("%s:%s:%s", customClear, sess.LobbyChannelID, actingUserID),
		},
	}}
	return []discordgo.MessageComponent{row1, row2}
}

func joinPrettyCards(cards []game.Card) string {
	out := make([]string, 0, len(cards))
	for _, c := range cards {
		out = append(out, prettyLabel(c))
	}
	return strings.Join(out, " ")
}

func (b *Bot) buildAllHandsReport(sess *store.Session) string {
	players := sess.Game.PlayersSnapshot()
	if len(players) == 0 {
		return ""
	}

	maxName := 0
	for _, p := range players {
		if l := len(p.Name); l > maxName {
			maxName = l
		}
	}

	var sb strings.Builder
	sb.WriteString("All hands:\n")
	sb.WriteString("```\n")
	for _, p := range players {
		sb.WriteString(fmt.Sprintf("%-*s - %s\n", maxName, p.Name, joinPrettyCards(p.Hand)))
	}
	sb.WriteString("```")
	return sb.String()
}
