package task

// TaskStatus has these responsibilities
// 1. The set of legal statuses — the enum.
// 2. The transition rules — which status can become which.
// 3. A few helpful queries — is this terminal? what's the string form?
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

// This determines the status flow and what they can transition into and to
//
//	from           → allowed next
//
//	pending        → running
//	running        → pending      (retry — attempts remain)
//	running        → completed    (success)
//	running        → failed       (give up — attempts exhausted)
//	completed     → (none — terminal)
//	failed         → (none — terminal)
var validTransitions = map[TaskStatus]map[TaskStatus]bool{
	StatusPending: {
		StatusRunning: true,
	},
	StatusRunning: {
		StatusPending:   true,
		StatusCompleted: true,
		StatusFailed:    true,
	},
	StatusCompleted: {},
	StatusFailed:    {},
}

func (s TaskStatus) CanTransitionTo(next TaskStatus) bool {
	return validTransitions[s][next]
}

func (s TaskStatus) String() string {
	return string(s)
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

func (s TaskStatus) IsTerminal() bool {
	return s == StatusCompleted || s == StatusFailed
}
