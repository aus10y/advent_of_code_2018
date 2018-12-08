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
	// Of all guards, which guard is most frequently asleep on the same minute?

	habits := make(map[int]map[int]int)

	{
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

				for m := fellAsleep; m < awoke; m++ {
					guardMap := habits[m]
					if guardMap == nil {
						guardMap = make(map[int]int)
					}

					guardMap[guard]++
					habits[m] = guardMap
				}
			default:
				log.Fatalf("unanticipated action")
			}
		}
	}

	asleep := 0
	minute := 0
	sleepiest := 0
	for m, habit := range habits {
		for g, s := range habit {
			if s > asleep {
				asleep = s
				sleepiest = g
				minute = m
			}
		}
	}

	fmt.Printf("Guard ID: %d, most often asleep during minute %d, Answer: %d\n", sleepiest, minute, sleepiest*minute)
}
