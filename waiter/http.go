package waiter

import (
	"errors"
	"fmt"
	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/internal/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type httpWaiter struct {
	urlString string
	attempts  int
	ticker    *time.Ticker
}

func NewHTTPWaiter() Waiter {
	return &httpWaiter{
		urlString: "",
		attempts:  0,
		ticker:    time.NewTicker(config.RetryDelayDefault),
	}
}

func (hw *httpWaiter) Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "HTTPWaiter",
		"function": "Wait",
	})
	success := false
	startTime := time.Now()
	log.Infof("Using retry delay of %s", retryDelay.String())
	hw.ticker = time.NewTicker(retryDelay)
	hw.urlString = url.String()
	hw.attempts = 0
	for hw.attempts < retryLimit {
		log.Infof("[%d/%d] Connecting to '%s'", hw.attempts+1, retryLimit, hw.urlString)
		err := hw.connectOnce(url)
		hw.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if hw.attempts >= retryLimit {
				log.Errorf("Connect error: %s; retry limit reached; giving up...", err)
				break
			}
			log.Errorf("Connect error: %s; delaying until next retry", err)
			hw.delayOnce()
			continue
		}
		// we're good
		log.Infof("Successfully connected to '%s' after %d of %d attempts; elapsed time: %s", hw.urlString, hw.attempts, retryLimit, time.Since(startTime).String())
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to '%s' after %d attempts; elapsed time: %s", hw.urlString, hw.attempts, time.Since(startTime).String())
		log.Errorf(errStr)
		return errors.New(errStr)
	}
	return nil
}

func (hw *httpWaiter) connectOnce(httpUrl url.URL) error {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "HTTPWaiter",
		"function": "connectOnce",
	})
	req, err := http.NewRequest(http.MethodGet, httpUrl.String(), nil)
	if err != nil {
		log.Errorf("error creating new request: %s", err)
		return err
	}
	log.Infof("connecting to %s", httpUrl.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("error executing request: %s", err)
		return err
	}
	defer func(logger logrus.FieldLogger, resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("error closing response body: %s", err)
		}
	}(log, res)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Errorf("request error; code: %d, status: %s", res.StatusCode, res.Status)
		return ErrConnection
	}
	log.Infof("successful request; code: %d, status: %s", res.StatusCode, res.Status)
	return nil
}

func (hw *httpWaiter) delayOnce() {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "HTTPWaiter",
		"function": "delayOnce",
	})
	log.Info("delaying until next attempt")
	<-hw.ticker.C
}
