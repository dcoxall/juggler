package juggler

// InstanceState is a numerical value representing the state of an Instancer
type InstanceState int

const (
	Stopped InstanceState = iota
	Starting
	Running
	Stopping
)
