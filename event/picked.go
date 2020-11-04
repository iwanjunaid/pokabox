package event

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type Picked struct {
	OutboxRecord *model.OutboxRecord
}

func (f Picked) String() string {
	id := f.OutboxRecord.ID
	groupID := f.OutboxRecord.GroupID

	return fmt.Sprintf("[%s:%s] Message with ID %s picked", PREFIX, groupID, id)
}

func (f Picked) GetGroupID() uuid.UUID {
	return f.OutboxRecord.GroupID
}

func (f Picked) GetRecord() *model.OutboxRecord {
	return f.OutboxRecord
}
