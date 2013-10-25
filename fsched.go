package fsched

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

type Scheduler struct {
	heap *eventHeap
	now  time.Time
}

func MakeScheduler() Scheduler {
	s := Scheduler{heap: new(eventHeap)}
	*s.heap = make([]event, 0)
	return s
}

func MakeSchedulerTime(t time.Time) Scheduler {
	s := Scheduler{new(eventHeap), t}
	*s.heap = make([]event, 0)
	return s
}

func (s Scheduler) Now() time.Time {
	return s.now
}

func (s Scheduler) Empty() bool {
	return s.heap.Len() < 1
}

// Returns error if t < e.Now()
func (s Scheduler) Schedule(f func(time.Time) interface{}, time time.Time) error {
	if time.Before(s.now) {
		return ErrPast
	}
	heap.Push(s.heap, event{f, time})
	return nil
}

// Returns an error if offset < 0
func (s Scheduler) ScheduleOffset(f func(time.Time) interface{}, offset time.Duration) error {
	return s.Schedule(f, s.now.Add(offset))
}

// Returns the timestamp on the next scheduled event
func (s Scheduler) PeekNext() (time.Time, error) {
	var t time.Time
	if s.Empty() {
		return t, ErrEmpty
	}
	return (*s.heap)[0].time, nil
}

func (s Scheduler) CallNext() interface{} {
	evt := heap.Pop(s.heap).(event)
	return evt.f(evt.time)
}
