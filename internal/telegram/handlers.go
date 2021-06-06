package telegram

import (
    "bot-for-fans/internal/models"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "go.uber.org/zap"
    "strings"
)

const (
	cmdStart     = "start"
	cmdSubscribe = "subscribe"
	cmdHelp      = "help"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) error {
    switch msg.Command() {
    case cmdStart:
        b.logger.Info(
            "The user entered the /start command",
            zap.Int64("chat_id", msg.Chat.ID),
            zap.String("first_name", msg.Chat.FirstName),
            zap.String("last_name", msg.Chat.LastName),
            zap.String("user_name", msg.Chat.UserName),
            zap.String("arguments", msg.CommandArguments()),
        )
        return b.handleStartCmd(msg)

    case cmdSubscribe:
        b.logger.Info(
            "The user entered the /subscribe command",
            zap.Int64("chat_id", msg.Chat.ID),
            zap.String("first_name", msg.Chat.FirstName),
            zap.String("last_name", msg.Chat.LastName),
            zap.String("user_name", msg.Chat.UserName),
            zap.String("arguments", msg.CommandArguments()),
        )
        return b.handleSubscribeCmd(msg)

    case cmdHelp:
        b.logger.Info(
            "The user entered the /help command",
            zap.Int64("chat_id", msg.Chat.ID),
            zap.String("first_name", msg.Chat.FirstName),
            zap.String("last_name", msg.Chat.LastName),
            zap.String("user_name", msg.Chat.UserName),
            zap.String("arguments", msg.CommandArguments()),
        )
        return b.handleHelpCmd(msg)

    default:
        b.logger.Warn(
            "The user entered the unknown command",
            zap.Int64("chat_id", msg.Chat.ID),
            zap.String("first_name", msg.Chat.FirstName),
            zap.String("last_name", msg.Chat.LastName),
            zap.String("user_name", msg.Chat.UserName),
            zap.String("arguments", msg.CommandArguments()),
        )
        return b.handleUnknownCmd(msg)
    }
}

func (b *Bot) handleStartCmd(msg *tgbotapi.Message) error {
    if _, err := b.userRepo.GetByChatID(msg.Chat.ID); err != nil {
        _, err = b.userRepo.Create(msg.Chat.ID, msg.From.FirstName, msg.From.LastName, msg.From.UserName)
        if err != nil {
            return err
        }
    }

    user := models.User{
        FirstName: msg.From.FirstName,
        LastName:  msg.From.LastName,
        UserName:  msg.From.UserName,
    }

    text := "<b>Hey " + user.GetAppeal() + " !</b>\n"
    text += "\n"
    text += "I provide a subscription to your favorite music authors.\n"
    text += "How it works? Everything is very simple! You just enter the command /subscribe [Author Name] and I will send you here links to only released tracks on your subscriptions!\n"
    text += "\n"
    text += "And if you want to know what else I can do, then just enter the command /help\n"
    text += "\n"
    text += "<b>Forward!</b>"

    answer := tgbotapi.NewMessage(msg.Chat.ID, text)
    answer.ParseMode = "html"
    _, err := b.bot.Send(answer)
    return err
}

func (b *Bot) handleSubscribeCmd(msg *tgbotapi.Message) error {
    if msg.CommandArguments() == "" {
        text := "Sorry, but need add author name 😬"

        answer := tgbotapi.NewMessage(msg.Chat.ID, text)
        _, err := b.bot.Send(answer)
        return err
    }

    user, err := b.userRepo.GetByChatID(msg.Chat.ID)
    if err != nil {
        return err
    }

    authorName := strings.TrimSpace(strings.ToLower(msg.CommandArguments()))
    spotifyAuthors, err := b.spotify.GetArtistsByName(authorName)
    if err != nil {
        return err
    }

    if len(spotifyAuthors) == 0 {
        return b.sendAuthorNotFound(msg.Chat.ID, msg.CommandArguments())
    }

    position := b.findAuthorInSpotifyResult(spotifyAuthors, authorName)
    if position == -1 {
        return b.sendRecommendedAuthors(msg.Chat.ID, spotifyAuthors)
    }

    author := spotifyAuthors[position]

    subscriptionID, err := b.subscriptionRepo.GetOrCreate(author.Name, author.ID)
    if err != nil {
        return err
    }

    err = b.userSubscriptionRepo.Create(user.ID, subscriptionID)
    if err != nil {
        return err
    }

    text := "<b>Congratulations!</b>\n"
    text += "You just signed up for music creation updates by author <u>" + author.Name + "</u> 🤟"

    answer := tgbotapi.NewMessage(msg.Chat.ID, text)
    answer.ParseMode = "html"
    _, err = b.bot.Send(answer)
    return err
}

func (b *Bot) handleHelpCmd(msg *tgbotapi.Message) error {
    text := "<b>Here's what I can do:</b>\n"
    text += "\n"
    text += "/start - Getting started with a bot\n"
    text += "/subscribe Author Name - Subscribe to author\n"
    text += "/help - View command list\n"
    text += "\n"
    text += "<b>Forward!</b> 😎"

    answer := tgbotapi.NewMessage(msg.Chat.ID, text)
    answer.ParseMode = "html"
    _, err := b.bot.Send(answer)
    return err
}

func (b *Bot) handleUnknownCmd(msg *tgbotapi.Message) error {
    const defaultCmdResp = "I'm sorry, but I don't know such a command 🤷"
    answer := tgbotapi.NewMessage(msg.Chat.ID, defaultCmdResp)
    _, err := b.bot.Send(answer)
    return err
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
    b.logger.Info(
        "The user entered the simple message",
        zap.Int64("chat_id", msg.Chat.ID),
        zap.String("first_name", msg.Chat.FirstName),
        zap.String("last_name", msg.Chat.LastName),
        zap.String("user_name", msg.Chat.UserName),
        zap.String("content", msg.Text),
    )

    const answerText = "I would like to chat with you, but unfortunately the developer forbade me this 🤐"
    answer := tgbotapi.NewMessage(msg.Chat.ID, answerText)
    _, err := b.bot.Send(answer)
    return err
}