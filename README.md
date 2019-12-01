# gowait: a golang service waiter

## Purpose
I started using Kubernetes often in my daily life and found a persistent irritation in waiting for scores of deployed 
pods to come up for local development. There are shell script-based solutions out there, and they work. However, shells 
in a production container are viewed as a security risk. To mitigate this, I wanted a statically-linked, single-file 
binary that could run in a very small container and give me the flexibility to configure it however I may need.

`gowait` is what I came up with. :)

## Usage
Add the docker image as an InitContainer to your favourite deployment file and configure appropriately.
I publish a pre-built image for each release here: https://hub.docker.com/repository/docker/neflyte/gowait

### CLI Arguments

#### Usage
`gowait -c|-configsource <env|yaml|json> [-f|-configfile <fileLocation>]`

- -c | -configsource
    - Where to read app configuration from
    - Values:
        - `env`: environment variables (the default)
        - `yaml`: YAML file
        - `json`: JSON file
- -f | -configfile
    - Location of configuration file to read
    - Used when -configsource is set to `yaml` or `json`

### Environment Variables

- `GOWAIT_RETRY_DELAY`
    - The amount of time to delay before retrying
    - Expressed as a string suitable for passing to time.ParseDuration()
    - e.g.: `GOWAIT_RETRY_DELAY="15s"`
- `GOWAIT_RETRY_LIMIT`
    - The maximum number of attempts to make
    - Expressed as a positive integer value greater than zero
    - e.g.: `GOWAIT_RETRY_LIMIT="20"`
- `GOWAIT_URL`
    - A URL representing the service to wait for
    - e.g.: `GOWAIT_URL="postgres://user@localhost:5432/database?ssl_mode=disable"`
    - Supported URL schemes:
        - `postgres`
            - Uses lib/pq to attempt a connection to a PostgreSQL database
        - `tcp`
            - Attempts a connection to a TCP port
            - If an established connection is alive for at least one second, the attempt succeeded
 - `GOWAIT_SECRET_SOURCE`
    - Where to read the secret value from
    - e.g.: `GOWAIT_SECRET_SOURCE="file"`
    - Supported values:
        - `env`: environment variables (the default)
        - `file`: a file
 - `GOWAIT_SECRET_FILENAME`
    - The name of the file to read the secret value from
    - Used when `GOWAIT_SECRET_SOURCE="file"`
    - e.g.: `GOWAIT_SECRET_FILENAME="/tmp/secret.txt"`
 - `GOWAIT_SECRET`
    - The secret value
    - Used by default and when `GOWAIT_SECRET_SOURCE="env"`
    - e.g.: `GOWAIT_SECRET="fnord"`
 - `GOWAIT_LOG_FORMAT`
    - The format of log messages
    - e.g.: `GOWAIT_LOG_FORMAT="text"`
    - Supported values:
        - `text`: Human-readable text (the default)
        - `json`: logstash-like JSON

### YAML Configuration Example

```yaml
---
url: "postgres://user@localhost:5432/database?ssl_mode=disable"
retryDelay: "15s"
retryLimit: 20
secretSource: "file"
secretFilename: "/tmp/secret.txt"
logFormat: "text"
```

### JSON Configuration Example

```json
{
  "url": "postgres://user@localhost:5432/database?ssl_mode=disable",
  "retryDelay": "15s",
  "retryLimit": 20,
  "secretSource": "file",
  "secretFilename": "/tmp/secret.txt",
  "logFormat": "text"
}
```
