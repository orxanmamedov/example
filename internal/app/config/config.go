package config

import (
	"context"
	"example/pkg/config"
	"example/pkg/logger"
	"os"
	"path"
)

func GetConfig() *Config {
	confOpts := config.Options{
		ReplaceFromEnvVars: true,
	}

	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = os.Getenv("CONFIG_FILE")
	}
	if configPath != "" {
		confOpts.Dir = path.Dir(configPath)
		confOpts.File = path.Base(configPath)
	}

	return NewConfig(confOpts)
}

func NewConfig(opts config.Options) *Config {
	var cfg = &Config{
		App: App{
			Env:      "undefined",
			Name:     "undefined",
			LogLevel: "debug",
		},
	}

	err := config.NewConfig(cfg, opts)
	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	if cfg.App.Name == "" {
		cfg.App.Name = config.AppName
	}

	if cfg.App.Env == "" {
		cfg.App.Env = config.EnvName
	}

	return cfg
}
