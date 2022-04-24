package app

import (
	// "context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/config"
	"github.com/bigbag/go-musthave-diploma/internal/utils"
)

type Server struct {
	l logrus.FieldLogger
	f *fiber.App
}

func New(l logrus.FieldLogger, cfg *config.Config) *Server {
	fiberCfg := fiber.Config{
		ReadTimeout: time.Second * cfg.Server.ReadTimeout,
		IdleTimeout: time.Second * cfg.Server.IdleTimeout,
		Immutable:   true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			l.WithError(err).Error("Unexpected API error")
			return utils.SendJSONError(ctx, fiber.StatusInternalServerError, err.Error())
		},
	}

	f := fiber.New(fiberCfg)

	f.Use(recover.New())

	f.Use(logger.New(logger.Config{
		Output: l.(*logrus.Logger).Writer(),
	}))

	f.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	// f.Use(userid.New(userid.Config{
	// 	Secret:     cfg.UserCookieSecret,
	// 	ContextKey: cfg.UserContextKey,
	// }))

	// ctxBg := context.Background()
	// urlStorage, _ := storage.NewStorageService(ctxBg, cfg.Storage)

	// urlRepository := url.NewURLRepository(urlStorage)

	// urlPool := url.NewTaskPool(ctxBg, l, urlRepository)
	// urlService := url.NewURLService(l, urlRepository, urlPool)
	// url.NewURLHandler(f.Group(""), urlService, cfg, l)

	return &Server{l: l, f: f}
}

func (s *Server) Start(addr string) error {
	return s.f.Listen(addr)
}

func (s *Server) Stop() error {
	return s.f.Shutdown()
}
