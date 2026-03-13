package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	MODE_DEV  = "DEV"
	MODE_TEST = "test"
	MODE_PROD = "prod"
)

type Configurations struct {
	Server *ServerConfiguration
	DB     *DBConfiguration
}

type ServerConfiguration struct {
	Host string
	Port string
	Name string
	Mode string
}

type DBConfiguration struct {
	host string
	port string
	user string
	pass string
	name string
}

func Load() (*Configurations, error) {
	// set default app config
	viper.SetDefault("SERVER_NAME", "CLAW_EMPIRE")
	viper.SetDefault("SERVER_HOST", "127.0.0.1")
	viper.SetDefault("SERVER_PORT", "8212")
	viper.SetDefault("SERVER_MODE", MODE_DEV)

	// read .env overwrite config
	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Warn("No .env file found, relying on system environment variables.")
		} else {
			return nil, fmt.Errorf("failed to read .env config: %s \n", err)
		}
	}

	serverCnf := initServerConfig()
	dbCnf := initDBConfig()

	return &Configurations{
		Server: serverCnf,
		DB:     dbCnf,
	}, nil
}

// -------------------- init function --------------------
func initServerConfig() *ServerConfiguration {
	return &ServerConfiguration{
		Host: viper.GetString("SERVER_HOST"),
		Port: viper.GetString("SERVER_PORT"),
		Name: viper.GetString("SERVER_NAME"),
		Mode: viper.GetString("SERVER_MODE"),
	}
}

func initDBConfig() *DBConfiguration {
	return &DBConfiguration{
		host: viper.GetString("DB_HOST"),
		port: viper.GetString("DB_PORT"),
		user: viper.GetString("DB_USER"),
		pass: viper.GetString("DB_PASS"),
		name: viper.GetString("DB_NAME"),
	}
}

// -------------------- helper function --------------------
func (db *DBConfiguration) GetConnStr() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		db.host, db.user, db.pass, db.name, db.port)
}
