package config

import (
	"example/pkg/config"
	"time"
)

type Config struct {
	App       App    `mapstructure:"app"`
	Server    Server `mapstructure:"server"`
	DB        DB     `mapstructure:"db"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

type App struct {
	Env        string  `mapstructure:"env"`
	Name       string  `mapstructure:"name"`
	LogLevel   string  `mapstructure:"log_level"`
	TraceRatio float64 `mapstructure:"trace_ratio"`
}

type Server struct {
	HTTPPort       *uint `mapstructure:"http_port"`
	GrpcPort       *uint `mapstructure:"grpc_port"`
	MonitoringPort uint  `mapstructure:"monitoring_port"`
	Logging        bool  `mapstructure:"logging"`
	WithSwagger    bool  `mapstructure:"with_swagger"`
}

type DB struct {
	Master     DBConfig `mapstructure:"master"`
	Slave      DBConfig `mapstructure:"slave"`
	Metrics    bool     `mapstructure:"metrics"`
	Migrations bool     `mapstructure:"migrations"`
}

type DBConfig struct {
	Host     string              `mapstructure:"host"`
	Port     string              `mapstructure:"port"`
	User     string              `mapstructure:"user"`
	Password config.SecretString `mapstructure:"password"`
	Database string              `mapstructure:"database"`
	MaxOpen  uint                `mapstructure:"max_open"`
	Timeout  time.Duration       `mapstructure:"timeout"`
}
