# flippo

A simple utility that let's you know when it's time to take a break from the computer. Made for OS X.  Must run in foreground (otherwise you won't get notified).

By default a notification comes after 40 minutes of not being idle (system-wide) and bugs you every minute until you become idle. You are considered idle after 10 seconds of no activity. A break must be at least 10 minutes.

The defaults are configurable via flags.

```
Usage of flippo:
  -break int
    	break alert interval (seconds) (default 2400)
  -debug
    	verbose display
  -freq int
    	notification frequency (seconds) (default 60)
  -idle-after int
    	time after which user is considered idle (seconds) (default 10)
  -min-break int
    	break length (seconds) (default 600)
  -sound string
    	sound name (from ~/Library/Sounds or /System/Library/Sounds) (default "Hero")
  -sound-after string
    	sound name after break (from ~/Library/Sounds or /System/Library/Sounds) (default "Purr")
```
