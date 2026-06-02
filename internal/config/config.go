package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	ZentaoMini ZentaoMiniConfig `mapstructure:"zentao_mini"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Log        LogConfig        `mapstructure:"log"`
	GitLab     GitLabConfig     `mapstructure:"gitlab"`
	Lanxin     LanxinConfig     `mapstructure:"lanxin"`
	Email      EmailConfig      `mapstructure:"email"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type ZentaoMiniConfig struct {
	BaseURL       string `mapstructure:"base_url"`
	Timeout       int    `mapstructure:"timeout"`
	ZentaoBaseURL string `mapstructure:"zentao_base_url"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type GitLabConfig struct {
	BaseURL       string `mapstructure:"base_url"`
	Token         string `mapstructure:"token"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

type LanxinConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
	Secret  string `mapstructure:"secret"`
	SkipSSL bool   `mapstructure:"skip_ssl"`
}

type EmailConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	SMTPHost      string `mapstructure:"smtp_host"`
	SMTPPort      int    `mapstructure:"smtp_port"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	SenderName    string `mapstructure:"sender_name"`
	SenderAddress string `mapstructure:"sender_address"`
	UseSSL        bool   `mapstructure:"use_ssl"`
	Recipients    string `mapstructure:"recipients"`
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
	v.SetDefault("gitlab.base_url", "")
	v.SetDefault("gitlab.token", "")
	v.SetDefault("gitlab.webhook_secret", "")
	v.SetDefault("lanxin.enabled", false)
	v.SetDefault("lanxin.url", "")
	v.SetDefault("lanxin.secret", "")
	v.SetDefault("lanxin.skip_ssl", false)
	v.SetDefault("email.enabled", false)
	v.SetDefault("email.smtp_host", "")
	v.SetDefault("email.smtp_port", 465)
	v.SetDefault("email.username", "")
	v.SetDefault("email.password", "")
	v.SetDefault("email.sender_name", "发布中心")
	v.SetDefault("email.sender_address", "")
	v.SetDefault("email.use_ssl", true)
	v.SetDefault("email.recipients", "")
}

func (c *Config) expandPaths() {
	if strings.HasPrefix(c.Database.Path, "~/") {
		home, _ := os.UserHomeDir()
		c.Database.Path = filepath.Join(home, c.Database.Path[2:])
	}
}
