package utils

import (
	"github.com/neflyte/gowait/internal/logger"
	"net/url"
)

func SanitizedURLString(urlWithCreds url.URL) string {
	log := logger.WithField("function", "SanitizedURLString")
	clone, err := url.Parse(urlWithCreds.String())
	if err != nil {
		log.Errorf("unable to clone url: %s", err)
		return urlWithCreds.String()
	}
	clone.User = url.User(clone.User.Username())
	return clone.String()
}
