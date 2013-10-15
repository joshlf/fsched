package eventheap

import (
	"math/rand"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	// Zero value
	var tm time.Time
	eh := MakeEventHeap()
	tmprime := eh.Now()
	if tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}

	tm = time.Now()
	eh = MakeEventHeapTime(tm)
	tmprime = eh.Now()
	if tmprime != tm {
		t.Errorf("Expected time %v; got %v", tm, tmprime)
	}
}

func TestEmpty(t *testing.T) {
	eh := MakeEventHeap()
	if !eh.Empty() {
		t.Error("EventHeap.Empty() returned false on empty event heap")
	}
	eh.Schedule(nil, time.Now())
	if eh.Empty() {
		t.Error("EventHeap.Empty() returned true on non-empty event heap")
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
	eh := MakeEventHeap()
	var tm time.Time
	for _, v := range times {
		eh.Schedule(f(v), tm.Add(time.Duration(v)))
	}
	for ind := 0; ind < len(times); ind++ {
		eh.CallNext()
	}
	if !eh.Empty() {
		t.Error("Event heap should be empty")
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
	eh := MakeEventHeap()
	for _, v := range times {
		eh.ScheduleOffset(f(v), time.Duration(v))
	}
	for ind := 0; ind < len(times); ind++ {
		p, _ := eh.PeekNext()
		var tm time.Time

		if p != tm.Add(time.Duration(ind)) {
			t.Errorf("Expected PeekNext() to return %v; returned %v", tm.Add(time.Duration(ind)), p)
		}
		eh.CallNext()
	}
	if !eh.Empty() {
		t.Error("Event heap should be empty")
	}
}

func TestSchedulePast(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now()
	eh := MakeEventHeapTime(t2)

	err := eh.Schedule(nil, t1)
	if err != ErrPast {
		t.Errorf("Expected error %v; got %v", ErrPast, err)
	}

	offset := t1.Sub(t2)
	err = eh.ScheduleOffset(nil, offset)
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
	eh := MakeEventHeap()
	for _, v := range times {
		eh.ScheduleOffset(f(v), time.Duration(v))
	}
	for ind := 0; ind < len(times); ind++ {
		p, _ := eh.PeekNext()
		var tm time.Time

		if p != tm.Add(time.Duration(ind)) {
			t.Errorf("Expected PeekNext() to return %v; returned %v", tm.Add(time.Duration(ind)), p)
		}
		eh.CallNext()
	}
	if !eh.Empty() {
		t.Error("Event heap should be empty")
	}
}

func TestPeekNextError(t *testing.T) {
	eh := MakeEventHeap()
	_, err := eh.PeekNext()
	if err != ErrEmpty {
		t.Errorf("Expected error %v; got %v", ErrEmpty, err)
	}
}
