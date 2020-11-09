package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PickerStarted struct {
	PickerGroupID uuid.UUID
	Timestamp     time.Time
}

func (p PickerStarted) String() string {
	return fmt.Sprintf("[%s:%s] Picker started", PREFIX, p.PickerGroupID)
}

func (p PickerStarted) GetPickerGroupID() uuid.UUID {
	return p.PickerGroupID
}

func (p PickerStarted) GetTimestamp() time.Time {
	return p.Timestamp
}
