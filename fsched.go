// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note that this doc comment shares text with the Scheduler documentation.
// Please keep in sync.

// Package fsched implements a sequential scheduler for function callbacks.
// The scheduler keeps track of the current time (different from real time).
// Each call to CallNext fast-forwards to the timestamp on the earliest event,
// and executes the associated callback function. Note that the scheduler's
// time is arbitrary and must only be internally consistent; it is unrelated
// to any real sense of time (ie, clock cycles, seconds since epoch, etc).
//
// While the scheduler interface is generic, it is designed to be used for
// simulating highly-parallel processes where accurate time is important,
// but the simulation is too processor-intensive to be run in real time.
//
// Note that the scheduler is NOT thread-safe. The intended usage of the
// scheduler is to call CallNext sequentially. In particular, this
// allows scheduled callbacks to safely interact with the scheduler,
// for example to schedule more events.
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

// Note that this doc comment shares text with
// the package overview. Please keep in sync.

// Scheduler allows for time- and offset-based
// scheduling of function callbacks.
//
// Note that Scheduler is NOT thread-safe.
// The intended usage of Scheduler is to call
// CallNext sequentially. In particular, this
// allows scheduled callbacks to safely interact
// with the Scheduler, for example to schedule more events.
type Scheduler struct {
	heap *eventHeap
	now  time.Time
}

// Returns a new Scheduler whose
// internal clock is set to the
// zero value of time.Time.
func NewScheduler() *Scheduler {
	s := Scheduler{heap: new(eventHeap)}
	*s.heap = make([]event, 0)
	return &s
}

// Returns a new Scheduler whose
// internal clock is set to t.
func NewSchedulerTime(t time.Time) *Scheduler {
	s := Scheduler{new(eventHeap), t}
	*s.heap = make([]event, 0)
	return &s
}

// Returns the current value
// of the internal clock.
func (s *Scheduler) Now() time.Time {
	return s.now
}

// Returns whether there are
// 0 events scheduled.
func (s *Scheduler) Empty() bool {
	return s.heap.Len() < 1
}

// Schedule f to be called when
// the internal clock reaches t.
//
// Returns ErrPast if t is before
// s.Now().
func (s *Scheduler) Schedule(f func(time.Time) interface{}, t time.Time) error {
	if t.Before(s.now) {
		return ErrPast
	}
	heap.Push(s.heap, event{f, t})
	return nil
}

// Schedule f to be called when
// offset has elapsed.
//
// Returns ErrPast if offset is
// negative.
func (s *Scheduler) ScheduleOffset(f func(time.Time) interface{}, offset time.Duration) error {
	return s.Schedule(f, s.now.Add(offset))
}

// Returns the timestamp on the next
// scheduled event, or the zero value
// and ErrEmpty if no events are
// scheduled.
func (s Scheduler) PeekNext() (time.Time, error) {
	var t time.Time
	if s.Empty() {
		return t, ErrEmpty
	}
	return (*s.heap)[0].time, nil
}

// Fast-forward the internal clock
// to match the next scheduled event,
// and call the associated callback,
// passing the (now updated) time
// as the single argument. Return
// the value returned from this call.
//
// If there are no events scheduled,
// return a nil interface value and
// ErrEmpty.
//
// Note that CallNext does not modify
// s after calling the callback. Thus,
// it is safe to call methods on s
// from within the callback.
func (s *Scheduler) CallNext() (interface{}, error) {
	if s.Empty() {
		return nil, ErrEmpty
	}
	evt := heap.Pop(s.heap).(event)
	s.now = evt.time
	return evt.f(evt.time), nil
}

// Remove the next scheduled event
// from the Scheduler, but do not
// alter the internal clock.
func (s *Scheduler) RemoveNext() {
	if !s.Empty() {
		heap.Pop(s.heap)
	}
}

// Fast-forward the internal clock
// to match the next scheduled event,
// and remove the event from the
// Scheduler. If there are no events
// scheduled, do not alter the clock.
func (s *Scheduler) RemoveNextUpdate() {
	if !s.Empty() {
		evt := heap.Pop(s.heap).(event)
		s.now = evt.time
	}
}

// Remove all scheduled events from
// the Scheduler, but do not alter
// the internal clock.
func (s *Scheduler) RemoveAll() {
	*s.heap = make([]event, 0)
}

// Remove all scheduled events from
// the Scheduler, fast-forwarding
// the internal clock to match the
// latest scheduled event. If there
// are no events scheduled, do not
// alter the clock.
func (s *Scheduler) RemoveAllUpdate() {
	if !s.Empty() {
		latest := (*s.heap)[0].time
		for _, evt := range (*s.heap)[1:] {
			if evt.time.After(latest) {
				latest = evt.time
			}
		}
		s.now = latest
	}
	*s.heap = make([]event, 0)
}
