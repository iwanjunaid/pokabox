package event

import (
	"fmt"

	"github.com/google/uuid"
)

type PickerStarted struct {
	GroupID uuid.UUID
}

func (p PickerStarted) String() string {
	return fmt.Sprintf("[%s:%s] Picker started", PREFIX, p.GroupID)
}

func (p PickerStarted) GetGroupID() uuid.UUID {
	return p.GroupID
}
