package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iwanjunaid/pokabox/model"
)

type ZombiePicked struct {
	PickerGroupID uuid.UUID
	OutboxRecord  *model.OutboxRecord
	Timestamp     time.Time
}

func (z ZombiePicked) String() string {
	pickerGroupID := z.PickerGroupID
	id := z.OutboxRecord.ID
	originGroupID := z.OutboxRecord.GroupID

	return fmt.Sprintf("[%s:%s] Zombie message with ID %s from Origin Group ID %s is picked",
		PREFIX, pickerGroupID, id, originGroupID)
}

func (z ZombiePicked) GetPickerGroupID() uuid.UUID {
	return z.PickerGroupID
}

func (z ZombiePicked) GetRecord() *model.OutboxRecord {
	return z.OutboxRecord
}

func (z ZombiePicked) GetTimestamp() time.Time {
	return z.Timestamp
}
