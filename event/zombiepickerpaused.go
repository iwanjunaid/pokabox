package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ZombiePickerPaused struct {
	PickerGroupID uuid.UUID
	Timestamp     time.Time
}

func (z ZombiePickerPaused) String() string {
	return fmt.Sprintf("[%s:%s] Zombie picker paused", PREFIX, z.PickerGroupID)
}

func (z ZombiePickerPaused) GetPickerGroupID() uuid.UUID {
	return z.PickerGroupID
}

func (z ZombiePickerPaused) GetTimestamp() time.Time {
	return z.Timestamp
}
