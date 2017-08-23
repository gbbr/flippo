package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	title = flag.String("title", "Take a break", "Notification title")
	body  = flag.String("body", "40 minutes passed since your last.", "Notification body")
	sound = flag.String("sound", "Blow", "Sound name (from ~/Library/Sounds or /System/Library/Sounds)")
)

const (
	minBreakLength = 15 * time.Second
	breakAlert     = 10 * time.Second
	notifyEvery    = 5 * time.Second
	isIdleAfter    = 2 * time.Second
)

func init() {
	flag.Parse()
}

func idleDuration() time.Duration {
	out, err := exec.Command("./idle_time.sh").Output()
	if err != nil {
		log.Printf("%+v", err)
	}
	s := strings.TrimRight(string(out), "\n")
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("%+v", err)
	}
	return time.Duration(sec) * time.Second
}

func notify() {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "%s"`, *title, *body, *sound)
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		log.Fatal("%+v", err)
	}
}

func main() {
	lastBreak := time.Now()
	notified := time.Now()
	isIdle := false

	for {
		time.Sleep(100 * time.Millisecond)
		sinceLast := time.Now().Sub(lastBreak)
		idle := idleDuration()

		if !isIdle && sinceLast > breakAlert && time.Now().Sub(notified) > notifyEvery {
			notify()
			notified = time.Now()
		}

		if idle > isIdleAfter {
			log.Println("idle")
			isIdle = true
		} else {
			isIdle = false
			log.Println("not idle")
		}

		if idle > minBreakLength {
			lastBreak = time.Now()
			log.Println("break")
		}
	}
}
