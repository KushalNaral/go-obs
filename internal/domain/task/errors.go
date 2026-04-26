package task

import "errors"

var (
	ErrInvalidMaxAttempts = errors.New("max attemps must be >= 1")
	ErrInvalidBaseDelay   = errors.New("base delay must be > 0")
	ErrInvalidMaxDelay    = errors.New("max delay must be >= base delay")
)
