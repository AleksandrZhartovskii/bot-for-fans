package telegram

import (
    "bot-for-fans/pkg/spotify/schemas"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

func (b *Bot) findAuthorInSpotifyResult(authors []schemas.Artist, authorName string) (index int) {
    for i := range authors {
        if strings.ToLower(strings.TrimSpace(authors[i].Name)) == authorName {
            return i
        }
    }
    return -1
}

func (b *Bot) sendRecommendedAuthors(chatID int64, authors []schemas.Artist) error {
    text := "<b>Sorry, but I have not found such an author</b>\n"
    text += "<b>But maybe one of these will suit you:</b>\n"
    text += "\n"
    for i, a := range authors {
        text += fmt.Sprintf("%d) %s\n", i+1, a.Name)
    }
    text += "\n"
    text += "<b>I hope I helped find what you wanted</b> 😅"

    answer := tgbotapi.NewMessage(chatID, text)
    answer.ParseMode = "html"
    _, err := b.bot.Send(answer)
    return err
}

func (b *Bot) sendAuthorNotFound(chatID int64, authorName string) error {
    text := fmt.Sprintf("Sorry, but I have not found an author like «%s» 😰", authorName)
    answer := tgbotapi.NewMessage(chatID, text)
    answer.ParseMode = "html"
    _, err := b.bot.Send(answer)
    return err
}