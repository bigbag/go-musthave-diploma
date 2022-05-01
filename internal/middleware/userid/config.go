package userid

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Next       func(c *fiber.Ctx) bool
	CookieName string
	ContextKey string
	Secret     string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:       nil,
	CookieName: "auth",
	ContextKey: "userid",
	Secret:     "secret",
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	if cfg.CookieName == "" {
		cfg.CookieName = ConfigDefault.CookieName
	}

	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigDefault.ContextKey
	}

	if cfg.Secret == "" {
		cfg.Secret = ConfigDefault.Secret
	}

	return cfg
}
