#!/bin/sh
/usr/sbin/ioreg -c IOHIDSystem | /usr/bin/awk '/HIDIdleTime/ {printf int($NF/1000000000); exit}'
