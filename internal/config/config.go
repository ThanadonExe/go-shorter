package config

import (
	"errors"

	"github.com/spf13/viper"
)

var AppConfig EnvConfig

type EnvConfig struct {
	Env    string `mapstructure:"APP_ENV"`
	Port   int    `mapstructure:"APP_PORT"`
	Domain string `mapstructure:"APP_DOMAIN"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`

	REDISHost     string `mapstructure:"REDIS_HOST"`
	REDISPort     string `mapstructure:"REDIS_PORT"`
	REDISPassword string `mapstructure:"REDIS_PASSWORD"`
	REDISDb       int    `mapstructure:"REDIS_DB"`

	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpiredIn int64  `mapstructure:"JWT_EXPIRED_IN"`
}

func InitAppConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("failed to load config")
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return errors.New("failed to parse config struct")
	}

	return nil
}
