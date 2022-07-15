package config

import (
	"time"
)

type Config struct {
	RunEvery time.Duration
	Sleep    time.Duration
	Token    string
	ChatId   int64
}
