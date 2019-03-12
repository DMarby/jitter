package jitter

import (
	"fmt"
	"math/rand"
	"time"
)

// Ticker is a ticker that emits events on a channel at the given interval, with an added delay up to the defined max jitter
// If the receiever doesn't keep up the events will be discarded
type Ticker struct {
	C <-chan time.Time // Channel which the events are delivered on

	interval time.Duration // Interval for the ticker to run at
	jitter   time.Duration // Max jitter to add to the interval

	stop   chan bool  // Channel used for stopping the timer
	random *rand.Rand // Local random for generating jitter
}

// NewTicker returns a new ticker with the given interval and jitter
func NewTicker(interval time.Duration, jitter time.Duration) *Ticker {
	if interval <= 0 {
		panic(fmt.Errorf("non-positive interval for NewTicker: %d", int(interval)))
	}

	if jitter <= 0 {
		panic(fmt.Errorf("non-positive jitter for NewTicker: %d", int(jitter)))
	}

	// Create a seeded random to use for the jitter
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	// Create a buffered channel for tick events
	c := make(chan time.Time, 1)
	ticker := &Ticker{
		C: c,

		interval: interval,
		jitter:   jitter,

		stop:   make(chan bool),
		random: random,
	}

	// Run the ticker
	// Ticker.C is a read-only channel, so we need to pass it
	go ticker.tick(c)

	return ticker
}

func (t Ticker) tick(c chan time.Time) {
loop:
	for {
		t.sleep() // Sleep for duration + jitter

		select {
		case <-t.stop: // Check for the stop signal and stop
			break loop
		case c <- time.Now(): // Send the time event to the ticker channel
		default: // Fall-through so that sending to the channel doesn't block
		}
	}
}

func (t Ticker) sleep() {
	jitter := int64(t.jitter)
	delay := time.Duration(t.random.Int63n(jitter))
	time.Sleep(t.interval + delay)
}

// Stop will stop the ticker. It does not close the channel.
func (t Ticker) Stop() {
	t.stop <- true
}
