package main

import (
	whoisparser "github.com/likexian/whois-parser-go"
	"time"
)

type Config struct {
	WhoisCache map[string]whoisparser.WhoisInfo
	RunEvery   string
	Sleep      time.Duration
	Token      string
	ChatId     int64
}
