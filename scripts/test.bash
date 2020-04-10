#!/usr/bin/env bash
GOWAIT_URL="postgres://postgres@localhost:5432/postgres?sslmode=disable"
GOWAIT_RETRY_DELAY="3s"
GOWAIT_RETRY_LIMIT="3"
GOWAIT_SECRET="postgres"
GOWAIT_LOG_FORMAT="text"
echo "Test Parameters:"
echo "    URL: $GOWAIT_URL, Retry Delay: $GOWAIT_RETRY_DELAY, Retry Limit: $GOWAIT_RETRY_LIMIT, Secret: $GOWAIT_SECRET, LogFormat: $GOWAIT_LOG_FORMAT"
bin/gowait
