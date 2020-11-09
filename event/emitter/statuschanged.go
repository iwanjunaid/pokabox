package emitter

import (
	"time"

	"github.com/iwanjunaid/pokabox/event"
	"github.com/iwanjunaid/pokabox/model"
)

func EmitEventStatusChanged(e event.EventHandler, timestamp time.Time,
	fromStatus string, toStatus string, record *model.OutboxRecord) {
	if e != nil {
		eventStatusChanged := event.StatusChanged{
			From:         fromStatus,
			To:           toStatus,
			OutboxRecord: record,
			Timestamp:    timestamp,
		}

		e(eventStatusChanged)
	}
}
