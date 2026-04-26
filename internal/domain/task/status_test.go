package task

import "testing"

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		status TaskStatus
		want   string
	}{
		{StatusPending, "pending"},
		{StatusRunning, "running"},
		{StatusCompleted, "completed"},
		{StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{"pending is valid", StatusPending, true},
		{"running is valid", StatusRunning, true},
		{"completed is valid", StatusCompleted, true},
		{"failed is valid", StatusFailed, true},
		{"zero value is invalid", TaskStatus(""), false},
		{"unknown status is invalid", TaskStatus("unknown"), false},
		{"uppercase is invalid", TaskStatus("PENDING"), false},
		{"typo is invalid", TaskStatus("pendng"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status TaskStatus
		want   bool
	}{
		{StatusPending, false},
		{StatusRunning, false},
		{StatusCompleted, true},
		{StatusFailed, true},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.IsTerminal(); got != tt.want {
				t.Errorf("IsTerminal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from TaskStatus
		to   TaskStatus
		want bool
	}{
		// from pending
		{"pending -> pending", StatusPending, StatusPending, false},
		{"pending -> running", StatusPending, StatusRunning, true},
		{"pending -> completed", StatusPending, StatusCompleted, false},
		{"pending -> failed", StatusPending, StatusFailed, false},

		// from running
		{"running -> pending (retry)", StatusRunning, StatusPending, true},
		{"running -> running", StatusRunning, StatusRunning, false},
		{"running -> completed", StatusRunning, StatusCompleted, true},
		{"running -> failed", StatusRunning, StatusFailed, true},

		// from completed (terminal)
		{"completed -> pending", StatusCompleted, StatusPending, false},
		{"completed -> running", StatusCompleted, StatusRunning, false},
		{"completed -> completed", StatusCompleted, StatusCompleted, false},
		{"completed -> failed", StatusCompleted, StatusFailed, false},

		// from failed (terminal)
		{"failed -> pending", StatusFailed, StatusPending, false},
		{"failed -> running", StatusFailed, StatusRunning, false},
		{"failed -> completed", StatusFailed, StatusCompleted, false},
		{"failed -> failed", StatusFailed, StatusFailed, false},

		// unknown / zero values fail safe
		{"zero -> running", TaskStatus(""), StatusRunning, false},
		{"pending -> unknown", StatusPending, TaskStatus("garbage"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("%s.CanTransitionTo(%s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
