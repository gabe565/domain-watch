package util

import "errors"

var ErrNotConfigured = errors.New("missing configuration")

var ErrUnexpectedStatus = errors.New("unexpected status code")
