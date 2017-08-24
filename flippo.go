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
	timeUnit   = time.Minute
)

var (
	sound       = flag.String("sound", "Hero", "sound name (from ~/Library/Sounds or /System/Library/Sounds)")
	soundAfter  = flag.String("sound-after", "Purr", "sound name after break (from ~/Library/Sounds or /System/Library/Sounds)")
	breakLength = flag.Int64("min-break", 10, "break length (minutes)")
	breakAlert  = flag.Int64("break", 40, "break alert interval (minutes)")
	notifyEvery = flag.Int64("freq", 1, "notification frequency (minutes)")
	idleAfter   = flag.Int64("idle-after", 10, "time after which user is considered idle (seconds)")
)

var idleDuration = func() time.Duration {
	out, err := exec.Command("./idle_time.sh").Output()
	if err != nil {
		log.Printf("%+v", err)
	}
	sec, err := strconv.ParseInt(string(out), 10, 64)
	if err != nil {
		log.Printf("%+v", err)
	}
	return time.Duration(sec) * time.Second
}

var notify = func() {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "%s"`, title, body, *sound)
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Fatal("%+v", err)
	}
}

var notifyBreak = func() {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "%s"`, titleBreak, bodyBreak, *soundAfter)
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Fatal("%+v", err)
	}
}

type timeTracker struct {
	lastBreak time.Time
	notified  time.Time
	isIdle    bool
	inBreak   bool
}

func newTimeTracker() *timeTracker {
	return &timeTracker{
		lastBreak: time.Now(),
		notified:  time.Now(),
	}
}

func (tt *timeTracker) check(t time.Time) {
	sinceLast := t.Sub(tt.lastBreak)
	lastNotified := t.Sub(tt.notified)
	idle := idleDuration()

	inBreak := idle > time.Duration(*breakLength)*timeUnit
	if inBreak {
		if !tt.inBreak {
			notifyBreak()
		}
		tt.lastBreak = t
	}
	tt.inBreak = inBreak
	tt.isIdle = idle > time.Duration(*idleAfter)*time.Second
	if !tt.isIdle && sinceLast > time.Duration(*breakAlert)*timeUnit &&
		lastNotified > time.Duration(*notifyEvery)*timeUnit {
		notify()
		tt.notified = t
	}
}

func main() {
	flag.Parse()
	tracker := newTimeTracker()

	for {
		time.Sleep(100 * time.Millisecond)
		tracker.check(time.Now())
	}
}
