package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBName            string `mapstructure:"DB_NAME"`
	DBUsername        string `mapstructure:"DB_USERNAME"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBPostgresSslMode string `mapstructure:"DB_POSTGRES_SSL_MODE"`

	JWTIssuer             string `mapstructure:"JWT_ISSUER"`
	JWTAccessExpireMinute string `mapstructure:"JWT_ACCESS_EXP_MINUTE"`
	JWTAccessSecret       string `mapstructure:"JWT_ACCESS_SECRET"`

	JWTRefreshExpireMinute string `mapstructure:"JWT_REFRESH_EXP_MINUTE"`
	JWTRefreshSecret       string `mapstructure:"JWT_REFRESH_SECRET"`
}

func initConfig() {
	ENV := os.Getenv("ENV")
	if ENV == "DEPLOY" {
		viper.AutomaticEnv()

		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("DB_NAME")
		viper.BindEnv("DB_USERNAME")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_POSTGRES_SSL_MODE")

		viper.BindEnv("JWT_ISSUER")

		viper.BindEnv("JWT_ACCESS_EXP_MINUTE")
		viper.BindEnv("JWT_ACCESS_SECRET")

		viper.BindEnv("JWT_REFRESH_EXP_MINUTE")
		viper.BindEnv("JWT_REFRESH_SECRET")
	} else {
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Error reading config file: ", err)
			os.Exit(1)
		}
	}
}

func InitConfigDsn() string {
	initConfig()

	var config Configuration

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("[Config][InitConfigDsn] Unable to decode into struct:", err)
	}

	dsn := `postgres://` + config.DBUsername + `:` + config.DBPassword + `@` + config.DBHost + `:` + config.DBPort + `/` + config.DBName + `?sslmode=` + config.DBPostgresSslMode

	return dsn
}

func InitConfigJwt() []string {
	initConfig()

	var config Configuration

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("[Config][InitConfigDsn] Unable to decode into struct:", err)
	}

	return []string{
		config.JWTIssuer,
		config.JWTAccessExpireMinute,
		config.JWTAccessSecret,
		config.JWTRefreshExpireMinute,
		config.JWTRefreshSecret,
	}
}
