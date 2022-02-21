package utils

import (
	"net/url"

	"github.com/neflyte/gowait/internal/logger"
)

// SanitizedURLString returns a parsed URL string with user credentials removed
func SanitizedURLString(urlWithCreds url.URL) string {
	log := logger.Function("SanitizedURLString")
	clone, err := url.Parse(urlWithCreds.String())
	if err != nil {
		// The URL itself should not be logged unless credentials are removed
		log.Err(err).
			Error("unable to clone url")
		return urlWithCreds.String()
	}
	if clone.User != nil {
		clone.User = url.User(clone.User.Username())
	}
	return clone.String()
}
