package juggler

import "os/exec"

type Bootstrapper interface {
	Bootstrap(int, string) (*exec.Cmd, error)
}
