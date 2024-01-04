package waiter

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/lib/logger"
	"github.com/neflyte/gowait/lib/utils"
)

const (
	SQLDriverName = "postgres"
)

type postgresWaiter struct {
	ticker     *time.Ticker
	urlString  string
	attempts   int
	retryDelay time.Duration
}

func NewPostgresWaiter() Waiter {
	return &postgresWaiter{
		urlString:  "",
		attempts:   0,
		retryDelay: config.RetryDelayDefault,
		ticker:     time.NewTicker(config.RetryDelayDefault),
	}
}

func (pg *postgresWaiter) Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	log := logger.Function("Wait").
		Field("waiter", "PostgresWaiter")
	success := false
	startTime := time.Now()
	log.Field("retryDelay", retryDelay.String()).
		Infof("Using retry delay")
	pg.ticker = time.NewTicker(retryDelay)
	pg.retryDelay = retryDelay
	pg.urlString = url.String()
	urlStr := utils.SanitizedURLString(url)
	pg.attempts = 0
	for pg.attempts < retryLimit {
		log.Field("url", urlStr).
			Infof("[%d/%d] Connecting", pg.attempts+1, retryLimit)
		err := pg.connectOnce()
		pg.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if pg.attempts >= retryLimit {
				log.Err(err).
					Error("Connect error: retry limit reached; giving up")
				break
			}
			log.Err(err).
				Error("Connect error; delaying until next retry")
			pg.delayOnce()
			continue
		}
		// we're good
		log.Fields(map[string]interface{}{
			"url":         urlStr,
			"attempts":    pg.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Infof("Successfully connected")
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to '%s' after %d attempts; elapsed time: %s", urlStr, pg.attempts, time.Since(startTime).String())
		log.Fields(map[string]interface{}{
			"url":         urlStr,
			"attempts":    pg.attempts,
			"retryLimit":  retryLimit,
			"elapsedTime": time.Since(startTime).String(),
		}).
			Error("Unable to connect")
		return errors.New(errStr)
	}
	return nil
}

func (pg *postgresWaiter) connectOnce() error {
	log := logger.Function("connectOnce").
		Field("waiter", "PostgresWaiter")
	db, err := sql.Open(SQLDriverName, pg.urlString)
	if err != nil {
		log.Err(err).
			Error("error opening database connection")
		return err
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Err(err).
				Errorf("error closing database")
		}
	}()
	// ping the DB
	err = db.Ping()
	if err != nil {
		log.Err(err).
			Error("error pinging database")
		return err
	}
	// we're good
	return nil
}

func (pg *postgresWaiter) delayOnce() {
	log := logger.Function("delayOnce").
		Field("waiter", "PostgresWaiter")
	log.Field("delay", pg.retryDelay.String()).
		Info("delaying until next attempt")
	<-pg.ticker.C
}
