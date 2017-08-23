package main

import (
	"testing"
	"time"
)

func withIdle(d time.Duration) {
	idleDuration = func() time.Duration {
		return d
	}
}

var (
	origIdleDuration = idleDuration
	origNotify       = notify
	origNotifyBreak  = notifyBreak
	origBreakLength  = *breakLength
	origBreakAlert   = *breakAlert
	origNotifyEvery  = *notifyEvery
	origIdleAfter    = *idleAfter
)

func teardown() {
	idleDuration = origIdleDuration
	notify = origNotify
	notifyBreak = origNotifyBreak
	*breakLength = origBreakLength
	*breakAlert = origBreakAlert
	*notifyEvery = origNotifyEvery
	*idleAfter = origIdleAfter
}

func TestIdle(t *testing.T) {
	t.Run("not idle", func(t *testing.T) {
		defer teardown()
		*idleAfter = 2
		tt := newTimeTracker()
		withIdle(time.Second)
		tt.check(time.Now())
		if tt.isIdle != false {
			t.Fatal("should not be idle")
		}
	})

	t.Run("idle", func(t *testing.T) {
		defer teardown()
		*idleAfter = 2
		tt := newTimeTracker()
		withIdle(3 * time.Second)
		tt.check(time.Now())
		if tt.isIdle != true {
			t.Fatal("should be idle")
		}
	})
}

func TestBreak(t *testing.T) {
	defer teardown()
	*breakLength = 5
	tt := newTimeTracker()
	startBreak := tt.lastBreak

	withIdle(2 * time.Second)
	tt.check(time.Now())
	if !tt.lastBreak.Equal(startBreak) {
		t.Fatal("should not have changed break")
	}

	var notified bool
	notifyBreak = func() {
		notified = true
	}
	withIdle(6 * time.Second)
	tt.check(time.Now())
	if !tt.lastBreak.After(startBreak) || !notified {
		t.Fatal("should have changed break")
	}

	notified = false
	notifyBreak = func() {
		notified = true
	}
	withIdle(25 * time.Second)
	tt.check(time.Now())
	if notified {
		t.Fatal("should not have notified again")
	}
}

func TestNotify(t *testing.T) {
	t.Run("notify", func(t *testing.T) {
		defer teardown()
		tt := newTimeTracker()
		*idleAfter = 3
		*breakAlert = 5
		withIdle(0)
		notified := false
		notify = func() {
			notified = true
		}
		m := time.Now().Add(6 * time.Second)

		tt.check(m)
		if !tt.notified.Equal(m) {
			t.Fatal("expected notification")
		}
	})

	t.Run("notify interval", func(t *testing.T) {
		defer teardown()
		tt := newTimeTracker()
		*idleAfter = 3
		*breakAlert = 5
		*notifyEvery = 2
		withIdle(0)
		notified := false
		notify = func() {
			notified = true
		}
		m := time.Now().Add(6 * time.Second)
		// after 6 seconds, first notification
		tt.check(m)
		if !tt.notified.Equal(m) || !notified {
			t.Fatal("expected notification")
		}
		// after one more second, no notification
		m = m.Add(time.Second)
		notified = false
		tt.check(m)
		if tt.notified.Equal(m) || notified {
			t.Fatal("notifed too early")
		}
		// after two more seconds, second notification
		m = m.Add(2 * time.Second)
		notified = false
		tt.check(m)
		if !tt.notified.Equal(m) || !notified {
			t.Fatal("expected notification after interval")
		}
	})
}
