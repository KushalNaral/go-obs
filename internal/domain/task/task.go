package task

import (
	"time"
)

// Rule of thumb: in a DDD aggregate,
// Almost every field should be unexported (lowercase).
// The outside world reads via methods (t.Status(), t.ID()) and mutates only via behavior methods t.MarkRunning()
type Task struct {
	ID TaskID
	// send_email, process_image, generate_report...
	Type        string
	Priority    Priority
	ScheduledAt time.Time
	Payload     []byte
	createdAt   time.Time
	attempts    []Attempt
	// Pending ->  Running -> Completed / Failed / Retried / Deadlettered
	status    TaskStatus
	events    []DomainEvent
	lastError string
	retry     Retry
}

type Priority int

const (
	PriorityHigh Priority = iota
	PriorityMedium
	PriorityLow
)

type Retry struct {
	Policy      Policy
	ScheduledAt time.Time
	Duration    time.Duration
	History     []map[int]History
}

type Policy struct {
	maxAttempts      int
	baseDelay        int
	backoffStrategry BackoffStrategy
}

type (
	Attempt         struct{}
	History         struct{}
	DomainEvent     struct{}
	BackoffStrategy struct{}
)
