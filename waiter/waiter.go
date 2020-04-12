package waiter

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

var (
	ErrConnection = errors.New("connection error")
)

type Waiter interface {
	Wait(url url.URL, retryDelay time.Duration, retryLimit int) error
}

func Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	// select the appropriate waiter
	var waiter Waiter
	switch url.Scheme {
	case "postgres":
		waiter = NewPostgresWaiter()
	case "tcp":
		waiter = NewTCPWaiter()
	case "http":
		waiter = NewHTTPWaiter()
	case "kafka":
		waiter = NewKafkaWaiter()
	default:
		return fmt.Errorf("unknown scheme: %s", url.Scheme)
	}
	// wait For IT!!
	return waiter.Wait(url, retryDelay, retryLimit)
}
