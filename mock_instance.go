package juggler

import (
	"fmt"
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
	stateChange chan int
	cmd         *exec.Cmd
}

func NewInstance(port int, ref string) *MockInstance {
	return &MockInstance{
		port:        port,
		ref:         ref,
		state:       Stopped,
		stateChange: make(chan int),
	}
}

func (i *MockInstance) Start() (<-chan int, error) {
	// We can't start if we aren't stopped
	if i.state != Stopped {
		return i.stateChange, fmt.Errorf("Unable to start")
	}

	// store and start the command
	i.cmd = exec.Command(
		"ping",
		fmt.Sprintf(":%s", i.port),
		i.ref,
	)
	if err := i.cmd.Start(); err != nil {
		return i.stateChange, err
	}
	i.state = Starting

	// in the background let's wait 5 seconds and trigger completion
	time.AfterFunc(
		2*time.Second,
		func() {
			i.state = Running
			i.stateChange <- i.state
		},
	)

	// return our channel that indicates a change in state
	return i.stateChange, nil
}

func (i *MockInstance) Stop() (<-chan int, error) {
	// We can't start if we aren't stopped
	if i.state != Running {
		return i.stateChange, fmt.Errorf("Unable to stop")
	}
	i.state = Stopping

	// wait for the program to exit so we free resources
	// and indicate state once completed
	go func() {
		i.cmd.Wait()
		i.state = Stopped
		i.stateChange <- i.state
	}()

	// signal the process to stop
	if err := i.cmd.Process.Signal(os.Kill); err != nil {
		return i.stateChange, err
	}

	// return our channel that indicates a change in state
	return i.stateChange, nil
}
