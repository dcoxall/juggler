package juggler

import (
	"github.com/dcoxall/juggler/utils"
	"testing"
	"time"
)

func TestStartingStoppingInstance(t *testing.T) {
	port := <-utils.FindAvailablePort()
	instance := NewInstance(port, "pong")
	ready, err := instance.Start()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	timeout := time.After(5 * time.Second)
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
	timeout = time.After(5 * time.Second)
	select {
	case state := <-stopped:
		if state != Stopped {
			t.Fatalf("instance returned %d instead of %d", state, Stopped)
		}
	case <-timeout:
		t.Fatalf("instance never stopped before timeout period")
	}
}

func TestInstanceStartErrors(t *testing.T) {
	port := <-utils.FindAvailablePort()
	instance := NewInstance(port, "pong")
	ready, _ := instance.Start()
	timeout := time.After(5 * time.Second)
	if _, err := instance.Start(); err == nil {
		t.Errorf("Expected an error when starting an already started instance")
	}
	select {
	case <-ready:
	case <-timeout:
		t.Fatalf("instance never started before timeout period")
	}
	instance.Stop()
	timeout = time.After(5 * time.Second)
	select {
	case <-ready:
	case <-timeout:
		t.Fatalf("instance never stopped within timeout period")
	}
}

func TestInstanceStopErrors(t *testing.T) {
	port := <-utils.FindAvailablePort()
	instance := NewInstance(port, "pong")
	if _, err := instance.Stop(); err == nil {
		t.Errorf("Expected an error when stopping an already stopped instance")
	}
}
