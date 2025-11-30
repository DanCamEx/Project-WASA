package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

// WebAPIConfiguration describes the web API configuration. This structure is automatically parsed by
// loadConfiguration() and values from flags, environment variable or configuration file will be loaded.
type WebAPIConfiguration struct {
	Config struct {
		Path string `conf:"default:/conf/config.yml"`
	}
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000"`
		DebugHost       string        `conf:"default:0.0.0.0:4000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	Debug bool
	DB    struct {
		Filename string `conf:"default:./database.db"`
	}
}

// loadConfiguration creates a WebAPIConfiguration starting from flags, environment variables and configuration file.
// It works by loading environment variables first, then updating the config using command line flags, finally loading the
// configuration file (specified in WebAPIConfiguration.Config.Path).
// So, CLI parameters will override the environment, and the configuration file will override everything.
// Note that the configuration file can be specified only via CLI or environment variable.
func loadConfiguration() (WebAPIConfiguration, error) {
	var cfg WebAPIConfiguration

	// Try to load configuration from environment variables and command line switches
	if err := conf.Parse(os.Args[1:], "CFG", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, err := conf.Usage("CFG", &cfg)
			if err != nil {
				return cfg, fmt.Errorf("generating config usage: %w", err)
			}
			fmt.Println(usage) //nolint:forbidding
			return cfg, conf.ErrHelpWanted
		}
		return cfg, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}
