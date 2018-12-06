package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

func parseRecord(line string) (record Record, err error) {
	// Parse the Time:
	// Ex. [1518-10-15 00:29] wakes up
	t, err := time.Parse("2006-01-02 15:04", line[1:17])
	if err != nil {
		return record, err
	}
	record.time = t

	// Parse the Guard ID, if present:
	// Guard #503 begins shift
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

func main() {
	input, err := ioutil.ReadFile("../../input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")

	records := make([]Record, len(lines))
	for i, line := range lines {
		record, err := parseRecord(line)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		records[i] = record
		fmt.Printf("%v\n", record)
	}

	// sort the records
}
