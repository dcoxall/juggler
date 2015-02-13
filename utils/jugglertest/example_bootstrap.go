package jugglertest

import (
	"fmt"
	"os/exec"
)

// ExampleBootstrap is an example of a basic Bootstrapper structure that uses a
// test webserver. It's primary use is to aid in the testing of juggler.
type ExampleBootstrap struct{}

// Bootstrap will start an example webserver on the provided port.
func (*ExampleBootstrap) Bootstrap(port int, reference string) (cmd *exec.Cmd, err error) {
	cmd = exec.Command(
		"/tmp/ping",
		fmt.Sprintf(":%d", port),
		reference,
	)
	err = cmd.Start()
	return
}
