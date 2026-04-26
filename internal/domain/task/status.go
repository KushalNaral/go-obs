package task

type TaskStatus int

const (
	StatePending TaskStatus = iota
	StateRunning
	StateCompleted
	StateFailed
	StateRetried
	StateDeadLettered
)

var TransistionTable map[TaskStatus]map[TaskStatus]bool

func CanTransitionTo() {
}
