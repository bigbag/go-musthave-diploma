package config

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Storage struct {
	DatabaseDSN string        `envconfig:"DATABASE_URI"`
	ConnTimeout time.Duration `envconfig:"STORAGE_CONNECTION_TIMEOUT" default:"3s"`
	StopTimeout time.Duration `envconfig:"STORAGE_STOP_TIMEOUT" default:"3s"`
}

type Auth struct {
	SecretKey   string        `envconfig:"AUTH_SECRET_KEY" default:"simple_auth"`
	CookieName  string        `envconfig:"AUTH_COOKIE_NAME" default:"auth"`
	ExpiresTime time.Duration `envconfig:"AUTH_EXPIRES_TIME" default:"24h"`
}

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"gophermart"`
	Auth        *Auth
	Server      struct {
		Listen      string        `envconfig:"RUN_ADDRESS"  default:":8080"`
		ReadTimeout time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT" default:"5s"`
	}
	AccrualURL string `envconfig:"ACCRUAL_SYSTEM_ADDRESS"`
	Storage    *Storage
	Logger     struct {
		Level  string `envconfig:"LOG_LEVEL" default:"info"`
		Output string `envconfig:"LOG_OUTPUT" default:"stdout"`
		Format string `envconfig:"LOG_FORMAT" default:"text"`
	}
}

// New parses environments and creates new instance of config
func New() (*Config, error) {
	cfg := new(Config)

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.Server.Listen, "a", cfg.Server.Listen, "listen address. env: RUN_ADDRESS")
	flag.StringVar(&cfg.Storage.DatabaseDSN, "d", cfg.Storage.DatabaseDSN, "database dsn. env: DATABASE_URI")
	flag.StringVar(&cfg.AccrualURL, "r", cfg.AccrualURL, "accrual service url. env: ACCRUAL_SYSTEM_ADDRESS")
	flag.Parse()

	return cfg, nil
}

func (c *Config) String() string {
	if out, err := json.MarshalIndent(&c, "", "  "); err == nil {
		return string(out)
	}
	return ""
}
