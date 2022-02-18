package config

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/neflyte/configmap"
	"github.com/neflyte/gowait/internal/logger"
	"gopkg.in/yaml.v3"
)

const (
	RetryLimitDefault   = 5
	RetryDelayDefault   = 10 * time.Second
	ConfSourceDefault   = ConfSourceEnv
	SecretSourceDefault = SecretSourceEnv

	ConfSourceEnv  = "env"
	ConfSourceYAML = "yaml"
	ConfSourceJSON = "json"

	KeyRetryDelay     = "retryDelay"
	KeyRetryLimit     = "retryLimit"
	KeyURL            = "url"
	KeySecretSource   = "secretSource"
	KeySecretFilename = "secretFilename"
	KeyLogFormat      = "logFormat"

	EnvRetryDelay     = "GOWAIT_RETRY_DELAY"
	EnvRetryLimit     = "GOWAIT_RETRY_LIMIT"
	EnvURL            = "GOWAIT_URL"
	EnvSecretSource   = "GOWAIT_SECRET_SOURCE"
	EnvSecretFilename = "GOWAIT_SECRET_FILENAME"
	EnvSecret         = "GOWAIT_SECRET"
	EnvLogFormat      = "GOWAIT_LOG_FORMAT"

	SecretSourceEnv  = "env"
	SecretSourceFile = "file"
)

var (
	EnvironmentVarMap = map[string]string{
		EnvRetryDelay:     KeyRetryDelay,
		EnvRetryLimit:     KeyRetryLimit,
		EnvURL:            KeyURL,
		EnvSecretSource:   KeySecretSource,
		EnvSecretFilename: KeySecretFilename,
		EnvLogFormat:      KeyLogFormat,
	}
)

// AppConfig represents the struct of application configuration info
type AppConfig struct {
	// internal fields
	ConfigSource   string `yaml:"-" json:"-"`
	ConfigFilename string `yaml:"-" json:"-"`
	Secret         string `yaml:"-" json:"-"`
	// data from the configuration file
	Url            url.URL       `yaml:"url" json:"url"`
	RetryDelay     time.Duration `yaml:"retryDelay" json:"retryDelay"`
	RetryLimit     int           `yaml:"retryLimit" json:"retryLimit"`
	SecretSource   string        `yaml:"secretSource" json:"secretSource"`
	SecretFilename string        `yaml:"secretFilename" json:"secretFilename"`
	LogFormat      string        `yaml:"logFormat" json:"logFormat"`
}

// AppConfigFile represents the configuration struct in a flat file
type AppConfigFile struct {
	Url            string `yaml:"url" json:"url"`
	RetryDelay     string `yaml:"retryDelay" json:"retryDelay"`
	RetryLimit     int    `yaml:"retryLimit" json:"retryLimit"`
	SecretSource   string `yaml:"secretSource" json:"secretSource"`
	SecretFilename string `yaml:"secretFilename" json:"secretFilename"`
	LogFormat      string `yaml:"logFormat" json:"logFormat"`
}

func ReadEnvironmentVariables(cm configmap.ConfigMap) {
	log := logger.WithField("function", "ReadEnvironmentVariables")
	for envVar, mapKey := range EnvironmentVarMap {
		val, ok := os.LookupEnv(envVar)
		if ok {
			log.Debugf("setting ConfigMap key %s = %s", mapKey, val)
			cm.Set(mapKey, val)
		}
	}
}

