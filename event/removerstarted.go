package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RemoverStarted struct {
	PickerGroupID uuid.UUID
	Timestamp     time.Time
}

func (r RemoverStarted) String() string {
	return fmt.Sprintf("[%s:%s] Remover started", PREFIX, r.PickerGroupID)
}

func (r RemoverStarted) GetPickerGroupID() uuid.UUID {
	return r.PickerGroupID
}

func (r RemoverStarted) GetTimestamp() time.Time {
	return r.Timestamp
}
