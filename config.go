package main

import (
	"time"
)

type Config struct {
	RunEvery string
	Sleep    time.Duration
	Token    string
	ChatId   int64
}
