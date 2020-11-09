package emitter

import (
	"time"

	"github.com/iwanjunaid/pokabox/event"
	"github.com/iwanjunaid/pokabox/model"
)

func EmitEventPicked(e event.EventHandler, timestamp time.Time, record *model.OutboxRecord) {
	if e != nil {
		eventPicked := event.Picked{
			OutboxRecord: record,
			Timestamp:    timestamp,
		}

		e(eventPicked)
	}
}
