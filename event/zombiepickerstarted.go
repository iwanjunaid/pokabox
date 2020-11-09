package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ZombiePickerStarted struct {
	PickerGroupID uuid.UUID
	Timestamp     time.Time
}

func (z ZombiePickerStarted) String() string {
	return fmt.Sprintf("[%s:%s] Zombie picker started", PREFIX, z.PickerGroupID)
}

func (z ZombiePickerStarted) GetPickerGroupID() uuid.UUID {
	return z.PickerGroupID
}

func (z ZombiePickerStarted) GetTimestamp() time.Time {
	return z.Timestamp
}
