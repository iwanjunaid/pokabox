package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type ZombieAcquired struct {
	PickerGroupID uuid.UUID
	OriginGroupID uuid.UUID
	OutboxRecord  *model.OutboxRecord
	Timestamp     time.Time
}

func (z ZombieAcquired) String() string {
	id := z.OutboxRecord.ID
	originGroupID := z.OriginGroupID
	groupID := z.OutboxRecord.GroupID

	return fmt.Sprintf("[%s:%s] Zombie message with ID %s from Origin Group ID %s is successfully acquired",
		PREFIX, groupID, id, originGroupID)
}

func (z ZombieAcquired) GetPickerGroupID() uuid.UUID {
	return z.PickerGroupID
}

func (z ZombieAcquired) GetOriginGroupID() uuid.UUID {
	return z.OriginGroupID
}

func (z ZombieAcquired) GetRecord() *model.OutboxRecord {
	return z.OutboxRecord
}

func (z ZombieAcquired) GetTimestamp() time.Time {
	return z.Timestamp
}
