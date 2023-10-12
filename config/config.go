package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const fileConfig = "../.env"

type Config struct {
	PostgresDriver   string `mapstructure:"POSTGRES_DRIVER"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresPath     string `mapstructure:"POSTGRES_PATH"`
	PostgresSslmode  string `mapstructure:"POSTGRES_SSLMODE"`
}

func ReadConfig() (*Config, error) {
	var cfg Config
	viper.SetConfigFile(fileConfig)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Não foi possível achar ou ler o arquivo de configuração")
		log.Println(err.Error())
	}
	if err := viper.UnmarshalExact(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s", c.PostgresDriver, c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresPath, c.PostgresSslmode)
}
