package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/config"
	"github.com/bigbag/go-musthave-diploma/internal/middleware/userid"
	"github.com/bigbag/go-musthave-diploma/internal/order"
	"github.com/bigbag/go-musthave-diploma/internal/storage"
	"github.com/bigbag/go-musthave-diploma/internal/user"
	"github.com/bigbag/go-musthave-diploma/internal/utils"
)

type Server struct {
	l  logrus.FieldLogger
	f  *fiber.App
	sr *storage.StorageRepository
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

	ctxBg := context.Background()
	sr, _ := storage.NewStorageRepository(
		ctxBg, cfg.Storage.DatabaseDSN, cfg.Storage.ConnTimeout,
	)
	ur := user.NewUserRepository(ctxBg, l, sr, cfg.Storage.ConnTimeout)
	us := user.NewUserService(l, ur)

	user.NewUserHandler(f.Group("/api/user/"), l, cfg, us)

	authMiddleware := userid.New(userid.Config{
		CookieName: cfg.Auth.CookieName,
		ContextKey: cfg.Auth.ContextKey,
		Secret:     cfg.Auth.SecretKey,
	})

	order.NewOrderHandler(f.Group("/api/user/", authMiddleware), l, cfg)

	return &Server{l: l, f: f, sr: sr}
}

func (s *Server) Start(addr string) error {
	return s.f.Listen(addr)
}

func (s *Server) Stop() error {
	s.sr.Close()
	return s.f.Shutdown()
}
