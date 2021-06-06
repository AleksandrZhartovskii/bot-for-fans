package repository

import (
    "bot-for-fans/internal/models"
    "time"
)

type SubscriptionRepository interface {
    GetOrCreate(authorName, authorID string) (id int64, err  error)
    GetAll() ([]models.Subscription, error)
    ChangeLastUpdateTime(id int64, t time.Time) error
}