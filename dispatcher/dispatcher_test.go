package dispatcher_test

import (
	"sync"
	"testing"

	"github.com/242617/other/dispatcher"
)

func TestBasic(t *testing.T) {
	var l sync.Mutex
	var i int

	EventTick := "tick"

	h1 := func(dispatcher.Event) {
		l.Lock()
		defer l.Unlock()
		i++
	}
	h2 := func(dispatcher.Event) {
		l.Lock()
		defer l.Unlock()
		i--
	}
	d := dispatcher.NewEventDispatcher()

	d.AddEventListener(EventTick, &h1)
	d.AddEventListener(EventTick, &h2)
	defer func() {
		d.RemoveEventListener(EventTick, &h1)
		d.RemoveEventListener(EventTick, &h2)
	}()

	d.Dispatch(dispatcher.Event{Type: EventTick})
	d.Dispatch(dispatcher.Event{Type: EventTick})
	d.Dispatch(dispatcher.Event{Type: EventTick})

	if i != 0 {
		t.Fatal("unexpected result")
	}
}
