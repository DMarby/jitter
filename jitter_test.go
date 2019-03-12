package jitter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DMarby/jitter"
)

func TestNewTicker(t *testing.T) {
	t.Run("panics on non-positive interval", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewTicker did not panic on non-positive interval")
			}
		}()

		jitter.NewTicker(0, time.Second)
	})

	t.Run("panics on non-positive jitter", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewTicker did not panic on non-positive jitter")
			}
		}()

		jitter.NewTicker(time.Second, 0)
	})
}

func TestInitialDelay(t *testing.T) {
	delay := time.Millisecond * 100
	ticker := jitter.NewTicker(delay, time.Millisecond)
	defer ticker.Stop()

	start := time.Now()
	<-ticker.C
	if time.Now().Sub(start) < delay {
		t.Errorf("ticked too early")
	}
}

func TestJitter(t *testing.T) {
	delay := 100 * time.Millisecond

	ticker := jitter.NewTicker(delay, delay)
	defer ticker.Stop()

	regularTicker := time.NewTicker(delay)
	defer regularTicker.Stop()

	ltHalf := false
	gtHalf := false

	// Test to make sure we get jitter values below and above half of the max possible jitter
	for i := 0; i < 15; i++ {
		diff := (<-ticker.C).Sub(<-regularTicker.C)
		if diff < delay/2 {
			ltHalf = true
		} else if diff > delay/2 {
			gtHalf = true
		}
	}

	if !ltHalf {
		t.Errorf("No jitter less then half of max")
	}

	if !gtHalf {
		t.Errorf("No jitter greater then half of max")
	}
}

func Example() {
	t := jitter.NewTicker(
		time.Second,
		time.Second*10,
	)

	for tick := range t.C {
		fmt.Println("Tick at", tick)
	}
}
