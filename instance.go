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

const (
	Stopped = iota
	Starting
	Running
	Stopping
)

type Instance struct {
	port      int
	ref       string
	state     int
	stateChan chan int
	forceStop chan int
	cmd       *exec.Cmd
	proxy     *httputil.ReverseProxy
}

func NewInstance(ref string) *Instance {
	return &Instance{
		port:      <-utils.FindAvailablePort(),
		ref:       ref,
		state:     Stopped,
		stateChan: make(chan int, 1),
		forceStop: make(chan int, 1),
	}
}

func (i *Instance) Started() bool {
	return i.state == Running
}

func (i *Instance) Stopped() bool {
	return i.state == Stopped
}

func (i *Instance) Start() (<-chan int, error) {
	// We can't start if we aren't stopped
	if !i.Stopped() {
		return i.stateChan, fmt.Errorf("Unable to start")
	}

	// store and start the command
	i.cmd = exec.Command(
		"/tmp/ping",
		fmt.Sprintf(":%d", i.port),
		i.ref,
	)
	i.cmd.Stdout = os.Stdout
	i.cmd.Stderr = os.Stderr
	if err := i.cmd.Start(); err != nil {
		return i.stateChan, err
	}
	i.state = Starting

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

func (i *Instance) Stop() (<-chan int, error) {
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

// This will forcably kill any process as well as stopping any starting process
func (i *Instance) ForceStop() error {
	i.forceStop <- 1
	return i.cmd.Process.Kill()
}

func (i *Instance) ReverseProxy() (*httputil.ReverseProxy, error) {
	if !i.Started() {
		return nil, fmt.Errorf("Instance not yet started")
	}
	return i.proxy, nil
}
