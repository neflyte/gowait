package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"gowait/internal/logger"
	"gowait/internal/types"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"time"
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

	LogFormatText = "text"
	LogFormatJSON = "json"
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

type AppConfig struct {
	ConfigSource   string        `yaml:"-" json:"-"`
	ConfigFilename string        `yaml:"-" json:"-"`
	Url            url.URL       `yaml:"url" json:"url"`
	RetryDelay     time.Duration `yaml:"retryDelay" json:"retryDelay"`
	RetryLimit     int           `yaml:"retryLimit" json:"retryLimit"`
	SecretSource   string        `yaml:"secretSource" json:"secretSource"`
	SecretFilename string        `yaml:"secretFilename" json:"secretFilename"`
	Secret         string        `yaml:"-" json:"-"`
	LogFormat      string        `yaml:"logFormat" json:"logFormat"`
}

func ReadEnvironmentVariables(cm types.ConfigMap) {
	log := logger.WithField("function", "ReadEnvironmentVariables")
	for envVar, mapKey := range EnvironmentVarMap {
		val, ok := os.LookupEnv(envVar)
		if ok {
			log.Debugf("setting ConfigMap key %s = %s", mapKey, val)
			cm.SetValue(mapKey, val)
		}
	}
}

func (ac *AppConfig) LoadFromConfigMap(cm types.ConfigMap) error {
	log := logger.WithField("function", "LoadFromConfigMap")
	// url
	rawUrl := cm.GetString(KeyURL)
	log.Debugf("rawUrl = %s", rawUrl)
	urlPtr, err := url.Parse(rawUrl)
	if err != nil {
		log.Warnf("unable to parse url from config: %s", err)
		ac.Url = url.URL{}
	} else {
		ac.Url = *urlPtr
	}
	// retryDelay
	retryDuration, err := time.ParseDuration(cm.GetString(KeyRetryDelay))
	if err != nil {
		log.Warnf("unable to parse retryDelay from config: %s; defaulting to %s", err, RetryDelayDefault.String())
		retryDuration = RetryDelayDefault
	}
	ac.RetryDelay = retryDuration
	// retryLimit
	limit, err := strconv.Atoi(cm.GetString(KeyRetryLimit))
	if err != nil {
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

func (ac *AppConfig) LoadFromYAML(fileName string) error {
	log := logger.WithField("function", "LoadFromYAML")
	log.Infof("reading YAML from file %s", fileName)
	rawYaml, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("unable to read file %s: %s", fileName, err)
		return err
	}
	log.Infof("umarshaling YAML")
	err = yaml.Unmarshal(rawYaml, ac)
	if err != nil {
		log.Errorf("unable to unmarshal YAML from file %s: %s", fileName, err)
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
	err = json.Unmarshal(rawJson, ac)
	if err != nil {
		log.Errorf("unable to unmarshal JSON from file %s: %s", fileName, err)
		return err
	}
	return nil
}
func (ac *AppConfig) LoadSecret() {
	log := logger.WithField("function", "LoadSecret")
	switch ac.SecretSource {
	case SecretSourceEnv:
		secretVal, ok := os.LookupEnv(EnvSecret)
		if ok {
			log.Debugf("setting Secret from source %s", SecretSourceEnv)
			ac.Secret = secretVal
		}
	case SecretSourceFile:
		rawSecret, err := ioutil.ReadFile(ac.SecretFilename)
		if err != nil {
			log.Errorf("error reading secret from file %s: %s", ac.SecretFilename, err)
		} else {
			log.Debugf("setting Secret from file %s", ac.SecretFilename)
			ac.Secret = string(rawSecret)
		}
	}
}
