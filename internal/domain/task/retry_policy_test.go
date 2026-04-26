package task

import (
	"errors"
	"testing"
	"time"
)

func TestNewRetryPolicy_Valid(t *testing.T) {
	in := RetryPolicyInput{
		MaxAttempts: 5,
		BaseDelay:   time.Second,
		MaxDelay:    30 * time.Second,
	}

	p, err := NewRetryPolicy(in)
	if err != nil {
		t.Fatalf("NewRetryPolicy(%+v) unexpected error: %v", in, err)
	}
	if p.MaxAttempts() != in.MaxAttempts {
		t.Errorf("MaxAttempts = %d, want %d", p.MaxAttempts(), in.MaxAttempts)
	}
	if p.BaseDelay() != in.BaseDelay {
		t.Errorf("BaseDelay = %v, want %v", p.BaseDelay(), in.BaseDelay)
	}
	if p.MaxDelay() != in.MaxDelay {
		t.Errorf("MaxDelay = %v, want %v", p.MaxDelay(), in.MaxDelay)
	}
}

func TestNewRetryPolicy_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		input   RetryPolicyInput
		wantErr error
	}{
		{
			name:    "zero max attempts",
			input:   RetryPolicyInput{MaxAttempts: 0, BaseDelay: time.Second, MaxDelay: time.Second},
			wantErr: ErrInvalidMaxAttempts,
		},
		{
			name:    "negative max attempts",
			input:   RetryPolicyInput{MaxAttempts: -1, BaseDelay: time.Second, MaxDelay: time.Second},
			wantErr: ErrInvalidMaxAttempts,
		},
		{
			name:    "zero base delay",
			input:   RetryPolicyInput{MaxAttempts: 3, BaseDelay: 0, MaxDelay: time.Second},
			wantErr: ErrInvalidBaseDelay,
		},
		{
			name:    "negative base delay",
			input:   RetryPolicyInput{MaxAttempts: 3, BaseDelay: -time.Second, MaxDelay: time.Second},
			wantErr: ErrInvalidBaseDelay,
		},
		{
			name:    "max delay below base delay",
			input:   RetryPolicyInput{MaxAttempts: 3, BaseDelay: 10 * time.Second, MaxDelay: time.Second},
			wantErr: ErrInvalidMaxDelay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRetryPolicy(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewRetryPolicy() err = %v, want %v", err, tt.wantErr)
			}
			if got != (RetryPolicy{}) {
				t.Errorf("expected zero RetryPolicy on error, got %+v", got)
			}
		})
	}
}

func TestDefaultRetryPolicy(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxAttempts() < 1 {
		t.Errorf("default MaxAttempts must be >= 1, got %d", p.MaxAttempts())
	}
	if p.BaseDelay() <= 0 {
		t.Errorf("default BaseDelay must be > 0, got %v", p.BaseDelay())
	}
	if p.MaxDelay() < p.BaseDelay() {
		t.Errorf("default MaxDelay (%v) must be >= BaseDelay (%v)", p.MaxDelay(), p.BaseDelay())
	}
}

func TestRetryPolicy_ShouldRetry(t *testing.T) {
	p, err := NewRetryPolicy(RetryPolicyInput{
		MaxAttempts: 3,
		BaseDelay:   time.Second,
		MaxDelay:    time.Second,
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name         string
		attemptsMade int
		want         bool
	}{
		{"no attempts yet", 0, true},
		{"one attempt made", 1, true},
		{"one less than max", 2, true},
		{"max reached", 3, false},
		{"beyond max", 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.ShouldRetry(tt.attemptsMade); got != tt.want {
				t.Errorf("ShouldRetry(%d) = %v, want %v", tt.attemptsMade, got, tt.want)
			}
		})
	}
}

func TestRetryPolicy_NextDelay_Exponential(t *testing.T) {
	p, err := NewRetryPolicy(RetryPolicyInput{
		MaxAttempts: 10,
		BaseDelay:   time.Second,
		MaxDelay:    time.Hour, // high cap so we observe pure exponential growth
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		attemptsMade int
		want         time.Duration
	}{
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
		{5, 16 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.want.String(), func(t *testing.T) {
			if got := p.NextDelay(tt.attemptsMade); got != tt.want {
				t.Errorf("NextDelay(%d) = %v, want %v", tt.attemptsMade, got, tt.want)
			}
		})
	}
}

func TestRetryPolicy_NextDelay_CappedAtMax(t *testing.T) {
	p, err := NewRetryPolicy(RetryPolicyInput{
		MaxAttempts: 100,
		BaseDelay:   time.Second,
		MaxDelay:    10 * time.Second,
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []int{5, 6, 10, 50, 100, 1000}
	for _, attempts := range tests {
		if got := p.NextDelay(attempts); got != p.MaxDelay() {
			t.Errorf("NextDelay(%d) = %v, want %v (max)", attempts, got, p.MaxDelay())
		}
	}
}

func TestRetryPolicy_NextDelay_ZeroOrNegative(t *testing.T) {
	p, err := NewRetryPolicy(RetryPolicyInput{
		MaxAttempts: 3,
		BaseDelay:   2 * time.Second,
		MaxDelay:    30 * time.Second,
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	for _, attempts := range []int{0, -1, -100} {
		if got := p.NextDelay(attempts); got != p.BaseDelay() {
			t.Errorf("NextDelay(%d) = %v, want %v (base)", attempts, got, p.BaseDelay())
		}
	}
}
