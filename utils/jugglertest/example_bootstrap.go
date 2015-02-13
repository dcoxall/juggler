package jugglertest

import (
	"fmt"
	"os/exec"
)

type ExampleBootstrap struct{}

func (*ExampleBootstrap) Bootstrap(port int, reference string) (cmd *exec.Cmd, err error) {
	cmd = exec.Command(
		"/tmp/ping",
		fmt.Sprintf(":%d", port),
		reference,
	)
	err = cmd.Start()
	return
}
