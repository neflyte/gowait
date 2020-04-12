package main

import (
	"flag"
	"github.com/neflyte/configmap"
	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/internal/logger"
	"github.com/neflyte/gowait/internal/utils"
	"github.com/neflyte/gowait/waiter"
	"github.com/sirupsen/logrus"
	"net/url"
)

const (
	AppVersion = "0.1.3"
)

var (
	cm  configmap.ConfigMap
	cfg *config.AppConfig
)

func init() {
	cfg = new(config.AppConfig)
	flag.StringVar(&cfg.ConfigSource, "configSource", config.ConfSourceDefault, "where to read the app config from; 'env' = environment vars, 'yaml' = yaml file, 'json' = json file")
	flag.StringVar(&cfg.ConfigSource, "c", config.ConfSourceDefault, "where to read the app config from; 'env' = environment vars, 'yaml' = yaml file, 'json' = json file (shorthand)")
	flag.StringVar(&cfg.ConfigFilename, "configFile", "", "path/name of file to read app config from")
	flag.StringVar(&cfg.ConfigFilename, "f", "", "path/name of file to read app config from (shorthand)")
	flag.Parse()
}

func main() {
	log := logger.WithField("function", "main")
	log.Warnf("gowait v%s - service readiness waiter", AppVersion)

	log.Info("Load configuration")
	cm = configmap.New()
	switch cfg.ConfigSource {
	case config.ConfSourceEnv:
		log.Info("read environment variables")
		config.ReadEnvironmentVariables(cm)
		log.Info("initialize configuration")
		err := cfg.LoadFromConfigMap(cm)
		if err != nil {
			log.Fatalf("unable to load configuration: %s; aborting...", err)
		}
	case config.ConfSourceJSON:
		if cfg.ConfigFilename == "" {
			log.Fatal("config source set to JSON but no config file specified; aborting...")
		}
		log.Infof("initialize configuration from JSON file %s", cfg.ConfigFilename)
		err := cfg.LoadFromJSON(cfg.ConfigFilename)
		if err != nil {
			log.Fatalf("unable to load configuration from JSON file %s: %s; aborting...", cfg.ConfigFilename, err)
		}
	case config.ConfSourceYAML:
		if cfg.ConfigFilename == "" {
			log.Fatal("config source set to YAML but no config file specified; aborting...")
		}
		log.Infof("initialize configuration from YAML file %s", cfg.ConfigFilename)
		err := cfg.LoadFromYAML(cfg.ConfigFilename)
		if err != nil {
			log.Fatalf("unable to load configuration from YAML file %s: %s; aborting...", cfg.ConfigFilename, err)
		}
	}

	// reconfigure logging
	switch cfg.LogFormat {
	case config.LogFormatText:
		logger.Logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		})
	case config.LogFormatJSON:
		logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// load secret
	log.Info("Load secret")
	cfg.LoadSecret()

	// take a copy of the sanitized URL as a string
	urlStr := utils.SanitizedURLString(cfg.Url)

	// add secret to URL if it's non-empty
	if cfg.Url.User != nil {
		log.Debug("adding secret to URL Userinfo")
		if cfg.Secret != "" {
			cfg.Url.User = url.UserPassword(cfg.Url.User.Username(), cfg.Secret)
		} else {
			cfg.Url.User = url.User(cfg.Url.User.Username())
		}
	}

	// go wait!
	log.Infof("Starting to wait for '%s', making at most %d attempts with a %s delay between each", urlStr, cfg.RetryLimit, cfg.RetryDelay.String())
	err := waiter.Wait(cfg.Url, cfg.RetryDelay, cfg.RetryLimit)
	if err != nil {
		log.Fatalf("Error waiting for %s: %s; aborting...", urlStr, err)
	}
	log.Infof("Successfully waited for %s; done.", urlStr)
}
