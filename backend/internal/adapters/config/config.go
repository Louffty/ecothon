package config

import (
	"fmt"
	"github.com/Louffty/green-code-moscow/internal/domain/entities"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Config struct {
	YandexLoginConfig oauth2.Config
	Database          *gorm.DB
	ListenPort        string
	Logging           bool
}

var AppConfig Config

func VKConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "52328706",
		ClientSecret: "1SwfVpuHQ3UGIkIFpClx",
		RedirectURL:  "https://nothypeproduction.space/api/v1/oauth/vk_callback",
		Scopes:       []string{"email"},
		Endpoint:     vk.Endpoint,
	}
}

func GetConfig() *Config {
	var appLogging bool
	logging := os.Getenv("LOGGING")
	switch logging {
	case "true":
		appLogging = true
	case "false":
		appLogging = false
	default:
		appLogging = false
	}

	var gormConfig *gorm.Config
	if appLogging {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
		gormConfig = &gorm.Config{
			Logger: newLogger,
		}
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=bizkit-database port=5432 sslmode=disable TimeZone=GMT",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	database, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	} else {
		log.Println("Успешно подключились к базе данных")
	}

	err = database.AutoMigrate(
		&entities.User{},
		&entities.Question{},
		&entities.Answer{},
		&entities.Conference{},
		&entities.Event{},
		&entities.Admin{},
		&entities.EventsUser{},
	)

	if err != nil {
		log.Panic(err)
	}

	port := os.Getenv("LISTEN_PORT")

	return &Config{
		Database:   database,
		ListenPort: port,
		Logging:    appLogging,
	}
}
