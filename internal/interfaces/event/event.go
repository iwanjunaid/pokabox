package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event interface {
	fmt.Stringer
	GetPickerGroupID() uuid.UUID
	GetTimestamp() time.Time
}
