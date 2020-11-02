package event

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type Fetched struct {
	OutboxRecord *model.OutboxRecord
}

func (f Fetched) String() string {
	id := f.OutboxRecord.ID
	groupID := f.OutboxRecord.GroupID

	return fmt.Sprintf("[%s:%s] Message with ID %s fetched", PREFIX, groupID, id)
}

func (f Fetched) GetGroupID() uuid.UUID {
	return f.OutboxRecord.GroupID
}

func (f Fetched) GetRecord() *model.OutboxRecord {
	return f.OutboxRecord
}
