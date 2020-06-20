#!/usr/bin/env bash
# default to postgres test
export GOWAIT_URL="postgres://postgres@localhost:5432/postgres?sslmode=disable"
export GOWAIT_RETRY_DELAY="3s"
export GOWAIT_RETRY_LIMIT="3"
export GOWAIT_SECRET="postgres"
export GOWAIT_LOG_FORMAT="text"
# default program args
declare PROGARGS=""
# if a different test was specified, set it up
[ -n "${1}" ] && {
  TESTOPT="${1}"
  case $TESTOPT in
    "postgres")
      export GOWAIT_URL="postgres://postgres@localhost:5432/postgres?sslmode=disable"
      export GOWAIT_RETRY_DELAY="3s"
      export GOWAIT_RETRY_LIMIT="3"
      export GOWAIT_SECRET="postgres"
      export GOWAIT_LOG_FORMAT="text"
      ;;
    "postgres-yaml")
      export -n GOWAIT_URL GOWAIT_RETRY_DELAY GOWAIT_RETRY_LIMIT GOWAIT_SECRET GOWAIT_LOG_FORMAT
      PROGARGS="-c yaml -f testdata/postgres/config.yaml"
      ;;
    "postgres-json")
      export -n GOWAIT_URL GOWAIT_RETRY_DELAY GOWAIT_RETRY_LIMIT GOWAIT_SECRET GOWAIT_LOG_FORMAT
      PROGARGS="-c json -f testdata/postgres/config.json"
      ;;
    "http")
      export GOWAIT_URL="http://localhost:8080/"
      export GOWAIT_RETRY_DELAY="3s"
      export GOWAIT_RETRY_LIMIT="3"
      export GOWAIT_SECRET=""
      export GOWAIT_LOG_FORMAT="text"
      ;;
    "http-yaml")
      export -n GOWAIT_URL GOWAIT_RETRY_DELAY GOWAIT_RETRY_LIMIT GOWAIT_SECRET GOWAIT_LOG_FORMAT
      PROGARGS="-c yaml -f testdata/http/config.yaml"
      ;;
    "http-json")
      export -n GOWAIT_URL GOWAIT_RETRY_DELAY GOWAIT_RETRY_LIMIT GOWAIT_SECRET GOWAIT_LOG_FORMAT
      PROGARGS="-c json -f testdata/http/config.json"
      ;;
    "tcp")
      export GOWAIT_URL="tcp://localhost:8080/"
      export GOWAIT_RETRY_DELAY="3s"
      export GOWAIT_RETRY_LIMIT="3"
      export GOWAIT_SECRET=""
      export GOWAIT_LOG_FORMAT="text"
      ;;
    "kafka")
      export GOWAIT_URL="kafka://localhost:9092/"
      export GOWAIT_RETRY_DELAY="3s"
      export GOWAIT_RETRY_LIMIT="3"
      export GOWAIT_SECRET=""
      export GOWAIT_LOG_FORMAT="text"
      ;;
    *)
      echo "*  unknown test ${TESTOPT}; aborting"
      exit 1
      ;;
  esac
}
# display the parameters that were set
echo -e "Test Parameters:\n\tURL: $GOWAIT_URL\n\tRetry Delay: $GOWAIT_RETRY_DELAY\n\tRetry Limit: $GOWAIT_RETRY_LIMIT\n\tSecret: $GOWAIT_SECRET\n\tLogFormat: $GOWAIT_LOG_FORMAT"
# run the test
bin/gowait ${PROGARGS}
# unexport the test variables
export -n GOWAIT_URL GOWAIT_RETRY_DELAY GOWAIT_RETRY_LIMIT GOWAIT_SECRET GOWAIT_LOG_FORMAT
echo "done."