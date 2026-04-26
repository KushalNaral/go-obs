package task

import (
	"time"
)

type RetryPolicy struct {
	maxAttempts int
	baseDelay   time.Duration
	maxDelay    time.Duration
}

type RetryPolicyInput struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func NewRetryPolicy(input RetryPolicyInput) (RetryPolicy, error) {
	if input.MaxAttempts < 1 {
		return RetryPolicy{}, ErrInvalidMaxAttempts
	}

	if input.BaseDelay <= 0 {
		return RetryPolicy{}, ErrInvalidBaseDelay
	}

	if input.MaxDelay < input.BaseDelay {
		return RetryPolicy{}, ErrInvalidMaxDelay
	}

	return RetryPolicy{
		maxAttempts: input.MaxAttempts,
		baseDelay:   input.BaseDelay,
		maxDelay:    input.MaxDelay,
	}, nil
}

func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		maxAttempts: 3,
		baseDelay:   time.Second * 1,
		maxDelay:    time.Second * 30,
	}
}

func (p RetryPolicy) ShouldRetry(attemptsMade int) bool {
	return attemptsMade < p.maxAttempts
}

// The Backoff Math
//
//	For exponential backoff capped at max:
//
//	delay = baseDelay * 2^(attemptsMade - 1)
//	if delay > maxDelay → delay = maxDelay
func (p RetryPolicy) NextDelay(attemptsMade int) time.Duration {
	if attemptsMade <= 0 {
		return p.baseDelay
	}

	if attemptsMade > 30 {
		return p.maxDelay
	}

	delay := p.baseDelay << (attemptsMade - 1)
	if delay > p.maxDelay || delay < 0 {
		return p.maxDelay
	}

	return delay
}

func (p RetryPolicy) MaxAttempts() int         { return p.maxAttempts }
func (p RetryPolicy) MaxDelay() time.Duration  { return p.maxDelay }
func (p RetryPolicy) BaseDelay() time.Duration { return p.baseDelay }
