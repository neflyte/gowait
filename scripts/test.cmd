@ECHO OFF
SET GOWAIT_URL=postgres://postgres@localhost:5432/postgres?sslmode=disable
SET GOWAIT_RETRY_DELAY=3s
SET GOWAIT_RETRY_LIMIT=3
SET GOWAIT_SECRET=postgres
SET GOWAIT_LOG_FORMAT=text
ECHO Test Parameters:
ECHO     URL: %GOWAIT_URL%, Retry Delay: %GOWAIT_RETRY_DELAY%, Retry Limit: %GOWAIT_RETRY_LIMIT%, Secret: %GOWAIT_SECRET%, LogFormat: %GOWAIT_LOG_FORMAT%
ECHO.
bin\gowait
