package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"shop-bot/constants"
	"shop-bot/internal/logger"
	"time"
)

type IConfig interface {
	AppEnv() constants.AppEnv

	BotToken() string

	ErrorStackTraceSizeInKb() int

	PostgresDSN() string
	RunMigrations() bool

	RequestTimeout() time.Duration
}

type configDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
}

type Config struct {
	logger logger.ILogger

	appEnv constants.AppEnv

	botToken string

	errorStackTraceSizeInKb int

	postgresDsn   string
	runMigrations bool

	requestTimeoutInSeconds int
}

func NewConfig(deps configDependencies) *Config {
	_logger := deps.Logger

	godotenv.Load() // ignore error, because in deployment we pass all env variables via docker run command

	config := &Config{
		logger: _logger,
	}

	appEnv := config.getRequiredString("APP_ENV")

	switch appEnv {
	case string(constants.DevelopmentEnv):
		config.appEnv = constants.DevelopmentEnv
	case string(constants.ProductionEnv):
		config.appEnv = constants.ProductionEnv
	default:
		panic(fmt.Sprintf("Invalid APP_ENV value: %s. Supported values: %s, %s", appEnv, constants.DevelopmentEnv, constants.ProductionEnv))
	}

	config.botToken = config.getRequiredString("BOT_TOKEN")
	config.postgresDsn = config.getRequiredString("POSTGRES_DSN")
	config.runMigrations = config.getOptionalBool("RUN_MIGRATIONS", false)
	config.errorStackTraceSizeInKb = config.getOptionalInt("ERROR_STACK_TRACE_SIZE_IN_KB", 4)
	config.requestTimeoutInSeconds = config.getOptionalInt("REQUEST_TIMEOUT_IN_SECONDS", 15)

	return config
}

func (c *Config) AppEnv() constants.AppEnv {
	return c.appEnv
}

func (c *Config) BotToken() string {
	return c.botToken
}

func (c *Config) ErrorStackTraceSizeInKb() int {
	return c.errorStackTraceSizeInKb
}

func (c *Config) PostgresDSN() string {
	return c.postgresDsn
}

func (c *Config) RunMigrations() bool {
	return c.runMigrations
}

func (c *Config) RequestTimeout() time.Duration {
	return time.Duration(c.requestTimeoutInSeconds) * time.Second
}
