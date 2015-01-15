package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Event struct {
	Event    string
	PeerID   string
	Session  string
	System   string
	Time     time.Time
	Duration time.Duration
}

func ParseEvent(evstr string) (Event, error) {
	var ev Event
	err := json.Unmarshal([]byte(evstr), &ev)
	if err != nil {
		return ev, err
	}

	return ev, nil
}

func LoadEvents(file string) ([]Event, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var events []Event
	scan := bufio.NewScanner(fi)
	for scan.Scan() {
		ev, err := ParseEvent(scan.Text())
		if err != nil {
			return nil, err
		}

		events = append(events, ev)
	}

	return events, nil
}

func Filter(evs []Event, f func(e Event) bool) []Event {
	out := make([]Event, 0, len(evs))
	for _, e := range evs {
		if f(e) {
			out = append(out, e)
		}
	}
	return out
}

type EventSorter struct {
	events   []Event
	lesscomp func(a, b Event) bool
	sort.Interface
}

func (e *EventSorter) Len() int {
	return len(e.events)
}

func (e *EventSorter) Swap(i, j int) {
	e.events[i], e.events[j] = e.events[j], e.events[i]
}

func (e *EventSorter) Less(i, j int) bool {
	return e.lesscomp(e.events[i], e.events[j])
}

func main() {
	events, err := LoadEvents(os.Args[1])
	if err != nil {
		panic(err)
	}

	events = Filter(events, func(e Event) bool {
		return e.Duration > 0
	})

	sort.Sort(&EventSorter{
		events: events,
		lesscomp: func(a, b Event) bool {
			return a.Duration < b.Duration
		},
	})

	for _, e := range events {
		fmt.Printf("%s: %s\n", e.Event, e.Duration)
	}
}
