package telegram

import (
    "bot-for-fans/internal/models"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) SendNotification(u models.User, newMusicURL string) error {
    msg := "<b>Hey " + u.GetAppeal() + "</b>\n" + "A new song has been released on one of your subscriptions:\n" + newMusicURL
    answer := tgbotapi.NewMessage(u.ChatID, msg)
    answer.ParseMode = "html"
    _, err := b.bot.Send(answer)
    return err
}