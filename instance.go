package juggler

import (
	"fmt"
	"github.com/dcoxall/juggler/utils"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"time"
)

// Instance is an Instancer structure that allows custom applications to be
// managed in asynchronously.
type Instance struct {
	port         int
	ref          string
	state        InstanceState
	stateChan    chan InstanceState
	forceStop    chan int
	cmd          *exec.Cmd
	proxy        *httputil.ReverseProxy
	bootstrapper Bootstrapper
}

// NewInstance takes a Bootstrapper and a reference string and returns an
// appropriately configured Instance.
func NewInstance(bootstrapper Bootstrapper, ref string) *Instance {
	return &Instance{
		port:         <-utils.FindAvailablePort(),
		ref:          ref,
		state:        Stopped,
		stateChan:    make(chan InstanceState, 1),
		forceStop:    make(chan int, 1),
		bootstrapper: bootstrapper,
	}
}

// Returns true if the Instance has started and is accepting requests.
func (i *Instance) Started() bool {
	return i.state == Running
}

// Returns true if the Instance has stopped and is not accepting requests.
func (i *Instance) Stopped() bool {
	return i.state == Stopped
}

// Returns true if the Instance is stopping and requests should stop being made.
func (i *Instance) Stopping() bool {
	return i.state == Stopping
}

// Returns true if the Instance is starting.
func (i *Instance) Starting() bool {
	return i.state == Starting
}

// Start will attempt to start the underlying web process returning an error if
// it is unable to. As it supports asynchronous startup the returned channel can
// be used to determine when the Instance state has changed to Running.
func (i *Instance) Start() (<-chan InstanceState, error) {
	// We can't start if we aren't stopped
	if !i.Stopped() {
		return i.stateChan, fmt.Errorf("Unable to start")
	}

	i.state = Starting

	// consider using sync.Once
	cmd, err := i.bootstrapper.Bootstrap(i.port, i.ref)
	i.cmd = cmd
	if err != nil {
		return i.stateChan, err
	}

	// in the background let's wait until we can connect and then trigger
	// completion on the channel
	go func(instance *Instance) {
		for !instance.Started() {
			select {
			case <-instance.forceStop:
				return
			default:
				if utils.IsPortFree(instance.port) {
					time.Sleep(500 * time.Millisecond)
				} else {
					instance.state = Running
				}
			}
		}
		instance.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", instance.port),
		})
		instance.stateChan <- instance.state
	}(i)

	// return our channel that indicates a change in state
	return i.stateChan, nil
}

// Stop will attempt to stop the underlying web process using the InstanceState
// channel to signal success. An error can be returned immediately if stopping
// can not be performed at that point.
func (i *Instance) Stop() (<-chan InstanceState, error) {
	// We can't start if we aren't stopped
	if !i.Started() {
		return i.stateChan, fmt.Errorf("Unable to stop")
	}
	i.state = Stopping

	// Begin watching this process and signal when it ends
	go func(instance *Instance) {
		instance.cmd.Wait()
		instance.state = Stopped
		instance.stateChan <- instance.state
	}(i)

	// send the kill signal to actually stop the process
	if err := i.cmd.Process.Signal(os.Kill); err != nil {
		return i.stateChan, err
	}

	// return our channel that indicates a change in state
	return i.stateChan, nil
}

// ForceStop will forcefully kill the underlying web process.
func (i *Instance) ForceStop() error {
	i.forceStop <- 1
	return i.cmd.Process.Kill()
}

// ReverseProxy returns a structure that can be used to make requests to the
// underlying web process. An error is returned if the Instance is not in the
// correct state.
func (i *Instance) ReverseProxy() (*httputil.ReverseProxy, error) {
	if !i.Started() {
		return nil, fmt.Errorf("Instance not yet started")
	}
	return i.proxy, nil
}
