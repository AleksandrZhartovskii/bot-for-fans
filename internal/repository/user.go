package repository

import (
    "bot-for-fans/internal/models"
)

type UserRepository interface {
    GetByChatID(chatID int64) (models.User, error)
    Create(chatID int64, firstName, lastName, userName string) (id int64, err error)
    GetAllBySubscription(subscriptionID int64) ([]models.User, error)
}