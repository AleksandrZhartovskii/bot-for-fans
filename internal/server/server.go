package server

import (
    "bot-for-fans/internal/repository"
    "bot-for-fans/internal/server/schemas"
    "bot-for-fans/internal/telegram"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/mailru/easyjson"
    "go.uber.org/zap"
)

type Server struct {
    bot      *telegram.Bot
    userRepo repository.UserRepository
    host     string
    logger   *zap.Logger
}

type Config struct {
    Bot      *telegram.Bot
    UserRepo repository.UserRepository
    Host     string
    Logger   *zap.Logger
}

func NewServer(c Config) *Server {
    return &Server{
        bot:      c.Bot,
        userRepo: c.UserRepo,
        host:     c.Host,
        logger:   c.Logger,
    }
}

func (s *Server) Start() error {
    server := fiber.New(fiber.Config{
        DisableStartupMessage: true,
        ReadBufferSize:        16384,
        ErrorHandler: func(ctx *fiber.Ctx, err error) error {
            c := fiber.StatusBadRequest
            if e, ok := err.(*fiber.Error); ok {
                c = e.Code
            }
            body, _ := easyjson.Marshal(&schemas.Response{
                Error: true,
                Data: schemas.ErrorResponse{
                    Message: err.Error(),
                },
            })
            return ctx.Status(c).Send(body)
        },
    })
    server.Use(cors.New())

    server.Use(func(ctx *fiber.Ctx) error {
        ctx.Response().Header.Set("Content-Type", "application/json")
        return ctx.Next()
    })

    server.Post("/notification", s.handleNotification)

    s.logger.Info("HTTP server successfully started on " + s.host)
    return server.Listen(s.host)
}