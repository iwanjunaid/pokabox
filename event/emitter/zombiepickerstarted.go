package emitter

import (
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventZombiePickerStarted(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID) {
	if e != nil {
		eventZombiePickerStarted := event.ZombiePickerStarted{
			PickerGroupID: pickerGroupID,
			Timestamp:     timestamp,
		}

		e(eventZombiePickerStarted)
	}
}
