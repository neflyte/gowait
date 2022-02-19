package main

import (
	"flag"
	"net/url"

	"github.com/neflyte/configmap"
	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/internal/logger"
	"github.com/neflyte/gowait/internal/utils"
	"github.com/neflyte/gowait/waiter"
)

const (
	AppVersion = "v0.1.4"
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
	log := logger.Function("main")
	log.Field("version", AppVersion).
		Warn("gowait - service readiness waiter")

	log.Info("Load configuration")
	cm = configmap.New()
	switch cfg.ConfigSource {
	case config.ConfSourceEnv:
		log.Info("read environment variables")
		config.ReadEnvironmentVariables(cm)
		log.Info("initialize configuration")
		err := cfg.LoadFromConfigMap(cm)
		if err != nil {
			log.Err(err).
				Fatalf("unable to load configuration; aborting")
		}
	case config.ConfSourceJSON:
		if cfg.ConfigFilename == "" {
			log.Fatal("config source set to JSON but no config file specified; aborting...")
		}
		log.Field("configFilename", cfg.ConfigFilename).
			Infof("initialize configuration from JSON file")
		err := cfg.LoadFromJSON(cfg.ConfigFilename)
		if err != nil {
			log.Field("configFilename", cfg.ConfigFilename).
				Err(err).
				Fatal("unable to load configuration from JSON file; aborting")
		}
	case config.ConfSourceYAML:
		if cfg.ConfigFilename == "" {
			log.Fatal("config source set to YAML but no config file specified; aborting...")
		}
		log.Infof("initialize configuration from YAML file %s", cfg.ConfigFilename)
		err := cfg.LoadFromYAML(cfg.ConfigFilename)
		if err != nil {
			log.Field("configFilename", cfg.ConfigFilename).
				Err(err).
				Fatal("unable to load configuration from YAML file; aborting")
		}
	}

	// reconfigure logging
	logger.ConfigureFormat(cfg.LogFormat)
	logger.ConfigureLevel(cfg.LogLevel)

	// do we have a URL to wait for?
	if cfg.Url.String() == "" {
		log.Fatal("no URL was specified; nothing to wait for")
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
	log.Fields(map[string]interface{}{
		"url":        urlStr,
		"maxRetries": cfg.RetryLimit,
		"retryDelay": cfg.RetryDelay.String(),
	}).
		Infof("Starting to wait")
	err := waiter.Wait(cfg.Url, cfg.RetryDelay, cfg.RetryLimit)
	if err != nil {
		log.Field("url", urlStr).
			Err(err).
			Fatal("Error waiting; aborting")
	}
	log.Field("url", urlStr).
		Info("Successfully waited; done.")
}
