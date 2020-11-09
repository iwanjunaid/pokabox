package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RemoverPaused struct {
	PickerGroupID uuid.UUID
	Timestamp     time.Time
}

func (r RemoverPaused) String() string {
	return fmt.Sprintf("[%s:%s] Remover paused", PREFIX, r.PickerGroupID)
}

func (r RemoverPaused) GetPickerGroupID() uuid.UUID {
	return r.PickerGroupID
}

func (r RemoverPaused) GetTimestamp() time.Time {
	return r.Timestamp
}
