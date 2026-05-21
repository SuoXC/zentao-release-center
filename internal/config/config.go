package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	ZentaoMini ZentaoMiniConfig `mapstructure:"zentao_mini"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Log       LogConfig       `mapstructure:"log"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type ZentaoMiniConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	Timeout      int    `mapstructure:"timeout"`
	ZentaoBaseURL string `mapstructure:"zentao_base_url"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./conf")

	v.SetEnvPrefix("ZENTAO_RELEASE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.expandPaths()
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("zentao_mini.base_url", "http://localhost:12345/api")
	v.SetDefault("zentao_mini.timeout", 120)
	v.SetDefault("database.path", "~/.zentao-release-center/release.db")
	v.SetDefault("log.level", "info")
}

func (c *Config) expandPaths() {
	if strings.HasPrefix(c.Database.Path, "~/") {
		home, _ := os.UserHomeDir()
		c.Database.Path = filepath.Join(home, c.Database.Path[2:])
	}
}
