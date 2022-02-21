package waiter

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/internal/logger"
)

// url: kafka://broker1:port/?brokers=broker2:port,broker3:port...

type kafkaWaiter struct {
	brokers  []string
	attempts int
	ticker   *time.Ticker
}

func NewKafkaWaiter() Waiter {
	return &kafkaWaiter{
		brokers:  make([]string, 0),
		attempts: 0,
		ticker:   time.NewTicker(config.RetryDelayDefault),
	}
}

func (kw *kafkaWaiter) Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	log := logger.Function("Wait").
		Field("waiter", "KafkaWaiter")
	// start with the url hostname
	kw.brokers = append(kw.brokers, url.Host)
	// add any extra brokers
	urlBrokers := url.Query().Get("urlBrokers")
	if len(urlBrokers) > 0 {
		toks := strings.Split(urlBrokers, ",")
		for _, tok := range toks {
			kw.brokers = append(kw.brokers, strings.TrimSpace(tok))
		}
	}
	success := false
	startTime := time.Now()
	log.Field("retryDelay", retryDelay.String()).
		Info("Using retry delay")
	kw.ticker = time.NewTicker(retryDelay)
	kw.attempts = 0
	for kw.attempts < retryLimit {
		log.Field("brokers", fmt.Sprintf("%#v", kw.brokers)).
			Infof("[%d/%d] Connecting", kw.attempts+1, retryLimit)
		err := kw.connectOnce()
		kw.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if kw.attempts >= retryLimit {
				log.Err(err).
					Error("Connect error: retry limit reached; giving up")
				break
			}
			log.Err(err).
				Error("Connect error; delaying until next retry")
			kw.delayOnce()
			continue
		}
		// we're good
		log.Fields(map[string]interface{}{
			"brokers":     fmt.Sprintf("%#v", kw.brokers),
			"attempts":    kw.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Info("Successfully connected")
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to %#v after %d attempts; elapsed time: %s", kw.brokers, kw.attempts, time.Since(startTime).String())
		log.Fields(map[string]interface{}{
			"brokers":     fmt.Sprintf("%#v", kw.brokers),
			"attempts":    kw.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Errorf("Unable to connect")
		return errors.New(errStr)
	}
	return nil
}

func (kw *kafkaWaiter) connectOnce() error {
	log := logger.Function("connectOnce").
		Field("waiter", "KafkaWaiter")
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = "gowait"
	broker := sarama.NewBroker(kw.brokers[0])
	err := broker.Open(saramaConfig)
	if err != nil {
		log.Err(err).
			Field("broker", kw.brokers[0]).
			Error("error opening connection to broker")
		return err
	}
	defer func() {
		if broker == nil {
			return
		}
		err = broker.Close()
		if err != nil {
			log.Err(err).
				Field("broker", kw.brokers[0]).
				Error("error closing broker connection")
		}
	}()
	connected, err := broker.Connected()
	if err != nil {
		log.Err(err).
			Error("broker connection error")
		return err
	}
	if !connected {
		log.Error("broker is not connected")
		return ErrConnection
	}
	log.Field("broker", kw.brokers[0]).
		Info("successfully connected to broker")
	return nil
}

func (kw *kafkaWaiter) delayOnce() {
	log := logger.Function("delayOnce").
		Field("waiter", "KafkaWaiter")
	log.Info("delaying until next attempt")
	<-kw.ticker.C
}
