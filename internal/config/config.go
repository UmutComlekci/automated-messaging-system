package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("APP_NAME", "automated-messaging-system")
	viper.SetDefault("RELEASE", "v1.0.0")

	viper.SetDefault("SERVER_PORT", "8081")

	viper.SetDefault("SQL_DRIVER", "postgres")
	viper.SetDefault("POSTGRESQL_HOST", "localhost")
	viper.SetDefault("POSTGRESQL_PORT", "5432")
	viper.SetDefault("POSTGRESQL_USER", "admin")
	viper.SetDefault("POSTGRESQL_PASSWORD", "admin")
	viper.SetDefault("POSTGRESQL_DATABASE", "messagedb")
	viper.SetDefault("POSTGRESQL_SSL_MODE", "disable")

	viper.SetDefault("REDIS_URL", "redis://localhost:6379")

	viper.SetDefault("PROVIDER", "http")
	viper.SetDefault("PROVIDER_HTTP_URL", "https://webhook.site/2b1660b7-2eda-4ad8-b9b7-b0c8b65131be")
}

// App
func GetAppName() string { return viper.GetString("APP_NAME") }
func GetRelease() string { return viper.GetString("RELEASE") }

// Api
func GetApiPort() string { return viper.GetString("SERVER_PORT") }

// Database
func GetSqlDriver() string { return viper.GetString("SQL_DRIVER") }
func GetConnectionString() string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		viper.GetString("POSTGRESQL_HOST"),
		viper.GetInt("POSTGRESQL_PORT"),
		viper.GetString("POSTGRESQL_USER"),
		viper.GetString("POSTGRESQL_PASSWORD"),
		viper.GetString("POSTGRESQL_DATABASE"),
		viper.GetString("POSTGRESQL_SSL_MODE"),
	)
	return psqlInfo
}

// Cache
func GetCacheConnectionString() string { return viper.GetString("REDIS_URL") }

// Provider
func GetProvider() string        { return viper.GetString("PROVIDER") }
func GetProviderHttpUrl() string { return viper.GetString("PROVIDER_HTTP_URL") }
