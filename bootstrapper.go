package juggler

import "os/exec"

// Bootstrapper is the interface used to plug-in custom applications to the
// Instance structure making it easier to use existing logic for communicating
// instance state.
type Bootstrapper interface {
	// Bootstrap takes an int representing the port and a string representing
	// the reference to the Instance and expects the command attatched to the
	// instance process to be returned or an error if it can't be started
	Bootstrap(int, string) (*exec.Cmd, error)
}
