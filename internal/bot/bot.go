package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "new api")
	}

	log.Printf("Authorized on account %s", api.Self.UserName)
	return &Bot{api: api}, nil
}

func (b *Bot) Run(_ context.Context, offset int) error {
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = 60 //todo

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return errors.Wrapf(err, "get updates chain")
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, _ = b.api.Send(msg) //todo validate and re-send
	}

	return nil
}

func (b *Bot) Debug(debug bool) {
	b.api.Debug = debug
}
