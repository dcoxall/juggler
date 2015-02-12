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

type MockInstance struct {
	port      int
	ref       string
	state     int
	stateChan chan int
	forceStop chan int
	cmd       *exec.Cmd
	proxy     *httputil.ReverseProxy
}

func NewInstance(port int, ref string) *MockInstance {
	return &MockInstance{
		port:      port,
		ref:       ref,
		state:     Stopped,
		stateChan: make(chan int, 1),
		forceStop: make(chan int, 1),
	}
}

func (i *MockInstance) Start() (<-chan int, error) {
	// We can't start if we aren't stopped
	if i.state != Stopped {
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
	go func(instance *MockInstance) {
		for instance.state != Running {
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

func (i *MockInstance) Stop() (<-chan int, error) {
	// We can't start if we aren't stopped
	if i.state != Running {
		return i.stateChan, fmt.Errorf("Unable to stop")
	}
	i.state = Stopping

	// Begin watching this process and signal when it ends
	go func(instance *MockInstance) {
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
func (i *MockInstance) ForceStop() error {
	i.forceStop <- 1
	return i.cmd.Process.Kill()
}

func (i *MockInstance) ReverseProxy() (*httputil.ReverseProxy, error) {
	if i.state != Running {
		return nil, fmt.Errorf("Instance not yet started")
	}
	return i.proxy, nil
}
