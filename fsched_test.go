package fsched

import (
	"math/rand"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	// Zero value
	var tm time.Time
	s := MakeScheduler()
	tmprime := s.Now()
	if tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}

	tm = time.Now()
	s = MakeSchedulerTime(tm)
	tmprime = s.Now()
	if tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}
}

func TestEmpty(t *testing.T) {
	s := MakeScheduler()
	if !s.Empty() {
		t.Error("Scheduler.Empty() returned false on empty scheduler")
	}
	s.Schedule(nil, time.Now())
	if s.Empty() {
		t.Error("Scheduler.Empty() returned true on non-empty scheduler")
	}
}

func TestSchedule(t *testing.T) {
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
	s := MakeScheduler()
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
	s := MakeScheduler()
	for _, v := range times {
		s.ScheduleOffset(f(v), time.Duration(v))
	}
	for ind := 0; ind < len(times); ind++ {
		p, _ := s.PeekNext()
		var tm time.Time

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
	t1 := time.Now()
	t2 := time.Now()
	s := MakeSchedulerTime(t2)

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
	s := MakeScheduler()
	for _, v := range times {
		s.ScheduleOffset(f(v), time.Duration(v))
	}
	for ind := 0; ind < len(times); ind++ {
		p, _ := s.PeekNext()
		var tm time.Time

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
	s := MakeScheduler()
	_, err := s.PeekNext()
	if err != ErrEmpty {
		t.Errorf("Expected error %v; got %v", ErrEmpty, err)
	}
}
