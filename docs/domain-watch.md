## domain-watch



```
domain-watch [flags] domain...
```

### Options

```
      --completion string        Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.
      --domains strings          List of domains to watch
  -e, --every duration           enable cron mode and configure update interval
      --gotify-token string      Gotify app token
      --gotify-url string        Gotify URL (include https:// and port if non-standard)
  -h, --help                     help for domain-watch
      --log-format string        log formatter (one of auto, color, plain, json) (default "auto")
  -l, --log-level string         log level (one of debug, info, warn, error) (default "info")
      --metrics-address string   Prometheus metrics API listen address (default ":9090")
      --metrics-enabled          Enables Prometheus metrics API
  -s, --sleep duration           sleep time between queries to avoid rate limits (default 3s)
      --telegram-chat int        Telegram chat ID
      --telegram-token string    Telegram token
  -t, --threshold ints           configure expiration notifications (default [1,7])
  -v, --version                  version for domain-watch
```

