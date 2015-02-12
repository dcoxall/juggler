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
	port        int
	ref         string
	state       int
	startChan   chan int
	stopChan    chan int
	forceStop   chan int
	cmd         *exec.Cmd
	proxy       *httputil.ReverseProxy
}

func NewInstance(port int, ref string) *MockInstance {
	return &MockInstance{
		port:        port,
		ref:         ref,
		state:       Stopped,
		startChan:   make(chan int, 1),
		stopChan:    make(chan int, 1),
		forceStop:   make(chan int, 1),
	}
}

func (i *MockInstance) Start() (<-chan int, error) {
	// We can't start if we aren't stopped
	if i.state != Stopped {
		return i.startChan, fmt.Errorf("Unable to start")
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
		return i.startChan, err
	}
	i.state = Starting

	// Begin watching this process and signal when it ends
	go func(){
		if err := i.cmd.Wait(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf("stopped\n")
		i.state = Stopped
		i.stopChan <- i.state
	}()

	// in the background let's wait until we can connect and then trigger
	// completion on the channel
	go func() {
		for {
			select {
			case <- i.forceStop:
				return
			default:
				if utils.IsPortFree(i.port) {
					time.Sleep(500 * time.Millisecond)
				} else {
					break
				}
			}
		}
		i.state = Running
		url := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", i.port),
		}
		i.proxy = httputil.NewSingleHostReverseProxy(url)
		i.startChan <- i.state
	}()

	// return our channel that indicates a change in state
	return i.startChan, nil
}

func (i *MockInstance) Stop() (<-chan int, error) {
	// We can't start if we aren't stopped
	if i.state != Running {
		return i.stopChan, fmt.Errorf("Unable to stop")
	}
	i.state = Stopping

	// we already have a go routine watching the command and that will update
	// the state when it ends so we just need to trigger a change in that
	// process to proceed with stopping the process.
	if err := i.cmd.Process.Signal(os.Kill); err != nil {
		return i.stopChan, err
	}

	// return our channel that indicates a change in state
	return i.stopChan, nil
}

// This will forcably kill any process as well as stopping any starting process
func (i *MockInstance) ForceStop() (error) {
	i.forceStop <- 1
	return i.cmd.Process.Kill()
}

func (i *MockInstance) ReverseProxy() (*httputil.ReverseProxy, error) {
	if i.state != Running {
		return nil, fmt.Errorf("Instance not yet started")
	}
	return i.proxy, nil
}
