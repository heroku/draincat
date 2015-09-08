package main

import (
	"math/rand"
	"regexp"
	"strconv"
)

type Latency struct {
	p99 int // in milliseconds
	p95 int
	p50 int
}

func NewLatencyFromSpec(spec string) (*Latency, error) {
	r := regexp.MustCompile(
		`(?P<p99>\d+)\:(?P<p95>\d+)\:(?P<p50>\d+)`).FindAllStringSubmatch(spec, -1)[0]
	p99, err := strconv.Atoi(r[1])
	if err != nil {
		return nil, err
	}
	p95, err := strconv.Atoi(r[2])
	if err != nil {
		return nil, err
	}
	p50, err := strconv.Atoi(r[3])
	if err != nil {
		return nil, err
	}
	return NewLatency(p99, p95, p50), nil
}

func NewLatency(p99, p95, p50 int) *Latency {
	return &Latency{p99, p95, p50}
}

func (l *Latency) Create() int {
	loc := rand.Intn(100)
	switch {
	case loc == 99:
		return l.p99
	case loc >= 95:
		return l.p95
	default:
		return l.p50
	}
}
