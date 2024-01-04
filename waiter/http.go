package waiter

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/lib/logger"
)

type httpWaiter struct {
	ticker    *time.Ticker
	urlString string
	attempts  int
}

func NewHTTPWaiter() Waiter {
	return &httpWaiter{
		urlString: "",
		attempts:  0,
		ticker:    time.NewTicker(config.RetryDelayDefault),
	}
}

func (hw *httpWaiter) Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	log := logger.Function("Wait").
		Field("waiter", "HTTPWaiter")
	success := false
	startTime := time.Now()
	log.Field("delay", retryDelay.String).
		Info("Using retry delay")
	hw.ticker = time.NewTicker(retryDelay)
	hw.urlString = url.String()
	hw.attempts = 0
	for hw.attempts < retryLimit {
		log.Field("url", hw.urlString).
			Infof("[%d/%d] Connecting", hw.attempts+1, retryLimit)
		err := hw.connectOnce(url)
		hw.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if hw.attempts >= retryLimit {
				log.Err(err).
					Error("Connect error: retry limit reached; giving up")
				break
			}
			log.Err(err).
				Errorf("Connect error; delaying until next retry")
			hw.delayOnce()
			continue
		}
		// we're good
		log.Fields(map[string]interface{}{
			"url":         hw.urlString,
			"attempts":    hw.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Info("Successfully connected")
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to '%s' after %d attempts; elapsed time: %s", hw.urlString, hw.attempts, time.Since(startTime).String())
		log.Fields(map[string]interface{}{
			"url":         hw.urlString,
			"attempts":    hw.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Errorf("Unable to connect")
		return errors.New(errStr)
	}
	return nil
}

func (hw *httpWaiter) connectOnce(httpUrl url.URL) error {
	log := logger.Function("connectOnce").
		Field("waiter", "HTTPWaiter")
	req, err := http.NewRequest(http.MethodGet, httpUrl.String(), nil)
	if err != nil {
		log.Err(err).
			Error("error creating new request")
		return err
	}
	log.Field("httpUrl", httpUrl.String()).
		Info("connecting")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Err(err).
			Error("error executing request")
		return err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Err(err).
				Error("error closing response body")
		}
	}()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Errorf("request error; code: %d, status: %s", res.StatusCode, res.Status)
		return ErrConnection
	}
	log.Fields(map[string]interface{}{
		"statusCode": res.StatusCode,
		"status":     res.Status,
	}).
		Info("successful request")
	return nil
}

func (hw *httpWaiter) delayOnce() {
	log := logger.Function("delayOnce").
		Field("waiter", "HTTPWaiter")
	log.Info("delaying until next attempt")
	<-hw.ticker.C
}
