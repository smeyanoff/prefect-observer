package config

import (
	"crm-uplift-ii24-backend/pkg/logging"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	DB   Dbconfig
	App  AppConfig
	CORS CORSConfig
}

type Dbconfig struct {
	Host     string
	Port     string
	User     string
	Pwd      string
	Database string
	SSLMode  string
}

type AppConfig struct {
	PrefectApiUrl           string
	InsecureSkipVerify      bool
	StageStatusQueryTimeout int
	NumWorkers              int
	Port                    string
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

var (
	cfg  Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		// Настраиваем Viper для работы с YAML-файлом
		viper.SetConfigName("backend-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("../configs")

		viper.SetEnvPrefix("observer")
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)

		if err := viper.ReadInConfig(); err != nil {
			logging.Fatal("Error reading config file", zap.String("err", err.Error()))
		}

		viper.AutomaticEnv()

		if err := viper.Unmarshal(&cfg); err != nil {
			logging.Fatal("Error config decode", zap.String("err", err.Error()))
		}

		logging.Info("Config downloaded...")
		logging.Debug("Config Fields", zap.Any("cfg", cfg))
	})

	return &cfg
}
