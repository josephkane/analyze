package analyze

import (
	"github.com/pkg/errors"

	"github.com/supergiant/analyze/pkg/storage/etcd"

	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/logger"
	"github.com/supergiant/analyze/pkg/plugin"
)

// Config  struct represents configuration of robot service
type Config struct {
	Logging logger.Config `mapstructure:"logging"`
	API     api.Config    `mapstructure:"api"`
	Plugin  plugin.Config `mapstructure:"plugin"`
	ETCD    etcd.Config   `mapstructure:"etcd"`
}

// Validate checks configuration instance for correctness
func (c *Config) Validate() error {
	if err := c.Logging.Validate(); err != nil {
		return err
	}

	if err := c.Plugin.Validate(); err != nil {
		return err
	}

	if len(c.ETCD.Endpoints) == 0 {
		return errors.New("etcd endpoints where not configured")
	}

	return nil
}
