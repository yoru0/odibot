package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/internal/store"
)

const prefix = "!odi"

type Bot struct {
	session    *discordgo.Session
	manager    *store.Manager
	ownerID    string
	announceCh string
}

func New(token, ownerID, annouceCh string) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsMessageContent

	b := &Bot{
		session:    s,
		manager:    store.NewManager(),
		ownerID:    ownerID,
		announceCh: annouceCh,
	}
	s.AddHandler(b.onInteractionCreate)
	s.AddHandler(b.onMessageCreate)
	return b, nil
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return err
	}
	b.announce("Odi is running")
	return nil
}

func (b *Bot) Stop() error {
	b.announce("Bye bye")
	return b.session.Close()
}

func (b *Bot) onMessageCreate(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if m.GuildID != "" {
		b.routeGuild(m)
		return
	}
	b.routeDM(m)
}

// helpers ---

func (b *Bot) dm(userID, content string) {
	ch, err := b.session.UserChannelCreate(userID)
	if err != nil {
		return
	}
	b.session.ChannelMessageSend(ch.ID, content)
}

func (b *Bot) dmChannelID(userID string) (string, error) {
	ch, err := b.session.UserChannelCreate(userID)
	if err != nil {
		return "", err
	}
	return ch.ID, nil
}

func (b *Bot) broadcast(session *store.Session, content string) {
	seen := map[string]bool{}
	for _, chID := range session.DMChannel {
		if seen[chID] {
			continue
		}
		seen[chID] = true
		b.session.ChannelMessageSend(chID, content)
	}
}

func (b *Bot) announce(msg string) {
	if b.announceCh == "" {
		return
	}
	b.session.ChannelMessageSend(b.announceCh, msg)
}
