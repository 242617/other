package dispatcher

import "sync"

type Event struct {
	Type  string
	Value interface{}
}

// type EventHandler func(event Event)
type EventDispatcher interface {
	Dispatch(Event)
	AddEventListener(eventType string, handler *func(event Event))
	RemoveEventListener(eventType string, handler *func(event Event))
}

func NewEventDispatcher() EventDispatcher {
	return &eventDispatcher{handlers: map[string][]*func(event Event){}}
}

type eventDispatcher struct {
	sync.Mutex
	handlers map[string][]*func(event Event)
}

func (ed *eventDispatcher) Dispatch(event Event) {
	ed.Lock()
	defer ed.Unlock()
	for _, handler := range ed.handlers[event.Type] {
		(*handler)(event)
	}
}

func (ed *eventDispatcher) AddEventListener(eventType string, handler *func(event Event)) {
	if _, ok := ed.handlers[eventType]; !ok {
		ed.handlers[eventType] = []*func(event Event){}
	}
	ed.Lock()
	defer ed.Unlock()
	ed.handlers[eventType] = append(ed.handlers[eventType], handler)
}

func (ed *eventDispatcher) RemoveEventListener(eventType string, handler *func(event Event)) {
	ed.Lock()
	defer ed.Unlock()
	i := indexOf(ed.handlers[eventType], handler)
	if i == -1 {
		return
	}
	ed.handlers[eventType] = remove(ed.handlers[eventType], i)
	if len(ed.handlers[eventType]) == 0 {
		delete(ed.handlers, eventType)
	}
}

func indexOf(handlers []*func(event Event), handler *func(event Event)) int {
	for i := 0; i < len(handlers); i++ {
		if handlers[i] == handler {
			return i
		}
	}
	return -1
}

func remove(handlers []*func(event Event), i int) []*func(event Event) {
	handlers[len(handlers)-1], handlers[i] = handlers[i], handlers[len(handlers)-1]
	return handlers[:len(handlers)-1]
}
