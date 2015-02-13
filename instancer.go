package juggler

import "net/http/httputil"

// Instancer is the interface expected to be stored in the Spawner. It
// represents a web server listening on a specific port that can be proxied to.
//
// A good example of this is if you wish to proxy requests to multiple versions
// of the same web application a string can indicate what version to create and
// proxy to.
type Instancer interface {
	// Returns true if the Instancer has started. (i.e. can be used)
	Started() bool
	// Returns true if the Instancer is in the process of
	// starting/bootstrapping.
	Starting() bool
	// Returns true if the Instancer has stopped and indicates requests will
	// fail.
	Stopped() bool
	// Returns true if the Instancer is in the process of stopping. We should no
	// longer be sending requests to this Instancer.
	Stopping() bool

	// Start should attempt to bootstrap and start the Instancer. This should be
	// asynchronous and use the InstanceState channel to signal completion.
	Start() (<-chan InstanceState, error)
	// Stop should shutdown the Instancer and tidy up any related processes.
	// Completion is signaled by adding the juggler.Stopped InstanceState to the
	// InstaceState channel.
	Stop() (<-chan InstanceState, error)
	// If all else fails to stop the process this should kill it and anything
	// linked to it. This should be used if the Instancer times-out during
	// bootstrapping.
	ForceStop() error

	// ReverseProxy returns a proxyable structure that can be easily used to
	// make requests to the underlying web process.
	ReverseProxy() (*httputil.ReverseProxy, error)
}
