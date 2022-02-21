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
	LogFormatDefault    = logger.LogFormatText
	LogLevelDefault     = logger.LogLevelInfo

	ConfSourceEnv  = "env"
	ConfSourceYAML = "yaml"
	ConfSourceJSON = "json"

	KeyRetryDelay     = "retryDelay"
	KeyRetryLimit     = "retryLimit"
	KeyURL            = "url"
	KeySecretSource   = "secretSource"
	KeySecretFilename = "secretFilename"
	KeyLogFormat      = "logFormat"
	KeyLogLevel       = "logLevel"

	EnvRetryDelay     = "GOWAIT_RETRY_DELAY"
	EnvRetryLimit     = "GOWAIT_RETRY_LIMIT"
	EnvURL            = "GOWAIT_URL"
	EnvSecretSource   = "GOWAIT_SECRET_SOURCE"
	EnvSecretFilename = "GOWAIT_SECRET_FILENAME"
	EnvSecret         = "GOWAIT_SECRET"
	EnvLogFormat      = "GOWAIT_LOG_FORMAT"
	EnvLogLevel       = "GOWAIT_LOG_LEVEL"

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
		EnvLogLevel:       KeyLogLevel,
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
	LogLevel       string        `yaml:"logLevel" json:"logLevel"`
}

// AppConfigFile represents the configuration struct in a flat file
type AppConfigFile struct {
	Url            string `yaml:"url" json:"url"`
	RetryDelay     string `yaml:"retryDelay" json:"retryDelay"`
	RetryLimit     int    `yaml:"retryLimit" json:"retryLimit"`
	SecretSource   string `yaml:"secretSource" json:"secretSource"`
	SecretFilename string `yaml:"secretFilename" json:"secretFilename"`
	LogFormat      string `yaml:"logFormat" json:"logFormat"`
	LogLevel       string `yaml:"logLevel" json:"logLevel"`
}

func ReadEnvironmentVariables(cm configmap.ConfigMap) {
	log := logger.Function("ReadEnvironmentVariables")
	for envVar, mapKey := range EnvironmentVarMap {
		val, ok := os.LookupEnv(envVar)
		if ok {
			log.Fields(map[string]interface{}{
				"key":   mapKey,
				"value": val,
			}).
				Debug("setting ConfigMap key")
			cm.Set(mapKey, val)
		}
	}
}

func (ac *AppConfig) LoadFromConfigMap(cm configmap.ConfigMap) error {
	log := logger.Function("LoadFromConfigMap")
	// url
	ac.Url = url.URL{}
	rawUrl := cm.GetString(KeyURL)
	log.Debugf("rawUrl: %s", rawUrl)
	if rawUrl != "" {
		urlPtr, err := url.Parse(rawUrl)
		if err != nil {
			log.Err(err).
				Warn("unable to parse url from config")
		} else {
			ac.Url = *urlPtr
		}
	}
	// retryDelay
	retryDuration, err := time.ParseDuration(cm.GetString(KeyRetryDelay))
	if err != nil && cm.GetString(KeyRetryDelay) != "" {
		log.Err(err).
			Field("default", RetryDelayDefault.String()).
			Warn("unable to parse retryDelay from config; using default")
		retryDuration = RetryDelayDefault
	}
	ac.RetryDelay = retryDuration
	// retryLimit
	limit, err := strconv.Atoi(cm.GetString(KeyRetryLimit))
	if err != nil && cm.GetString(KeyRetryLimit) != "" {
		log.Err(err).
			Field("default", RetryLimitDefault).
			Warn("unable to parse retryLimit from config; using default")
		limit = RetryLimitDefault
	}
	ac.RetryLimit = limit
	// secretSource
	secSrc := cm.GetString(KeySecretSource)
	if secSrc == "" {
		log.Field("default", SecretSourceDefault).
			Info("no secret source specified; using default")
		secSrc = SecretSourceDefault
	}
	ac.SecretSource = secSrc
	// secretFilename
	ac.SecretFilename = cm.GetString(KeySecretFilename)
	// logFormat
	ac.LogFormat = LogFormatDefault
	if cm.GetString(KeyLogFormat) != "" {
		ac.LogFormat = cm.GetString(KeyLogFormat)
	}
	// logLevel
	ac.LogLevel = LogLevelDefault
	if cm.GetString(KeyLogLevel) != "" {
		ac.LogLevel = cm.GetString(KeyLogLevel)
	}
	// done.
	return nil
}

