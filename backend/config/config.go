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
	Server   *ServerConfiguration
	DB       *DBConfiguration
	OpenClaw *OpenClawConfiguration
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

type OpenClawConfiguration struct {
	GatewayWSURL       string
	GatewayToken       string
	CommandTimeoutMS   int
	ActiveWindowMinute int
	DefaultSessionKey  string
	DeviceStorePath    string
	ClientID           string
	ClientVersion      string
	ClientPlatform     string
	ClientDeviceFamily string
	ReconnectDelayMS   int
}

func Load() (*Configurations, error) {
	// set default app config
	viper.SetDefault("SERVER_NAME", "CLAW_EMPIRE")
	viper.SetDefault("SERVER_HOST", "127.0.0.1")
	viper.SetDefault("SERVER_PORT", "8212")
	viper.SetDefault("SERVER_MODE", MODE_DEV)
	viper.SetDefault("OPENCLAW_GATEWAY_WS_URL", "ws://127.0.0.1:18789/ws")
	viper.SetDefault("OPENCLAW_COMMAND_TIMEOUT_MS", 15000)
	viper.SetDefault("OPENCLAW_ACTIVE_WINDOW_MINUTE", 15)
	viper.SetDefault("OPENCLAW_DEFAULT_SESSION_KEY", "main")
	viper.SetDefault("OPENCLAW_DEVICE_STORE_PATH", "/app/data/openclaw/device.json")
	viper.SetDefault("OPENCLAW_CLIENT_ID", "openclaw-empire")
	viper.SetDefault("OPENCLAW_CLIENT_VERSION", "1.0.0")
	viper.SetDefault("OPENCLAW_CLIENT_PLATFORM", "linux")
	viper.SetDefault("OPENCLAW_CLIENT_DEVICE_FAMILY", "server")
	viper.SetDefault("OPENCLAW_RECONNECT_DELAY_MS", 3000)

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
	openClawCnf := initOpenClawConfig()

	return &Configurations{
		Server:   serverCnf,
		DB:       dbCnf,
		OpenClaw: openClawCnf,
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

func initOpenClawConfig() *OpenClawConfiguration {
	return &OpenClawConfiguration{
		GatewayWSURL:       viper.GetString("OPENCLAW_GATEWAY_WS_URL"),
		GatewayToken:       viper.GetString("OPENCLAW_GATEWAY_TOKEN"),
		CommandTimeoutMS:   viper.GetInt("OPENCLAW_COMMAND_TIMEOUT_MS"),
		ActiveWindowMinute: viper.GetInt("OPENCLAW_ACTIVE_WINDOW_MINUTE"),
		DefaultSessionKey:  viper.GetString("OPENCLAW_DEFAULT_SESSION_KEY"),
		DeviceStorePath:    viper.GetString("OPENCLAW_DEVICE_STORE_PATH"),
		ClientID:           viper.GetString("OPENCLAW_CLIENT_ID"),
		ClientVersion:      viper.GetString("OPENCLAW_CLIENT_VERSION"),
		ClientPlatform:     viper.GetString("OPENCLAW_CLIENT_PLATFORM"),
		ClientDeviceFamily: viper.GetString("OPENCLAW_CLIENT_DEVICE_FAMILY"),
		ReconnectDelayMS:   viper.GetInt("OPENCLAW_RECONNECT_DELAY_MS"),
	}
}

// -------------------- helper function --------------------
func (db *DBConfiguration) GetConnStr() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		db.host, db.user, db.pass, db.name, db.port)
}
