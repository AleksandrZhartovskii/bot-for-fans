package models

import "time"

type Subscription struct {
    ID              int64
    AuthorSpotifyID string
    AuthorName      string
    LastUpdateTime  time.Time
}