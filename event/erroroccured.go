package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ErrorOccured struct {
	PickerGroupID uuid.UUID
	Error         error
	Timestamp     time.Time
}

func (e ErrorOccured) String() string {
	return fmt.Sprintf("[%s:%s] Error occured: %s", PREFIX, e.PickerGroupID, e.Error.Error())
}

func (e ErrorOccured) GetPickerGroupID() uuid.UUID {
	return e.PickerGroupID
}

func (e ErrorOccured) GetError() error {
	return e.Error
}

func (e ErrorOccured) GetTimestamp() time.Time {
	return e.Timestamp
}
