package tracker

import (
    "bot-for-fans/internal/repository"
    "bot-for-fans/pkg/spotify"
    "go.uber.org/zap"
    "time"
)

type T struct {
    botURL  string
    repo    repository.SubscriptionRepository
    spotify *spotify.Client
    logger  *zap.Logger
}

type Config struct {
    BotURL           string
    SubscriptionRepo repository.SubscriptionRepository
    Spotify          SpotifyCfg
    Logger           *zap.Logger
}

type SpotifyCfg struct {
    AppSecret           string
    AppClientID         string
    GetArtistsURI       string
    GetArtistsAlbumsURI string
    GetSeveralAlbumsURI string
}

func NewTracker(c Config) *T {
    return &T{
        botURL:  c.BotURL,
        repo:    c.SubscriptionRepo,
        logger:  c.Logger,
        spotify: spotify.NewClient(spotify.Config{
            AppSecret:           c.Spotify.AppSecret,
            AppClientID:         c.Spotify.AppClientID,
            GetArtistsURI:       c.Spotify.GetArtistsURI,
            GetArtistsAlbumsURI: c.Spotify.GetArtistsAlbumsURI,
            GetSeveralAlbumsURI: c.Spotify.GetSeveralAlbumsURI,
        }),
    }
}

func (t *T) Start() {
    t.logger.Info("Release-tracker service successfully started")
    ticker := time.Tick(time.Hour*24)
    for {
        select {
        case <-ticker:
            t.checkSubscriptions()
        }
    }
}