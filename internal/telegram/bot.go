package telegram

import (
	"bot-for-fans/internal/repository"
	"bot-for-fans/pkg/spotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type Bot struct {
	bot                  *tgbotapi.BotAPI
	userRepo             repository.UserRepository
	subscriptionRepo     repository.SubscriptionRepository
	userSubscriptionRepo repository.UserSubscriptionRepository
	spotify              *spotify.Client
	logger               *zap.Logger
}

type Config struct {
	B       *tgbotapi.BotAPI
	U       repository.UserRepository
	S       repository.SubscriptionRepository
	US      repository.UserSubscriptionRepository
	Spotify SpotifyCfg
	Logger  *zap.Logger
}

type SpotifyCfg struct {
	AppSecret           string
	AppClientID         string
	GetArtistsURI       string
	GetArtistsAlbumsURI string
	GetSeveralAlbumsURI string
}

func NewBot(c Config) *Bot {
	return &Bot{
		bot:                  c.B,
		userRepo:             c.U,
		subscriptionRepo:     c.S,
		userSubscriptionRepo: c.US,
		logger:               c.Logger,
		spotify: spotify.NewClient(spotify.Config{
			AppSecret:           c.Spotify.AppSecret,
			AppClientID:         c.Spotify.AppClientID,
			GetArtistsURI:       c.Spotify.GetArtistsURI,
			GetArtistsAlbumsURI: c.Spotify.GetArtistsAlbumsURI,
			GetSeveralAlbumsURI: c.Spotify.GetSeveralAlbumsURI,
		}),
	}
}

func (b *Bot) Start() error {
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.logger.Info("Bot successfully started")

	b.handleUpdatesChannel(updates)
	return nil
}

func (b *Bot) handleUpdatesChannel(u tgbotapi.UpdatesChannel) {
	for update := range u {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.logger.Error(
					"Error handling user command",
					zap.String("error", err.Error()),
					zap.Int64("chat_id", update.Message.Chat.ID),
					zap.String("first_name", update.Message.Chat.FirstName),
					zap.String("last_name", update.Message.Chat.LastName),
					zap.String("user_name", update.Message.Chat.UserName),
					zap.String("command", update.Message.Command()),
					zap.String("arguments", update.Message.CommandArguments()),
				)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.logger.Error(
				"Error handling user message",
				zap.String("error", err.Error()),
				zap.Int64("chat_id", update.Message.Chat.ID),
				zap.String("first_name", update.Message.Chat.FirstName),
				zap.String("last_name", update.Message.Chat.LastName),
				zap.String("user_name", update.Message.Chat.UserName),
				zap.String("content", update.Message.Text),
			)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}
