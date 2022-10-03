## domain-watch



```
domain-watch [flags] domain...
```

### Options

```
      --completion string       Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.
  -e, --every duration          enable cron mode and configure update interval
  -h, --help                    help for domain-watch
      --log-format string       log formatter (text, json) (default "text")
  -l, --log-level string        log level (trace, debug, info, warning, error, fatal, panic) (default "info")
  -s, --sleep duration          sleep time between queries to avoid rate limits (default 3s)
      --telegram-chat int       Telegram chat ID
      --telegram-token string   Telegram token
  -t, --threshold ints          configure expiration notifications (default [1,7])
```

