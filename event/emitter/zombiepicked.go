package emitter

import (
	"time"

	"github.com/iwanjunaid/pokabox/model"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/event"
)

func EmitEventZombiePicked(e event.EventHandler, timestamp time.Time,
	pickerGroupID uuid.UUID, record *model.OutboxRecord) {
	if e != nil {
		eventZombiePicked := event.ZombiePicked{
			PickerGroupID: pickerGroupID,
			OutboxRecord:  record,
			Timestamp:     timestamp,
		}

		e(eventZombiePicked)
	}
}
