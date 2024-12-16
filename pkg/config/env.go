package config

import "os"

var (
	AppName = "unknown"
	EnvName = "unknown"

	HostName = "unknown"
)

const (
	appNameKey  = "APP_NAME"
	envNameKey  = "ENVIRONMENT_NAME"
	hostNameKey = "HOSTNAME"
)

func envParams() {
	if appName, ok := os.LookupEnv(appNameKey); ok && appName != "" {
		AppName = appName
	}

	if envName, ok := os.LookupEnv(envNameKey); ok && envName != "" {
		EnvName = envName
	}

	hostName := os.Getenv(hostNameKey)
	if hostName == "" {
		hostName, _ = os.Hostname()
	}

	if hostName != "" {
		HostName = hostName
	}
}
