package event

import "fmt"

type Event interface {
	fmt.Stringer
}
