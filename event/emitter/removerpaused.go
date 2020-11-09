package emitter

import (
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventRemoverPaused(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID) {
	if e != nil {
		eventRemoverPaused := event.RemoverPaused{
			PickerGroupID: pickerGroupID,
			Timestamp:     timestamp,
		}

		e(eventRemoverPaused)
	}
}
