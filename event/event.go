package event

import (
	"github.com/iwanjunaid/pokabox/internal/interfaces/event"
)

const PREFIX = "pokabox"

type EventHandler func(event.Event)
