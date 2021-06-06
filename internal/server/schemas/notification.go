package schemas

type Notification struct {
    SubscriptionID int64  `json:"subscriptionID"`
    MusicURL       string `json:"musicURL"`
}