package pgdb

import (
    "github.com/jmoiron/sqlx"
)

type UserSubscriptionRepository struct {
    db *sqlx.DB
}

func NewUserSubscriptionRepository(db *sqlx.DB) *UserSubscriptionRepository {
    return &UserSubscriptionRepository{db: db}
}

const createUserSubscriptionQuery = `
insert into telegram.user_subscriptions (user_id, subscription_id) 
values ($1, $2);
`

func (r *UserSubscriptionRepository) Create(userID, subscriptionID int64) error {
    rows, err := r.db.Queryx(createUserSubscriptionQuery, userID, subscriptionID)
    defer func() { _ = rows.Close() }()
    return err
}