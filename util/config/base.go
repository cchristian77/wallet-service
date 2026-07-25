package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var c *Config

type (
	Config struct {
		App      AppConfig      `env:"app"`
		Database DatabaseConfig `env:"database"`
		Context  ContextConfig  `env:"context"`
	}

	AppConfig struct {
		Name      string `env:"name"`
		Env       string `env:"env"`
		Version   string `env:"version"`
		Port      int32  `env:"port"`
		APIPrefix string `env:"api_prefix"`
	}

	DatabaseConfig struct {
		Host     string `env:"host"`
		Port     int32  `env:"port"`
		User     string `env:"user"`
		Password string `env:"password"`
		Name     string `env:"name"`
	}

	ContextConfig struct {
		Timeout string `env:"timeout"`
	}
)

func LoadConfig() error {
	k := koanf.New(".")

	if err := k.Load(file.Provider("./env.json"), json.Parser()); err != nil {
		return err
	}

	if err := k.UnmarshalWithConf("env", &c, koanf.UnmarshalConf{Tag: "env"}); err != nil {
		return err
	}

	return nil
}

func Env() *Config {
	return c
}

func (c *Config) DSN() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
	)
}

func (c *Config) DatabaseURL() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}
