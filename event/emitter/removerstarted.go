package emitter

import (
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventRemoverStarted(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID) {
	if e != nil {
		eventRemoverStarted := event.RemoverStarted{
			PickerGroupID: pickerGroupID,
			Timestamp:     timestamp,
		}

		e(eventRemoverStarted)
	}
}
