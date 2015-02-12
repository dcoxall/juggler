package juggler

import "net/http/httputil"

type Instancer interface {
	Started() bool
	Starting() bool
	Stopped() bool
	Stopping() bool

	Start() (<-chan InstanceState, error)
	Stop() (<-chan InstanceState, error)
	ForceStop() error

	ReverseProxy() (*httputil.ReverseProxy, error)
}
