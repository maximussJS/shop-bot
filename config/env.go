package config

import (
	"fmt"
	"os"
	"shop-bot/utils"
	"strconv"
)

func (c *Config) getRequiredString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		c.logger.Error(`Environment variable "` + key + `" not found`)
		panic(fmt.Sprintf(`Environment variable "%s" not found`, key))
	}

	return value
}

func (c *Config) getOptionalInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		c.logger.Warn(fmt.Sprintf(`Environment variable "%s" not found, used default %d`, key, defaultValue))
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)

	utils.PanicIfError(err)

	return valueInt
}

func (c *Config) getOptionalBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)

	if value == "" {
		c.logger.Warn(fmt.Sprintf(`Environment variable "%s" not found, used default %t`, key, defaultValue))
		return defaultValue
	}

	valueBool, err := strconv.ParseBool(value)

	utils.PanicIfError(err)

	return valueBool
}
