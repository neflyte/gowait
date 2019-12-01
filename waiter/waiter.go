package waiter

import (
	"fmt"
	"net/url"
	"time"
)

type Waiter interface {
	Wait(url url.URL, retryDelay time.Duration, retryLimit int) error
}

func Wait(url url.URL, retryDelay time.Duration, retryLimit int) error {
	// log := logger.WithField("function", "Wait")
	// select the appropriate waiter
	var waiter Waiter
	switch url.Scheme {
	case "postgres":
		waiter = NewPostgresWaiter()
	case "tcp":
		waiter = NewTCPWaiter()
	default:
		return fmt.Errorf("unknown scheme: %s", url.Scheme)
	}
	// wait For IT!!
	return waiter.Wait(url, retryDelay, retryLimit)
}
