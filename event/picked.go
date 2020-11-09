package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type Picked struct {
	OutboxRecord *model.OutboxRecord
	Timestamp    time.Time
}

func (f Picked) String() string {
	id := f.OutboxRecord.ID
	groupID := f.OutboxRecord.GroupID

	return fmt.Sprintf("[%s:%s] Message with ID %s picked", PREFIX, groupID, id)
}

func (p Picked) GetPickerGroupID() uuid.UUID {
	return p.OutboxRecord.GroupID
}

func (p Picked) GetRecord() *model.OutboxRecord {
	return p.OutboxRecord
}

func (p Picked) GetTimestamp() time.Time {
	return p.Timestamp
}
