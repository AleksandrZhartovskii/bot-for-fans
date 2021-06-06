package main

import (
    "bot-for-fans/internal/repository/pgdb"
    "bot-for-fans/internal/tracker"
    "bot-for-fans/pkg/logger"
    "github.com/jmoiron/sqlx"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "log"
    "os"
    "time"
)

func main() {
    l, err := logger.NewLogger("bot", logger.LogLevelDebug)
    if err != nil {
        log.Fatal(err)
    }

    db, err := sqlx.Open(os.Getenv("DB_DRIVER_NAME"), os.Getenv("DB_DATA_SOURCE"))
    if err != nil {
        log.Fatal(err)
    }
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10 / 2)
    db.SetConnMaxLifetime(5 * time.Minute)
    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    tracker.NewTracker(tracker.Config{
        BotURL:           "http://localhost:8080",
        SubscriptionRepo: pgdb.NewSubscriptionRepository(db),
        Logger:           l,
        Spotify:          tracker.SpotifyCfg{
            AppSecret:           os.Getenv("SPOTIFY_APP_SECRET"),
            AppClientID:         os.Getenv("SPOTIFY_APP_CLIENT_ID"),
            GetArtistsURI:       os.Getenv("SPOTIFY_GET_ARTISTS_URI"),
            GetArtistsAlbumsURI: os.Getenv("SPOTIFY_GET_ARTISTS_ALBUMS_URI"),
            GetSeveralAlbumsURI: os.Getenv("SPOTIFY_GET_SEVERAL_ALBUMS_URI"),
        },
    }).Start()
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal(err)
    }
}