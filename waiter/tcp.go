package waiter

import (
	"errors"
	"fmt"
	"gowait/config"
	"gowait/internal/logger"
	"net"
	"net/url"
	"time"
)

type tcpWaiter struct {
	urlString string
	attempts  int
	ticker    *time.Ticker
}

func NewTCPWaiter() Waiter {
	return &tcpWaiter{
		urlString: "",
		attempts:  0,
		ticker:    time.NewTicker(config.RetryDelayDefault),
	}
}

func (tw *tcpWaiter) Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "TCPWaiter",
		"function": "Wait",
	})
	success := false
	startTime := time.Now()
	log.Infof("Using retry delay of %s", retryDelay.String())
	tw.ticker = time.NewTicker(retryDelay)
	tw.urlString = url.String()
	tw.attempts = 0
	for tw.attempts < retryLimit {
		log.Infof("[%d/%d] Connecting to '%s'", tw.attempts+1, retryLimit, tw.urlString)
		err := tw.connectOnce(url.Host)
		tw.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if tw.attempts >= retryLimit {
				log.Errorf("Connect error: %s; retry limit reached; giving up...", err)
				break
			}
			log.Errorf("Connect error: %s; delaying until next retry", err)
			tw.delayOnce()
			continue
		}
		// we're good
		log.Infof("Successfully connected to '%s' after %d of %d attempts; elapsed time: %s", tw.urlString, tw.attempts, retryLimit, time.Now().Sub(startTime).String())
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to '%s' after %d attempts; elapsed time: %s", tw.urlString, tw.attempts, time.Now().Sub(startTime).String())
		log.Errorf(errStr)
		return errors.New(errStr)
	}
	return nil
}

func (tw *tcpWaiter) connectOnce(host string) error {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "TCPWaiter",
		"function": "connectOnce",
	})
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Errorf("unable to connect to tcp address %s: %s", host, err)
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Errorf("error closing tcp connection: %s", err)
		}
	}()
	return nil
}

func (tw *tcpWaiter) delayOnce() {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "TCPWaiter",
		"function": "delayOnce",
	})
	log.Info("delaying until next attempt")
	<-tw.ticker.C
}
