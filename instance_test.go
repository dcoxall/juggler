package juggler

import (
	"github.com/dcoxall/juggler/utils"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestStartingStoppingInstance(t *testing.T) {
	port := <-utils.FindAvailablePort()
	instance := NewInstance(port, "startstop")
	if instance.Started() {
		t.Errorf("Instance should not have started")
	}
	ready, err := instance.Start()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	timeout := time.After(2 * time.Second)
	select {
	case state := <-ready:
		if state != Running {
			t.Fatalf("instance returned %d instead of %d", state, Running)
		}
	case <-timeout:
		instance.ForceStop()
		t.Fatalf("instance state not updated")
	}
	if !instance.Started() {
		t.Errorf("Instance claims to have not started")
	}
	if instance.Stopped() {
		t.Errorf("Instance should not be stopped")
	}
	stopped, err := instance.Stop()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	timeout = time.After(2 * time.Second)
	select {
	case state := <-stopped:
		if state != Stopped {
			t.Fatalf("instance returned %d instead of %d", state, Stopped)
		}
	case <-timeout:
		instance.ForceStop()
		t.Fatalf("instance never stopped before timeout period")
	}
	if !instance.Stopped() {
		t.Errorf("Instance claims to have not stopped")
	}
}

func TestInstanceStartErrors(t *testing.T) {
	port := <-utils.FindAvailablePort()
	instance := NewInstance(port, "starterrors")
	ready, err := instance.Start()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	timeout := time.After(2 * time.Second)
	if _, err := instance.Start(); err == nil {
		t.Errorf("Expected an error when starting an already started instance")
	}
	select {
	case <-ready:
	case <-timeout:
		instance.ForceStop()
		t.Fatalf("instance never started before timeout period")
	}
	stopped, _ := instance.Stop()
	timeout = time.After(2 * time.Second)
	select {
	case <-stopped:
	case <-timeout:
		instance.ForceStop()
		t.Fatalf("instance never stopped within timeout period")
	}
}

func TestInstanceStopErrors(t *testing.T) {
	port := <-utils.FindAvailablePort()
	instance := NewInstance(port, "stoperrors")
	if _, err := instance.Stop(); err == nil {
		t.Errorf("Expected an error when stopping an already stopped instance")
	}
}

func TestInstanceProxying(t *testing.T) {
	var wg sync.WaitGroup
	instances := map[string]*Instance{
		"foo": NewInstance(<-utils.FindAvailablePort(), "foo"),
		"bar": NewInstance(<-utils.FindAvailablePort(), "bar"),
	}
	for ref, i := range instances {
		wg.Add(1)
		go func(instance *Instance) {
			ready, _ := instance.Start()
			timeout := time.After(2 * time.Second)
			select {
			case <-ready:
				wg.Done()
			case <-timeout:
				wg.Done()
				t.Fatalf("State change failed to occur within timeout for %s", ref)
			}
		}(i)
	}

	wg.Wait() // wait until our 2 instances have started

	// make a request to both instances and assert different responses
	for ref, i := range instances {
		proxy, _ := i.ReverseProxy()
		req, _ := http.NewRequest("GET", "http://example.com/", nil)
		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
		if body := w.Body.String(); body != ref {
			t.Errorf("Expected response body to be %v but got %v", ref, body)
		}

		wg.Add(1)
		go func(instance *Instance) {
			stopped, _ := instance.Stop()
			timeout := time.After(2 * time.Second)
			select {
			case <-stopped:
				wg.Done()
			case <-timeout:
				wg.Done()
				t.Fatalf("State change failed to occur within timeout for %s", ref)
			}
		}(i)
	}

	wg.Wait() // wait until instances have been removed
}
