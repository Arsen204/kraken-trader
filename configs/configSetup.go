package configs

import (
	"course-project/pkg/kraken"
	tg "course-project/pkg/telegram"

	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	ErrInitConfig     = errors.New("init config error")
	ErrKrakenConfig   = errors.New("kraken config error")
	ErrTelegramConfig = errors.New("telegram config error")
	ErrDBConfig       = errors.New("db config error")
)

func init() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
}

type AppConfig struct {
	Host     string
	Port     string
	Kraken   kraken.Config
	Telegram tg.TelegramConfig
	DB       DBConfig
}

type DBConfig struct {
	Server   string
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func (d *DBConfig) DNS() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=%v",
		d.Server,
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.SSLMode,
	)
}

func NewAppConfig() (*AppConfig, error) {
	// Init
	err := viper.ReadInConfig()
	if err != nil {
		return nil, ErrInitConfig
	}

	err = godotenv.Load()
	if err != nil {
		return nil, ErrInitConfig
	}

	// Kraken config
	APIKey := os.Getenv("API_KEY")
	if APIKey == "" {
		return nil, ErrKrakenConfig
	}

	APISecret := os.Getenv("API_SECRET")
	if APISecret == "" {
		return nil, ErrKrakenConfig
	}

	krakenConfig := kraken.Config{
		BaseURI:      viper.GetString("kraken.base_uri"),
		BaseDemoURI:  viper.GetString("kraken.base_demo_uri"),
		WSEndpoint:   viper.GetString("kraken.ws_endpoint"),
		RestEndpoint: viper.GetString("kraken.rest_endpoint"),
		APIKey:       APIKey,
		APISecret:    APISecret,
	}

	// Telegram config
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, ErrTelegramConfig
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		return nil, ErrTelegramConfig
	}

	telegramConfig := tg.TelegramConfig{
		Token:  token,
		ChatID: chatID,
	}

	// DB config
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return nil, ErrDBConfig
	}

	dbConfig := DBConfig{
		Server:   viper.GetString("db.server"),
		User:     viper.GetString("db.user"),
		Password: password,
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	// App config
	appConfig := AppConfig{
		Host:     viper.GetString("host"),
		Port:     viper.GetString("port"),
		Kraken:   krakenConfig,
		Telegram: telegramConfig,
		DB:       dbConfig,
	}

	return &appConfig, nil
}
