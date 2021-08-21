package config

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Configuration defines the interface for configuration objects
type Configuration interface {
	GetInt(key string) int

	GetString(key string) string

	IsSet(key string) bool

	LoadConfiguration(cfgFile string) error
}

// NewConfiguration returns a new Configuration instance
func NewConfiguration() Configuration {
	return &concreteConfig{
		cfg: viper.New(),
	}
}

// concreteConfig implements the Configuration interface
type concreteConfig struct {
	cfg *viper.Viper
}

func (c *concreteConfig) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

func (c *concreteConfig) GetString(key string) string {
	return c.cfg.GetString(key)
}

func (c *concreteConfig) IsSet(key string) bool {
	return c.cfg.IsSet(key)
}

// LoadConfiguration loads the configuration for the application from different configuration sources
func (c *concreteConfig) LoadConfiguration(cfgFile string) error {
	log.Debug("Reading configuration ...")

	// From the environment
	c.cfg.SetEnvPrefix("PROVISION")
	c.cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.cfg.AutomaticEnv()

	if cfgFile != "" {
		log.Debug(
			fmt.Sprintf(
				"Reading configuration from: %s",
				cfgFile))

		c.cfg.SetConfigFile(cfgFile)
	}

	if err := c.cfg.ReadInConfig(); err != nil {
		log.Fatal(
			fmt.Sprintf(
				"Configuration invalid. Error was %v",
				err))
		return err
	}

	return nil
}
