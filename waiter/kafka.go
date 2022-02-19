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
	"github.com/sirupsen/logrus"
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
	log := logger.Function("Wait").Field("waiter", "KafkaWaiter")
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
	log.Infof("Using retry delay of %s", retryDelay.String())
	kw.ticker = time.NewTicker(retryDelay)
	kw.attempts = 0
	for kw.attempts < retryLimit {
		log.Infof("[%d/%d] Connecting to %#v", kw.attempts+1, retryLimit, kw.brokers)
		err := kw.connectOnce()
		kw.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if kw.attempts >= retryLimit {
				log.Errorf("Connect error: %s; retry limit reached; giving up...", err)
				break
			}
			log.Errorf("Connect error: %s; delaying until next retry", err)
			kw.delayOnce()
			continue
		}
		// we're good
		log.Infof("Successfully connected to %#v after %d of %d attempts; elapsed time: %s", kw.brokers, kw.attempts, retryLimit, time.Since(startTime).String())
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to %#v after %d attempts; elapsed time: %s", kw.brokers, kw.attempts, time.Since(startTime).String())
		log.Errorf(errStr)
		return errors.New(errStr)
	}
	return nil

}

func (kw *kafkaWaiter) connectOnce() error {
	log := logger.Function("connectOnce").Field("waiter", "KafkaWaiter")
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = "gowait"
	broker := sarama.NewBroker(kw.brokers[0])
	err := broker.Open(saramaConfig)
	if err != nil {
		log.Errorf("error opening connection to broker %s: %s", kw.brokers[0], err)
		return err
	}
	defer func(logger logrus.FieldLogger, brkr *sarama.Broker) {
		if brkr == nil {
			return
		}
		err := brkr.Close()
		if err != nil {
			logger.Errorf("error closing broker connection: %s", err)
		}
	}(log, broker)
	connected, err := broker.Connected()
	if err != nil {
		log.Errorf("broker connection error: %s", err)
		return err
	}
	if !connected {
		log.Errorf("broker is not connected")
		return ErrConnection
	}
	log.Infof("successfully connected to broker %s", kw.brokers[0])
	return nil
}

func (kw *kafkaWaiter) delayOnce() {
	log := logger.Function("delayOnce").Field("waiter", "KafkaWaiter")
	log.Info("delaying until next attempt")
	<-kw.ticker.C
}
