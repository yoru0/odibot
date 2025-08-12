package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yoru0/odibot/internal/store"
)

const prefix = "!odi"

type Bot struct {
	session *discordgo.Session
	manager *store.Manager
	ownerID string
}

func New(token, ownerID string) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsMessageContent

	b := &Bot{
		session: s,
		manager: store.NewManager(),
		ownerID: ownerID,
	}
	s.AddHandler(b.onMessageCreate)
	return b, nil
}

func (b *Bot) Start() error {
	return b.session.Open()
}

func (b *Bot) Stop() error {
	return b.session.Close()
}

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
	for userID := range session.DMChannel {
		b.session.ChannelMessageSend(session.DMChannel[userID], content)
	}
}

func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if m.GuildID != "" {
		b.routeGuild(m)
		return
	}
	b.routeDM(m)
}
