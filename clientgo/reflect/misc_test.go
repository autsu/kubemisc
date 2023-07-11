package main

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
)

func TestBackoff(t *testing.T) {
	realClock := &clock.RealClock{}
	manager := wait.NewExponentialBackoffManager(800*time.Millisecond, 30*time.Second, 2*time.Minute, 2.0, 1.0, realClock)
	wait.BackoffUntil(func() {
		t.Log("tick")
	}, manager, true, wait.NeverStop)
}
