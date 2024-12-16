package config

import (
	"context"
	"encoding/json"
	"example/pkg/logger"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	configPathApp = "CONFIG_PATH_APP"
	configPathKey = "CONFIG_PATH_KEY"
)

func fomConsul(opts *Options) (bool, error) {
	pathApp := os.Getenv(configPathApp)
	pathKey := os.Getenv(configPathKey)

	if os.Getenv(consul.HTTPAddrEnvName) == "" || pathApp == "" || pathKey == "" {
		logger.Debug(context.Background(), "Config: skip Consul - no ENV found")
		return false, nil
	}

	cfg := consul.DefaultConfig()
	cfg.WaitTime = 10 * time.Second
	client, err := consul.NewClient(cfg)
	if err != nil {
		return false, errors.Wrapf(err, "new client")
	}
	kv := client.KV()
	path := fmt.Sprintf("%s/%s", pathApp, pathKey)
	kvp, _, err := kv.Get(path, nil)
	if err != nil {
		return false, nil
	}
	if kvp == nil {
		logger.Debugf(context.Background(), "Config: Consul path(%s) not found", path)
		return false, nil
	}

	viperConfig := make(map[string]interface{})

	if err := json.Unmarshal(kvp.Value, &viperConfig); err != nil {
		return false, err
	}
	if err := viper.MergeConfigMap(viperConfig); err != nil {
		return false, err
	}

	logger.Debugf(context.Background(), "Config: load from Consul (%s)", path)
	return true, nil
}
