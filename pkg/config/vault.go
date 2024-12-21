package config

import (
	"context"
	"example/pkg/logger"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"os"
	"time"
)

const (
	vaultRoleID   = "VAULT_ROLE_ID"
	vaultSecretID = "VAULT_SECRET_ID"

	vaultTimeout = 10 * time.Second
)

func readVault(pathApp, pathKey string) (map[string]interface{}, error) {
	if os.Getenv(vault.EnvVaultAddress) == "" {
		logger.Debug(context.Background(), "Config: skip Vault - no ENV")
		return nil, nil
	}
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Timeout = vaultTimeout
	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "new client")
	}

	requestPath := "auth/approle/login"
	options := map[string]interface{}{
		"role_id":   os.Getenv(vaultRoleID),
		"secret_id": os.Getenv(vaultSecretID),
	}

	secret, err := vaultClient.Logical().Write(requestPath, options)
	if err != nil {
		return nil, err
	}

	vaultClient.SetToken(secret.Auth.ClientToken)
	data, err := vaultClient.Logical().Read(fmt.Sprintf("%s/data/%s", pathApp, pathKey))
	if err != nil {
		return nil, err
	}
	if data == nil {
		logger.Debugf(context.Background(), "Config: skip Vault - no data found (%s)", pathKey)
		return nil, nil
	}

	resp, ok := data.Data["data"]
	if !ok {
		return nil, errors.New("no date section from vault")
	}

	m, ok := resp.(map[string]interface{})
	if !ok {
		return nil, errors.New("date section does not match type map[string]interface{} from vault")
	}

	return m, nil
}

func GetSecret(path, key string) (map[string]interface{}, error) {
	return readVault(path, key)
}
