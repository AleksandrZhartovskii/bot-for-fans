package repository

type UserSubscriptionRepository interface {
    Create(userID, subscriptionID int64) error
}