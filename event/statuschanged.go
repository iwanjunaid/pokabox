package event

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type StatusChanged struct {
	From         string
	To           string
	OutboxRecord *model.OutboxRecord
}

func (s StatusChanged) String() string {
	id := s.OutboxRecord.ID
	groupID := s.OutboxRecord.GroupID

	return fmt.Sprintf("[%s:%s] Status changed from %s to %s for message with ID %s",
		PREFIX, groupID, s.From, s.To, id)
}

func (s StatusChanged) GetGroupID() uuid.UUID {
	return s.OutboxRecord.GroupID
}

func (s StatusChanged) GetFrom() string {
	return s.From
}

func (s StatusChanged) GetTo() string {
	return s.To
}

func (s StatusChanged) GetRecord() *model.OutboxRecord {
	return s.OutboxRecord
}
