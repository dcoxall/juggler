package utils

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestIsPortFree(t *testing.T) {
	port := 3456
	addr := fmt.Sprintf("localhost:%d", port)
	if !IsPortFree(port) {
		t.Fatalf("Expected port (%d) to be free", port)
	}
	cmd := exec.Command("ping", addr, "pong")
	go func() { cmd.Run() }()
	<-time.After(2 * time.Second)
	if IsPortFree(port) {
		t.Fatalf("Expected port (%d) to not be free", port)
	}
	cmd.Process.Kill()
}
