package eventheap

import (
	"container/heap"
	"errors"
	"time"
)

type event struct {
	f    func(time.Time) interface{}
	time time.Time
}

type eventHeap []event

func (e eventHeap) Len() int            { return len(e) }
func (e eventHeap) Less(i, j int) bool  { return e[i].time.Before(e[j].time) }
func (e eventHeap) Swap(i, j int)       { e[i], e[j] = e[j], e[i] }
func (e *eventHeap) Push(x interface{}) { *e = append(*e, x.(event)) }
func (e *eventHeap) Pop() interface{} {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}

var (
	ErrPast  = errors.New("Event scheduled in the past")
	ErrEmpty = errors.New("Empty")
)

type EventHeap struct {
	heap *eventHeap
	now  time.Time
}

func MakeEventHeap() EventHeap {
	eh := EventHeap{heap: new(eventHeap)}
	*eh.heap = make([]event, 0)
	return eh
}

func MakeEventHeapTime(t time.Time) EventHeap {
	eh := EventHeap{new(eventHeap), t}
	*eh.heap = make([]event, 0)
	return eh
}

func (e EventHeap) Now() time.Time {
	return e.now
}

func (e EventHeap) Empty() bool {
	return e.heap.Len() < 1
}

// Returns error if t < e.Now()
func (e EventHeap) Schedule(f func(time.Time) interface{}, time time.Time) error {
	if time.Before(e.now) {
		return ErrPast
	}
	heap.Push(e.heap, event{f, time})
	return nil
}

// Returns an error if offset < 0
func (e EventHeap) ScheduleOffset(f func(time.Time) interface{}, offset time.Duration) error {
	return e.Schedule(f, e.now.Add(offset))
}

// Returns the timestamp on the next scheduled event
func (e EventHeap) PeekNext() (time.Time, error) {
	var t time.Time
	if e.Empty() {
		return t, ErrEmpty
	}
	return (*e.heap)[0].time, nil
}

func (e EventHeap) CallNext() interface{} {
	evt := heap.Pop(e.heap).(event)
	return evt.f(evt.time)
}
