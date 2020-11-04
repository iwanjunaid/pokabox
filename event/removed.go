package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type Removed struct {
	PickerGroupID uuid.UUID
	OutboxRecord  *model.OutboxRecord
	Timestamp     time.Time
}

func (r Removed) String() string {
	return fmt.Sprintf("[%s:%s] Message with ID %s successfully removed",
		PREFIX, r.PickerGroupID, r.OutboxRecord.ID)
}

func (r Removed) GetPickerGroupID() uuid.UUID {
	return r.PickerGroupID
}

func (r Removed) GetRecord() *model.OutboxRecord {
	return r.OutboxRecord
}

func (r Removed) GetTimestamp() time.Time {
	return r.Timestamp
}
