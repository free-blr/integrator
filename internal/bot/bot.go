package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"

	"free.blr/integrator/internal/model"
)

type commandHandler = func(context.Context, *tgbotapi.Message) error

type Bot struct {
	api            *tgbotapi.BotAPI
	tagService     tagService
	requestService requestService
}

func NewBot(token string, tagService tagService, requestService requestService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "new api")
	}

	log.Printf("Authorized on account %s", api.Self.UserName)
	return &Bot{
		api:            api,
		tagService:     tagService,
		requestService: requestService,
	}, nil
}

func (b *Bot) Run(ctx context.Context, offset int) error {
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = 60 //todo

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return errors.Wrapf(err, "get updates chain")
	}

	for msg := range updates {
		if msg.CallbackQuery != nil {
			if err = b.handleCallbackQuery(ctx, msg.CallbackQuery); err != nil {
				log.Println("ERROR", err)
			}
		}

		if msg.Message != nil {
			if err = b.handleMessage(ctx, msg.Message); err != nil {
				log.Println("ERROR", err)
			}
		}
	}

	return nil
}

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		return nil
	}
	return b.resolveHandler(msg.Command())(ctx, msg)
}

func (b *Bot) resolveHandler(command string) commandHandler {
	switch command {
	case CommandHelp:
		return b.handleHelp
	case CommandAskHelp:
		return b.handleAskHelp
	case CommandOfferHelp:
		return b.handleOfferHelp
	case CommandGetAskHelp:
		return b.handleGetAskHelp
	case CommandGetOfferHelp:
		return b.handleGetOfferHelp
	default:
		return b.handleUnknown
	}
}

func (b *Bot) handleHelp(_ context.Context, msg *tgbotapi.Message) error {
	rsp := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Введите /%s или /%s.", CommandAskHelp, CommandOfferHelp))
	_, err := b.api.Send(rsp)
	return err
}

func (b *Bot) handleAskHelp(ctx context.Context, msg *tgbotapi.Message) error {
	selectors, err := b.tagSelector(ctx, CommandAskHelp)
	if err != nil {
		return err
	}

	rsp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите раздел")
	rsp.ReplyMarkup = selectors
	_, err = b.api.Send(rsp)
	return err
}

func (b *Bot) handleOfferHelp(ctx context.Context, msg *tgbotapi.Message) error {
	selectors, err := b.tagSelector(ctx, CommandOfferHelp)
	if err != nil {
		return errors.Wrap(err, "tag selector")
	}

	rsp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите раздел")
	rsp.ReplyMarkup = selectors
	_, err = b.api.Send(rsp)
	return err
}

func (b *Bot) handleGetAskHelp(ctx context.Context, msg *tgbotapi.Message) error {
	selectors, err := b.tagSelector(ctx, CommandGetAskHelp)
	if err != nil {
		return err
	}

	rsp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите раздел")
	rsp.ReplyMarkup = selectors
	_, err = b.api.Send(rsp)
	return err
}

func (b *Bot) handleGetOfferHelp(ctx context.Context, msg *tgbotapi.Message) error {
	selectors, err := b.tagSelector(ctx, CommandGetOfferHelp)
	if err != nil {
		return errors.Wrap(err, "tag selector")
	}

	rsp := tgbotapi.NewMessage(msg.Chat.ID, "Выберите раздел")
	rsp.ReplyMarkup = selectors
	_, err = b.api.Send(rsp)
	return err
}

func (b *Bot) handleUnknown(_ context.Context, msg *tgbotapi.Message) error {
	rsp := tgbotapi.NewMessage(msg.Chat.ID, "Введите /askhelp или /offerhelp.")
	_, err := b.api.Send(rsp)
	return err
}

func (b *Bot) handleCallbackQuery(ctx context.Context, msg *tgbotapi.CallbackQuery) error {
	parts := strings.Split(msg.Data, ":")
	if len(parts) != 3 {
		return fmt.Errorf("'%s' not valid data command", msg.Data)
	}
	action, subject, idstr := Command(parts[0]), parts[1], parts[2]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("'%s' is not valid id", idstr)
	}

	//todo refactor this shit
	switch subject {
	case "tag":
		switch action {
		case CommandAskHelp, CommandOfferHelp:
			return b.addRequest(ctx, msg, action, id)
		case CommandGetAskHelp, CommandGetOfferHelp:
			return b.showRequests(ctx, msg, action, id)
		default:
			return fmt.Errorf("'%s' not valid action", action)
		}
	default:
		return fmt.Errorf("'%s' not valid subject", subject)
	}

	return nil
}

func (b *Bot) addRequest(ctx context.Context, msg *tgbotapi.CallbackQuery, action Command, id int) error {
	r := &model.Request{
		TgUsername: msg.From.UserName, //todo may be save more info (or not for security reason?)
		TagID:      id,
	}
	switch action {
	case CommandAskHelp:
		r.Type = model.RequestTypeIn
	case CommandOfferHelp:
		r.Type = model.RequestTypeOut
	default:
		return fmt.Errorf("'%s' not valid action", action)
	}

	err := b.requestService.Insert(ctx, r)
	if err != nil {
		return errors.Wrap(err, "request service insert")
	}

	rsp := tgbotapi.NewMessage(msg.Message.Chat.ID, "Ваша заявка принята")
	rsp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	_, err = b.api.Send(rsp)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) showRequests(ctx context.Context, msg *tgbotapi.CallbackQuery, action Command, id int) error {
	opts := model.RequestQueryOptions{
		TagID: []int{id},
	}
	switch action {
	case CommandGetAskHelp:
		opts.Type = append(opts.Type, model.RequestTypeIn)
	case CommandGetOfferHelp:
		opts.Type = append(opts.Type, model.RequestTypeOut)
	default:
		return fmt.Errorf("'%s' not valid action", action)
	}

	requests, err := b.requestService.GetByOptions(ctx, opts)
	if err != nil {
		return err
	}

	if requests == nil {
		rsp := tgbotapi.NewMessage(msg.Message.Chat.ID, "На данный момент заявок по данному тэгу нет")
		rsp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		_, err = b.api.Send(rsp)
		return err
	}

	byTag := map[model.Tag][]*model.Request{}
	for _, request := range requests {
		byTag[request.Tag] = append(byTag[request.Tag], request)
	}

	var responseText string
	for tag, requests := range byTag {
		responseText += fmt.Sprintf("%s:\n", tag.Name)
		var unames []string
		for _, request := range requests {
			unames = append(unames, fmt.Sprintf("@%s", request.TgUsername))
		}
		responseText += fmt.Sprintf("%s\n", strings.Join(unames, ", "))
	}

	rsp := tgbotapi.NewMessage(msg.Message.Chat.ID, responseText)
	rsp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	_, err = b.api.Send(rsp)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) tagSelector(ctx context.Context, action string) (tgbotapi.InlineKeyboardMarkup, error) {
	tags, err := b.tagService.GetByOptions(ctx, model.TagQueryOptions{})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(tags))
	for _, tag := range tags {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(tag.Name, fmt.Sprintf("%s:tag:%d", action, tag.ID))))
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...), nil
}

func (b *Bot) Debug(debug bool) {
	b.api.Debug = debug
}
