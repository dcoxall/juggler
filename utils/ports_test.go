package utils

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestIsPortFree(t *testing.T) {
	port := <-FindAvailablePort()
	addr := fmt.Sprintf(":%d", port)
	if !IsPortFree(port) {
		t.Fatalf("Expected port (%d) to be free", port)
	}
	cmd := exec.Command("/tmp/ping", addr, "pong")
	cmd.Start()
	go func() { cmd.Wait() }()
	<-time.After(time.Second)
	if IsPortFree(port) {
		t.Fatalf("Expected port (%d) to not be free", port)
	}
	cmd.Process.Kill()
}
