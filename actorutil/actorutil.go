package actorutil

import (
	"time"

	"github.com/anthdm/hollywood/actor"
)

// SendAfter sends msg to pid once d has passed. Timers run on their own goroutine,
// so an actor uses this to hear back from one via a message instead of having it
// touch the actor's state directly
func SendAfter(e *actor.Engine, pid *actor.PID, d time.Duration, msg any) *time.Timer {
	return time.AfterFunc(d, func() { e.Send(pid, msg) })
}
