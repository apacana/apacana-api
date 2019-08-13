package helper

import "time"

type Tick struct {
	eventname string
	delayms   time.Duration
}

type TickerInfo struct {
	TickerOptions
	tickList         []*Tick
	largeLatencyList []*Tick
	latesttme        time.Time
	starttime        time.Time
}

type TickerOptions struct {
	LargeLatencyThreshold time.Duration
}

func NewTickerInfo() *TickerInfo {
	now := time.Now()
	ti := &TickerInfo{
		tickList:         make([]*Tick, 0),
		largeLatencyList: make([]*Tick, 0),
		latesttme:        now,
		starttime:        now,
	}
	return ti
}

func (ti *TickerInfo) TickMs() float64 {
	currenttime := time.Now()
	interval := float64(currenttime.Sub(ti.latesttme).Nanoseconds()) / float64(1000*1000)
	ti.latesttme = currenttime
	return interval
}
