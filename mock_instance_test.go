package juggler

import (
	"testing"
	"time"
)

func TestStartingStoppingInstance(t *testing.T) {
	instance := NewInstance()
	ready, err := instance.Start()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	timeout := time.After(10 * time.Second)
	select {
	case state := <-ready:
		if state != Running {
			t.Fatalf("instance returned %d instead of %d", state, Running)
		}
	case <-timeout:
		t.Fatalf("instance never started before timeout period")
	}
	stopped, err := instance.Stop()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	timeout = time.After(10 * time.Second)
	select {
	case state := <-stopped:
		if state != Stopped {
			t.Fatalf("instance returned %d instead of %d", state, Stopped)
		}
	case <-timeout:
		t.Fatalf("instance never stopped before timeout period")
	}
}
