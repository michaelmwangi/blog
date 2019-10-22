#!/bin/bash
echo  "starting ...."
taskset -c 1 go run main.go &
sudo perf stat -C 0  -- sleep 1

