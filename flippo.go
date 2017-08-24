package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

const (
	title      = "Take a break"
	body       = "40 minutes passed since your last."
	titleBreak = "You took a break"
	bodyBreak  = "Good!"
)

var (
	sound       = flag.String("sound", "Hero", "sound name (from ~/Library/Sounds or /System/Library/Sounds)")
	soundBreak  = flag.String("sound-after", "Purr", "sound name after break (from ~/Library/Sounds or /System/Library/Sounds)")
	breakLength = flag.Int64("min-break", 600, "break length (seconds)")
	breakAlert  = flag.Int64("break", 2400, "break alert interval (seconds)")
	notifyEvery = flag.Int64("freq", 60, "notification frequency (seconds)")
	idleAfter   = flag.Int64("idle-after", 10, "time after which user is considered idle (seconds)")
	debug       = flag.Bool("debug", false, "verbose display")
)

// idleDuration returns system idle time
var idleDuration = func() time.Duration {
	awk := exec.Command("/usr/bin/awk", "/HIDIdleTime/ {printf int($NF/1000000000); exit}")
	stdin, err := awk.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	ioreg, err := exec.Command("/usr/sbin/ioreg", "-c", "IOHIDSystem").Output()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		_, err := stdin.Write(ioreg)
		if err != nil {
			log.Fatal(err)
		}
	}()
	out, err := awk.Output()
	if err != nil {
		log.Fatal(err)
	}
	sec, err := strconv.ParseInt(string(out), 10, 64)
	if err != nil {
		log.Printf("%+v", err)
	}
	return time.Duration(sec) * time.Second
}

// notify notifies the using the given title, body and system sound
var notify = func(title, body, sound string) {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "%s"`, title, body, sound)
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Fatal("%+v", err)
	}
}

type timeTracker struct {
	// lastBreak is the time when the most recent break was interrupted,
	// as in the user became active again
	lastBreak time.Time
	// notified is the last time the user was sent a break notification
	notified time.Time
	// isIdle is true if the user is idle
	isIdle bool
	// inBreak will be true if the user is in a break
	inBreak bool
}

// newTimeTracker creates a new tracker to track the users activity
func newTimeTracker() *timeTracker {
	return &timeTracker{
		lastBreak: time.Now(),
		notified:  time.Now(),
	}
}

// check checks the users activity status at time t.
func (tt *timeTracker) check(t time.Time) {
	sinceLast := t.Sub(tt.lastBreak)
	lastNotified := t.Sub(tt.notified)
	idle := idleDuration()
	secs := func(s *int64) time.Duration {
		return time.Duration(*s) * time.Second
	}
	if *debug {
		log.Printf("Idle: %ds", idle)
	}
	inBreak := idle > secs(breakLength)
	if inBreak {
		if !tt.inBreak {
			notify(titleBreak, bodyBreak, *soundBreak)
			if *debug {
				log.Println("Completed break.")
			}
		}
		if *debug {
			log.Println("In break.")
		}
		tt.lastBreak = t
	}
	tt.inBreak = inBreak
	tt.isIdle = idle > secs(idleAfter)
	if !tt.isIdle && sinceLast > secs(breakAlert) && lastNotified > secs(notifyEvery) {
		if *debug {
			log.Println("Notified to take break.")
		}
		notify(title, body, *sound)
		tt.notified = t
	}
}

func main() {
	flag.Parse()
	tracker := newTimeTracker()
	for {
		time.Sleep(time.Second)
		tracker.check(time.Now())
	}
}
