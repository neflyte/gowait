package waiter

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/neflyte/gowait/config"
	"github.com/neflyte/gowait/internal/logger"
	"github.com/neflyte/gowait/internal/utils"
	"net/url"
	"time"
)

const (
	SQLDriverName = "postgres"
)

type postgresWaiter struct {
	urlString  string
	attempts   int
	retryDelay time.Duration
	ticker     *time.Ticker
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
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "PostgresWaiter",
		"function": "Wait",
	})
	success := false
	startTime := time.Now()
	log.Infof("Using retry delay of %s", retryDelay.String())
	pg.ticker = time.NewTicker(retryDelay)
	pg.retryDelay = retryDelay
	pg.urlString = url.String()
	urlStr := utils.SanitizedURLString(url)
	pg.attempts = 0
	for pg.attempts < retryLimit {
		log.Infof("[%d/%d] Connecting to '%s'", pg.attempts+1, retryLimit, urlStr)
		err := pg.connectOnce()
		pg.attempts++ // no matter what happens, we made an attempt
		if err != nil {
			if pg.attempts >= retryLimit {
				log.Errorf("Connect error: %s; retry limit reached; giving up...", err)
				break
			}
			log.Errorf("Connect error: %s; delaying until next retry", err)
			pg.delayOnce()
			continue
		}
		// we're good
		log.Infof("Successfully connected to '%s' after %d of %d attempts; elapsed time: %s", urlStr, pg.attempts, retryLimit, time.Now().Sub(startTime).String())
		success = true
		break
	}
	if !success {
		errStr := fmt.Sprintf("Unable to connect to '%s' after %d attempts; elapsed time: %s", urlStr, pg.attempts, time.Now().Sub(startTime).String())
		log.Errorf(errStr)
		return errors.New(errStr)
	}
	return nil
}

func (pg *postgresWaiter) connectOnce() error {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "PostgresWaiter",
		"function": "connectOnce",
	})
	db, err := sql.Open(SQLDriverName, pg.urlString)
	if err != nil {
		log.Errorf("error opening database connection: %s", err)
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Errorf("error closing database: %s", err)
		}
	}()
	// ping the DB
	err = db.Ping()
	if err != nil {
		log.Errorf("error pinging database: %s", err)
		return err
	}
	// we're good
	return nil
}

func (pg *postgresWaiter) delayOnce() {
	log := logger.WithFields(map[string]interface{}{
		"waiter":   "PostgresWaiter",
		"function": "delayOnce",
	})
	log.Infof("delaying %s until next attempt", pg.retryDelay.String())
	<-pg.ticker.C
}
