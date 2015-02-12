package utils

import (
	"fmt"
	"os"
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
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()
	go func() { cmd.Wait() }()
	<-time.After(2 * time.Second)
	if IsPortFree(port) {
		t.Fatalf("Expected port (%d) to not be free", port)
	}
	cmd.Process.Kill()
}
