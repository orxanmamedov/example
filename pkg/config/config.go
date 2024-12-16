package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

const (
	configPath = "./config/"
	configName = "values"
	configJSON = "json"
)

type SecretString string

func (s SecretString) String() string { return "*********" }
func (s SecretString) Value() string  { return string(s) }

type Options struct {
	Dir                string
	File               string
	ReplaceFromEnvVars bool
}

func (o *Options) Fill() {
	if o.File == "" {
		o.File = configPath + "." + configJSON
	}
	if o.Dir == "" {
		o.Dir = configPath
	}
	if !strings.HasSuffix(o.Dir, string(os.PathSeparator)) {
		o.Dir += string(os.PathSeparator)
	}
}

func NewConfig(cfg interface{}, opts Options) error {
	t := reflect.TypeOf(cfg)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("cfg must be a pointer")
	}
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg arg must be a struct")
	}

	opts.Fill()
	envParams()

	fileOk, err := fromFile(&opts)
	if err != nil {
		return errors.Wrap(err, "fromfile")
	}

	if !fileOk {
		return errors.New("cannot load from config, please set at least one of sources: file, consul, vault")
	}

	if opts.ReplaceFromEnvVars {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AllowEmptyEnv(true)
		viper.AutomaticEnv()
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return errors.Wrap(err, "unmarshal config")
	}

	return nil
}