func (ac *AppConfig) LoadFromConfigMap(cm configmap.ConfigMap) error {
	log := logger.WithField("function", "LoadFromConfigMap")
	// url
	ac.Url = url.URL{}
	rawUrl := cm.GetString(KeyURL)
	log.Debugf("rawUrl = %s", rawUrl)
	if rawUrl != "" {
		urlPtr, err := url.Parse(rawUrl)
		if err != nil {
			log.Warnf("unable to parse url from config: %s", err)
		} else {
			ac.Url = *urlPtr
		}
	}
	// retryDelay
	retryDuration, err := time.ParseDuration(cm.GetString(KeyRetryDelay))
	if err != nil && cm.GetString(KeyRetryDelay) != "" {
		log.Warnf("unable to parse retryDelay from config: %s; defaulting to %s", err, RetryDelayDefault.String())
		retryDuration = RetryDelayDefault
	}
	ac.RetryDelay = retryDuration
	// retryLimit
	limit, err := strconv.Atoi(cm.GetString(KeyRetryLimit))
	if err != nil && cm.GetString(KeyRetryLimit) != "" {
		log.Warnf("unable to parse retryLimit from config: %s; defaulting to %d", err, RetryLimitDefault)
		limit = RetryLimitDefault
	}
	ac.RetryLimit = limit
	// secretSource
	secSrc := cm.GetString(KeySecretSource)
	if secSrc == "" {
		log.Infof("no secret source specified; defaulting to %s", SecretSourceDefault)
		secSrc = SecretSourceDefault
	}
	ac.SecretSource = secSrc
	// secretFilename
	ac.SecretFilename = cm.GetString(KeySecretFilename)
	// logFormat
	ac.LogFormat = cm.GetString(KeyLogFormat)
	// done.
	return nil
}

func (ac *AppConfig) PopulateFromAppConfigFile(fileCfg *AppConfigFile) error {
	log := logger.WithField("function", "PopulateFromAppConfigFile")
	if fileCfg == nil {
		log.Warnf("nil AppConfigFile; nothing to do")
		return nil
	}
	// copy the data over to ac
	waitUrl, err := url.Parse(fileCfg.Url)
	if err != nil {
		log.Errorf("error parsing URL: %s", err)
		return err
	}
	ac.Url = *waitUrl
	ac.RetryDelay, err = time.ParseDuration(fileCfg.RetryDelay)
	if err != nil {
		log.Errorf("error parsing RetryDelay: %s", err)
		return err
	}
	ac.RetryLimit = fileCfg.RetryLimit
	ac.SecretSource = fileCfg.SecretSource
	ac.SecretFilename = fileCfg.SecretFilename
	ac.LogFormat = fileCfg.LogFormat
	return nil
}

func (ac *AppConfig) LoadFromYAML(fileName string) error {
	log := logger.WithField("function", "LoadFromYAML")
	log.Infof("reading YAML from file %s", fileName)
	rawYaml, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("unable to read file %s: %s", fileName, err)
		return err
	}
	log.Infof("umarshaling YAML")
	fileCfg := &AppConfigFile{}
	err = yaml.Unmarshal(rawYaml, fileCfg)
	if err != nil {
		log.Errorf("unable to unmarshal YAML from file %s: %s", fileName, err)
		return err
	}
	// copy the data over to ac
	err = ac.PopulateFromAppConfigFile(fileCfg)
	if err != nil {
		log.Errorf("error populating config from configfile: %s", err)
		return err
	}
	return nil
}

func (ac *AppConfig) LoadFromJSON(fileName string) error {
	log := logger.WithField("function", "LoadFromJSON")
	log.Infof("reading JSON from file %s", fileName)
	rawJson, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("unable to read file %s: %s", fileName, err)
		return err
	}
	log.Infof("umarshaling JSON")
	fileCfg := &AppConfigFile{}
	err = json.Unmarshal(rawJson, fileCfg)
	if err != nil {
		log.Errorf("unable to unmarshal JSON from file %s: %s", fileName, err)
		return err
	}
	// copy the data over to ac
	err = ac.PopulateFromAppConfigFile(fileCfg)
	if err != nil {
		log.Errorf("error populating config from configfile: %s", err)
		return err
	}
	return nil
}

func (ac *AppConfig) LoadSecret() {
	log := logger.WithField("function", "LoadSecret")
	switch ac.SecretSource {
	case SecretSourceEnv:
		secretVal, ok := os.LookupEnv(EnvSecret)
		if ok && secretVal != "" {
			log.Debugf("setting Secret from source %s", SecretSourceEnv)
			ac.Secret = secretVal
		}
	case SecretSourceFile:
		rawSecret, err := ioutil.ReadFile(ac.SecretFilename)
		if err != nil {
			log.Errorf("error reading secret from file %s: %s", ac.SecretFilename, err)
		} else {
			if string(rawSecret) != "" {
				log.Debugf("setting Secret from file %s", ac.SecretFilename)
				ac.Secret = string(rawSecret)
			}
		}
	}
}
