package waiter

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/internal/logger"
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
	log := logger.Function("Wait").
		Field("waiter", "TCPWaiter")
	success := false
	startTime := time.Now()
	log.Field("retryDelay", retryDelay.String()).
		Info("Using retry delay")
	tw.ticker = time.NewTicker(retryDelay)
	tw.urlString = url.String()
	tw.attempts = 0
	for tw.attempts < retryLimit {
		log.Field("url", tw.urlString).
			Infof("[%d/%d] Connecting", tw.attempts+1, retryLimit)
		err := tw.connectOnce(url.Host)
		tw.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if tw.attempts >= retryLimit {
				log.Err(err).
					Error("Connect error: retry limit reached; giving up")
				break
			}
			log.Err(err).
				Error("Connect error; delaying until next retry")
			tw.delayOnce()
			continue
		}
		// we're good
		log.Fields(map[string]interface{}{
			"url":         tw.urlString,
			"attempts":    tw.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Info("Successfully connected")
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to '%s' after %d attempts; elapsed time: %s", tw.urlString, tw.attempts, time.Since(startTime).String())
		log.Fields(map[string]interface{}{
			"url":         tw.urlString,
			"attempts":    tw.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Errorf("Unable to connect")
		return errors.New(errStr)
	}
	return nil
}

func (tw *tcpWaiter) connectOnce(host string) error {
	log := logger.Function("connectOnce").
		Field("waiter", "TCPWaiter")
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Err(err).
			Field("host", host).
			Error("unable to connect to tcp address")
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Err(err).
				Error("error closing tcp connection")
		}
	}()
	return nil
}

func (tw *tcpWaiter) delayOnce() {
	log := logger.Function("delayOnce").
		Field("waiter", "TCPWaiter")
	log.Info("delaying until next attempt")
	<-tw.ticker.C
}
