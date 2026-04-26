package task

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewTaskID_NotZero(t *testing.T) {
	id := NewTaskID()
	if id.IsZero() {
		t.Errorf("NewTaskID() returned zero value: %q", id)
	}
}

func TestNewTaskID_Unique(t *testing.T) {
	id1 := NewTaskID()
	id2 := NewTaskID()
	if id1 == id2 {
		t.Errorf("NewTaskID() returned duplicate ids: %q == %q", id1, id2)
	}
}

func TestParseTaskID(t *testing.T) {
	validLower := uuid.New().String()
	validUpper := strings.ToUpper(uuid.New().String())

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid lowercase uuid", input: validLower, wantErr: false},
		{name: "valid uppercase uuid", input: validUpper, wantErr: false},
		{name: "empty string", input: "", wantErr: true},
		{name: "garbage string", input: "abc123", wantErr: true},
		{name: "malformed uuid", input: uuid.New().String() + "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTaskID(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseTaskID(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseTaskID(%q) unexpected error: %v", tt.input, err)
			}
			if got.IsZero() {
				t.Errorf("ParseTaskID(%q) returned zero TaskID", tt.input)
			}
		})
	}
}

func TestTaskID_RoundTrip(t *testing.T) {
	original := NewTaskID()
	parsed, err := ParseTaskID(original.String())
	if err != nil {
		t.Fatalf("round-trip parse failed: %v", err)
	}
	if parsed != original {
		t.Errorf("round-trip mismatch: got %q, want %q", parsed, original)
	}
}

func TestTaskID_Equality(t *testing.T) {
	raw := uuid.New().String()

	a, err := ParseTaskID(raw)
	if err != nil {
		t.Fatalf("ParseTaskID(%q) failed: %v", raw, err)
	}
	b, err := ParseTaskID(raw)
	if err != nil {
		t.Fatalf("ParseTaskID(%q) failed: %v", raw, err)
	}

	if a != b {
		t.Errorf("expected equal TaskIDs to compare equal: %q != %q", a, b)
	}
}

func TestTaskID_IsZero(t *testing.T) {
	tests := []struct {
		name string
		id   TaskID
		want bool
	}{
		{name: "zero value", id: TaskID{}, want: true},
		{name: "freshly created", id: NewTaskID(), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskID_UUID(t *testing.T) {
	raw := uuid.New()
	id, err := ParseTaskID(raw.String())
	if err != nil {
		t.Fatalf("ParseTaskID(%q) failed: %v", raw, err)
	}
	if id.UUID() != raw {
		t.Errorf("UUID() = %v, want %v", id.UUID(), raw)
	}
}
