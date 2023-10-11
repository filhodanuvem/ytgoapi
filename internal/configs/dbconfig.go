package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig() string {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s",
		viper.Get("DB_DRIVER"),
		viper.Get("DB_USER"),
		viper.Get("DB_PASSWORD"),
		viper.Get("DB_HOST"),
		viper.Get("DB_PORT"),
		viper.Get("DB_PATH"),
		viper.Get("DB_SSLMODE"),
	)
}
