// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fsched

import (
	"math/rand"
	"testing"
	"time"
)

var (
	// Use these fixed values instead of time.Now()
	// since some machines fix the return value of
	// time.Now(), which can cause unexpected behavior.
	Zero          = time.Time{}
	NanoAfterZero = Zero.Add(time.Nanosecond)
)

func TestNow(t *testing.T) {
	tm := Zero
	s := NewScheduler()
	tmprime := s.Now()
	if tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}

	tm = NanoAfterZero
	s = NewSchedulerTime(tm)
	tmprime = s.Now()
	if tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}
}

func TestEmpty(t *testing.T) {
	s := NewScheduler()
	if !s.Empty() {
		t.Error("Scheduler.Empty() returned false on empty scheduler")
	}

	s.Schedule(nil, NanoAfterZero)
	if s.Empty() {
		t.Error("Scheduler.Empty() returned true on non-empty scheduler")
	}
}

func TestSchedule(t *testing.T) {
	// This test verifies that the scheduler
	// is ordering things properly by scheduling
	// f(j) at offset j. Since all j in the
	// range are used, f(j) should be the jth
	// function to be called. Calling f(j)
	// increments a global counter, and compares
	// it to j. Thus, after j calls to f, the
	// global counter should have been
	// incremented j times.

	i := 0
	f := func(j int) func(tm time.Time) interface{} {
		return func(tm time.Time) interface{} {
			// This i is the same i as in
			// TestSchedule's scope
			if i != j {
				t.Errorf("Expected i = %v; got %v", j, i)
			}
			i++
			return nil
		}
	}
	times := rand.Perm(100)
	s := NewScheduler()
	var tm time.Time
	for _, v := range times {
		s.Schedule(f(v), tm.Add(time.Duration(v)))
	}
	for ind := 0; ind < len(times); ind++ {
		s.CallNext()
	}
	if !s.Empty() {
		t.Error("Scheduler should be empty")
	}
}

func TestScheduleOffset(t *testing.T) {
	// This test verifies that the scheduler
	// is ordering things properly by scheduling
	// f(j) at offset j. Since all j in the
	// range are used, f(j) should be the jth
	// function to be called. Calling f(j)
	// increments a global counter, and compares
	// it to j. Thus, after j calls to f, the
	// global counter should have been
	// incremented j times.

	i := 0
	f := func(j int) func(tm time.Time) interface{} {
		return func(tm time.Time) interface{} {
			// This i is the same i as in
			// TestScheduleOffset's scope
			if i != j {
				t.Errorf("Expected i = %v; got %v", j, i)
			}
			i++
			return nil
		}
	}
	times := rand.Perm(100)
	s := NewScheduler()
	for _, v := range times {
		s.ScheduleOffset(f(v), time.Duration(v))
	}
	for ind := 0; ind < len(times); ind++ {
		p, _ := s.PeekNext()
		tm := Zero

		if p != tm.Add(time.Duration(ind)) {
			t.Errorf("Expected PeekNext() to return %v; returned %v", tm.Add(time.Duration(ind)), p)
		}
		s.CallNext()
	}
	if !s.Empty() {
		t.Error("Scheduler should be empty")
	}
}

func TestSchedulePast(t *testing.T) {
	t1 := Zero
	t2 := NanoAfterZero
	s := NewSchedulerTime(t2)

	err := s.Schedule(nil, t1)
	if err != ErrPast {
		t.Errorf("Expected error %v; got %v", ErrPast, err)
	}

	offset := t1.Sub(t2)
	err = s.ScheduleOffset(nil, offset)
	if err != ErrPast {
		t.Errorf("Expected error %v; got %v", ErrPast, err)
	}
}

func TestPeekNext(t *testing.T) {
	// This test verifies that the scheduler
	// is ordering things properly by scheduling
	// f(j) at offset j. Since all j in the
	// range are used, f(j) should be the jth
	// function to be called. Calling f(j)
	// increments a global counter, and compares
	// it to j. Thus, after j calls to f, the
	// global counter should have been
	// incremented j times.

	i := 0
	f := func(j int) func(tm time.Time) interface{} {
		return func(tm time.Time) interface{} {
			// This i is the same i as in
			// TestPeekNext's scope
			if i != j {
				t.Errorf("Expected i = %v; got %v", j, i)
			}
			i++
			return nil
		}
	}
	times := rand.Perm(100)
	s := NewScheduler()
	for _, v := range times {
		s.ScheduleOffset(f(v), time.Duration(v))
	}
	for ind := 0; ind < len(times); ind++ {
		p, _ := s.PeekNext()
		tm := Zero

		if p != tm.Add(time.Duration(ind)) {
			t.Errorf("Expected PeekNext() to return %v; returned %v", tm.Add(time.Duration(ind)), p)
		}
		s.CallNext()
	}
	if !s.Empty() {
		t.Error("Scheduler should be empty")
	}
}

func TestPeekNextError(t *testing.T) {
	s := NewScheduler()
	_, err := s.PeekNext()
	if err != ErrEmpty {
		t.Errorf("Expected error %v; got %v", ErrEmpty, err)
	}
}

func TestCallNextError(t *testing.T) {
	s := NewScheduler()
	_, err := s.CallNext()
	if err != ErrEmpty {
		t.Errorf("Expected error %v; got %v", ErrEmpty, err)
	}
}

func TestRemoveNext(t *testing.T) {
	s := NewScheduler()
	s.Schedule(nil, Zero)
	s.RemoveNext()
	if !s.Empty() {
		t.Errorf("Scheduler should be empty")
	}
}

func TestRemoveNextUpdate(t *testing.T) {
	s := NewScheduler()
	s.Schedule(nil, Zero)
	s.RemoveNextUpdate()
	if !s.Empty() {
		t.Errorf("Scheduler should be empty")
	}

	s = NewScheduler()
	tm := NanoAfterZero
	s.Schedule(nil, tm)
	s.RemoveNextUpdate()
	if tmprime := s.Now(); tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}

	tm = NanoAfterZero
	s = NewSchedulerTime(tm)
	s.RemoveNextUpdate()
	if tmprime := s.Now(); tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}
}

func TestRemoveAll(t *testing.T) {
	s := NewScheduler()
	s.Schedule(nil, Zero)
	s.RemoveAll()
	if !s.Empty() {
		t.Errorf("Scheduler should be empty")
	}
}

func TestRemoveAllUpdate(t *testing.T) {
	s := NewScheduler()
	s.Schedule(nil, Zero)
	s.RemoveAllUpdate()
	if !s.Empty() {
		t.Errorf("Scheduler should be empty")
	}

	s = NewScheduler()
	tm := NanoAfterZero
	s.Schedule(nil, tm)
	s.RemoveAllUpdate()
	if tmprime := s.Now(); tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}

	tm = NanoAfterZero
	s = NewSchedulerTime(tm)
	if tmprime := s.Now(); tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}
}
