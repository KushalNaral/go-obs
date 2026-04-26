package task

import (
	"fmt"

	"github.com/google/uuid"
)

// This is a Value Object, wraps primitive data, adds meaning, owns its own rules, prevents invalid states from entering the domain
type TaskID struct {
	id uuid.UUID
}

func NewTaskID() TaskID {
	return TaskID{
		id: uuid.New(),
	}
}

func (t TaskID) String() string {
	return t.id.String()
}

func (t TaskID) IsZero() bool {
	return t.id == uuid.Nil
}

func (t TaskID) UUID() uuid.UUID {
	return t.id
}

func ParseTaskID(id string) (TaskID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return TaskID{}, fmt.Errorf("parse task id %q: %w", id, err)
	}
	return TaskID{
		id: parsedID,
	}, nil
}
