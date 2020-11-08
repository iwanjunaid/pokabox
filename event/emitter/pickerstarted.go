package emitter

import (
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventPickerStarted(e event.EventHandler, timestamp time.Time, pickerGroupID uuid.UUID) {
	if e != nil {
		pickerStarted := event.PickerStarted{
			PickerGroupID: pickerGroupID,
			Timestamp:     timestamp,
		}

		e(pickerStarted)
	}
}
