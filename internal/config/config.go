package config

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ilyakaznacheev/cleanenv"

	"git.burning.moe/celediel/burning.moe/internal/models"
)

// AppConfig contains data to be accessed across the app.
type AppConfig struct {
	CacheTimer    time.Duration
	ListenPort    uint16
	Logger        *log.Logger
	LogLevel      log.Level
	TemplateCache models.TemplateCache
	UseCache      bool
}

// ConfigDatabase contains data to be loaded from environmental variables
type ConfigDatabase struct {
	CacheTimer time.Duration `env:"CACHETIMER" env-default:"12h" env-description:"How often to automatically regenerate template cache."`
	LogLevel   string        `env:"LOGLEVEL" env-default:"warn" env-description:"Logging level. Default: warn, Possible values: debug info warn error fatal none"`
	Port       uint16        `env:"PORT" env-default:"9001" env-description:"server port"`
	UseCache   bool          `env:"CACHE" env-default:"true" env-description:"Use template cache"`
}

// Initialises the app wide AppConfig, loads values from environment, and set up the Logger
func Initialise() AppConfig {
	app := AppConfig{}

	// Setup logger
	app.Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})

	// load values from config
	if cfg, err := loadConfig(); err == nil {
		app.ListenPort = cfg.Port
		app.UseCache = cfg.UseCache
		app.CacheTimer = cfg.CacheTimer
		app.LogLevel, err = log.ParseLevel(cfg.LogLevel)
		if err != nil {
			app.LogLevel = log.WarnLevel
		}
	} else {
		app.Logger.Print("Failed loading config from environment", "err", err)
	}

	app.Logger.SetLevel(app.LogLevel)
	app.Logger.Debug("Loaded config from environment:",
		"port", app.ListenPort, "useCache", app.UseCache,
		"log_level", app.LogLevel, "cacheTimer", app.CacheTimer,
	)

	return app
}

// loadConfig utilises cleanenv to load config values from the environment
func loadConfig() (ConfigDatabase, error) {
	var cfg ConfigDatabase
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return ConfigDatabase{}, err
	} else {
		return cfg, nil
	}
}
