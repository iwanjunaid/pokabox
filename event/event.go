package event

import (
	"github.com/iwanjunaid/pokabox/internal/interfaces/event"
)

const PREFIX = "outbox"

type EventHandler func(event.Event)
