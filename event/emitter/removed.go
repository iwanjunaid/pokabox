package emitter

import (
	"time"

	"github.com/iwanjunaid/pokabox/model"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventRemoved(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID, record *model.OutboxRecord) {
	if e != nil {
		eventRemoved := event.Removed{
			PickerGroupID: pickerGroupID,
			OutboxRecord:  record,
			Timestamp:     timestamp,
		}

		e(eventRemoved)
	}
}