func (ac *AppConfig) PopulateFromAppConfigFile(fileCfg *AppConfigFile) error {
	log := logger.Function("PopulateFromAppConfigFile")
	if fileCfg == nil {
		log.Warnf("nil AppConfigFile; nothing to do")
		return nil
	}
	// copy the data over to ac
	waitUrl, err := url.Parse(fileCfg.Url)
	if err != nil {
		log.Err(err).
			Field("url", fileCfg.Url).
			Error("error parsing URL")
		return err
	}
	ac.Url = *waitUrl
	ac.RetryDelay, err = time.ParseDuration(fileCfg.RetryDelay)
	if err != nil {
		log.Err(err).
			Field("retryDelay", fileCfg.RetryDelay).
			Error("error parsing RetryDelay")
		return err
	}
	ac.RetryLimit = fileCfg.RetryLimit
	ac.SecretSource = fileCfg.SecretSource
	ac.SecretFilename = fileCfg.SecretFilename
	ac.LogFormat = fileCfg.LogFormat
	ac.LogLevel = fileCfg.LogLevel
	return nil
}

func (ac *AppConfig) LoadFromYAML(fileName string) error {
	log := logger.Function("LoadFromYAML")
	log.Field("file", fileName).
		Debug("reading YAML from file")
	rawYaml, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Err(err).
			Field("file", fileName).
			Error("unable to read file")
		return err
	}
	log.Debug("umarshaling YAML")
	fileCfg := &AppConfigFile{}
	err = yaml.Unmarshal(rawYaml, fileCfg)
	if err != nil {
		log.Err(err).
			Field("file", fileName).
			Error("unable to unmarshal YAML from file")
		return err
	}
	// copy the data over to ac
	err = ac.PopulateFromAppConfigFile(fileCfg)
	if err != nil {
		log.Err(err).
			Error("error populating config from configfile")
		return err
	}
	return nil
}

func (ac *AppConfig) LoadFromJSON(fileName string) error {
	log := logger.Function("LoadFromJSON")
	log.Field("file", fileName).
		Debug("reading JSON from file")
	rawJson, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Err(err).
			Field("file", fileName).
			Error("unable to read file")
		return err
	}
	log.Debug("umarshaling JSON")
	fileCfg := &AppConfigFile{}
	err = json.Unmarshal(rawJson, fileCfg)
	if err != nil {
		log.Err(err).
			Field("file", fileName).
			Error("unable to unmarshal JSON from file")
		return err
	}
	// copy the data over to ac
	err = ac.PopulateFromAppConfigFile(fileCfg)
	if err != nil {
		log.Err(err).
			Error("error populating config from configfile")
		return err
	}
	return nil
}

func (ac *AppConfig) LoadSecret() {
	log := logger.Function("LoadSecret")
	switch ac.SecretSource {
	case SecretSourceEnv:
		secretVal, ok := os.LookupEnv(EnvSecret)
		if ok && secretVal != "" {
			log.Field("source", SecretSourceEnv).
				Debug("setting Secret from source")
			ac.Secret = secretVal
		}
	case SecretSourceFile:
		rawSecret, err := ioutil.ReadFile(ac.SecretFilename)
		if err != nil {
			log.Err(err).
				Field("file", ac.SecretFilename).
				Error("error reading secret from file")
		} else {
			if string(rawSecret) != "" {
				log.Field("file", ac.SecretFilename).
					Debug("setting Secret from file")
				ac.Secret = string(rawSecret)
			}
		}
	}
}
