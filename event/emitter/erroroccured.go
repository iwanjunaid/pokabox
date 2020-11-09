package emitter

import (
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventErrorOccured(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID, err error) {
	if e != nil {
		eventErrorOccured := event.ErrorOccured{
			PickerGroupID: pickerGroupID,
			Error:         err,
			Timestamp:     timestamp,
		}

		e(eventErrorOccured)
	}
}
