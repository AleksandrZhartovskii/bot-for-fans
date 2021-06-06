package server

import (
    "bot-for-fans/internal/models"
    "bot-for-fans/internal/server/schemas"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

func (s *Server) handleNotification(ctx *fiber.Ctx) error {
    var notification schemas.Notification
    if err := notification.UnmarshalJSON(ctx.Body()); err != nil {
        return err
    }

    users, err := s.userRepo.GetAllBySubscription(notification.SubscriptionID)
    if err != nil {
        s.logger.Error(
            "Error getting all users by subscription",
            zap.String("error", err.Error()),
            zap.Int64("subscription_id", notification.SubscriptionID),
        )
        return err
    }

    for i := range users {
        go func(user models.User) {
            if err = s.bot.SendNotification(user, notification.MusicURL); err != nil {
                s.logger.Error(
                    "Error sending notification to user",
                    zap.String("error", err.Error()),
                    zap.Int64("chat_id", user.ID),
                    zap.String("first_name", user.FirstName),
                    zap.String("last_name", user.LastName),
                    zap.String("user_name", user.UserName),
                    zap.String("music_url", notification.MusicURL),
                )
            }
        }(users[i])
    }

    return ctx.SendStatus(fiber.StatusOK)
}