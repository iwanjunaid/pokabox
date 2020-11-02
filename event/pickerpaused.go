package event

import (
	"fmt"

	"github.com/google/uuid"
)

type PickerPaused struct {
	GroupID uuid.UUID
}

func (p PickerPaused) String() string {
	return fmt.Sprintf("[%s:%s] Picker paused", PREFIX, p.GroupID)
}

func (p PickerPaused) GetGroupID() uuid.UUID {
	return p.GroupID
}
