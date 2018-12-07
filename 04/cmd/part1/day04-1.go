package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Action int

const (
	Begins Action = iota
	Sleeps
	Wakes
)

type Record struct {
	time   time.Time
	guard  int
	action Action
}

type recordSlice []Record

func (p recordSlice) Len() int {
	return len(p)
}

func (p recordSlice) Less(i, j int) bool {
	a, b := p[i].time, p[j].time
	return a.Before(b) || (a.Equal(b) && i < j)
}

func (p recordSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func parseRecord(line string) (record Record, err error) {
	// Parse the Time:
	// Ex. "[1518-10-15 00:29] wakes up"
	t, err := time.Parse("2006-01-02 15:04", line[1:17])
	if err != nil {
		return record, err
	}
	record.time = t

	// Parse the Guard ID, if present:
	// Ex., ..." Guard #503 begins shift"
	start := strings.Index(line, "#")
	if start != -1 {
		end := strings.Index(line, " begins")
		record.guard, err = strconv.Atoi(line[start+1 : end])
		if err != nil {
			return record, err
		}
		record.action = Begins
		return
	}

	// Parse the Action:
	if strings.Index(line, "falls asleep") != -1 {
		record.action = Sleeps
	} else if strings.Index(line, "wakes up") != -1 {
		record.action = Wakes
	} else {
		return record, errors.New("could not parse action")
	}

	return
}

func parseRecords(lines []string) recordSlice {
	records := make(recordSlice, len(lines))
	for i, line := range lines {
		record, err := parseRecord(line)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		records[i] = record
	}
	return records
}

type guardHabit struct {
	asleep int
	times  map[int]int
}

func main() {
	var records recordSlice
	{
		input, err := ioutil.ReadFile("../../input.txt")
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(string(input), "\n")
		records = parseRecords(lines)
	}

	sort.Sort(records)

	// Strategy:
	// Find the guard that has the most minutes asleep.
	// What minute does that guard spend asleep the most?
	//
	// Walk the records, for each guard,
	// Find their next period of sleep,
	// Add the minutes asleep to that guard's cummulative sleep,
	// Increment the guardHabit.times counter for each individual minute of sleep.
	habits := make(map[int]guardHabit)

	guard := 0
	fellAsleep := 0
	awoke := 0
	for _, record := range records {
		if record.guard != 0 {
			guard = record.guard

			// Records containing the guard ID are always the 'Begins' action.
			continue
		}

		switch record.action {
		case Sleeps:
			fellAsleep = record.time.Minute()
		case Wakes:
			awoke = record.time.Minute()

			for i := fellAsleep; i < awoke; i++ {
				habit := habits[guard]
				if habit.times == nil {
					habit.times = make(map[int]int)
				}
				habit.asleep++
				habit.times[i] = habit.times[i] + 1
				habits[guard] = habit
			}
		default:
			log.Fatalf("unanticipated action")
		}
	}

	// Find the sleepiest guard.
	sleepiest := 0
	{
		mostSleep := 0
		for guard, habit := range habits {
			if habit.asleep > mostSleep {
				sleepiest = guard
				mostSleep = habit.asleep
			}
		}
	}

	// Find the minute during which the guard is most often asleep.
	minute := 0
	{
		sleep := 0
		for m, s := range habits[sleepiest].times {
			if s > sleep {
				minute = m
				sleep = s
			}
		}
	}

	fmt.Printf("Guard ID: %d, most often asleep during minute %d, Answer: %d\n", sleepiest, minute, sleepiest*minute)
}
