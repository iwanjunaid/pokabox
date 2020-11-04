package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PickerPaused struct {
	PickerGroupID uuid.UUID
	Timestamp     time.Time
}

func (p PickerPaused) String() string {
	return fmt.Sprintf("[%s:%s] Picker paused", PREFIX, p.PickerGroupID)
}

func (p PickerPaused) GetPickerGroupID() uuid.UUID {
	return p.PickerGroupID
}

func (p PickerPaused) GetTimestamp() time.Time {
	return p.Timestamp
}
