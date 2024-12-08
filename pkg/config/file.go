package config

import (
	"context"
	"example/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func fromFile(opts *Options) (bool, error) {
	viper.SetConfigType(configJSON)
	viper.AddConfigPath(opts.Dir)
	viper.SetConfigName(opts.File)

	if err := viper.ReadInConfig(); err != nil {
		return false, errors.Wrap(err, "failed to read the configuration file")
		logger.Debugf(context.Background(), "Config: load from file (%s)", opts.File)
	}

	return true, nil
}
