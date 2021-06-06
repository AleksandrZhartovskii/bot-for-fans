package main

import (
	"bot-for-fans/internal/repository/pgdb"
	"bot-for-fans/internal/server"
	"bot-for-fans/internal/telegram"
	"bot-for-fans/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := dialDB()
	if err != nil {
		log.Fatal(err)
	}

	if err = runPgMigrations(); err != nil {
		log.Fatal(err)
	}

	userRepo := pgdb.NewUserRepository(db)
	subscriptionRepo := pgdb.NewSubscriptionRepository(db)
	userSubscriptionRepo := pgdb.NewUserSubscriptionRepository(db)

	telegramBot := telegram.NewBot(telegram.Config{
		B:      bot,
		U:      userRepo,
		S:      subscriptionRepo,
		US:     userSubscriptionRepo,
		Logger: l,
		Spotify: telegram.SpotifyCfg{
			AppSecret:           os.Getenv("SPOTIFY_APP_SECRET"),
			AppClientID:         os.Getenv("SPOTIFY_APP_CLIENT_ID"),
			GetArtistsURI:       os.Getenv("SPOTIFY_GET_ARTISTS_URI"),
			GetArtistsAlbumsURI: os.Getenv("SPOTIFY_GET_ARTISTS_ALBUMS_URI"),
			GetSeveralAlbumsURI: os.Getenv("SPOTIFY_GET_SEVERAL_ALBUMS_URI"),
		},
	})
	go func() {
		if err = telegramBot.Start(); err != nil {
			l.Fatal(err.Error())
		}
	}()

	notificationServer := server.NewServer(server.Config{
		Bot:      telegramBot,
		UserRepo: userRepo,
		Host:     os.Getenv("SERVER_HOST"),
		Logger:   l,
	})
	if err = notificationServer.Start(); err != nil {
		l.Fatal(err.Error())
	}
}

func dialDB() (*sqlx.DB, error) {
	db, err := sqlx.Open(os.Getenv("DB_DRIVER_NAME"), os.Getenv("DB_DATA_SOURCE"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10 / 2)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}

func runPgMigrations() error {
	m, err := migrate.New(os.Getenv("PATH_TO_MIGRATIONS"), os.Getenv("DB_DATA_SOURCE"))
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
