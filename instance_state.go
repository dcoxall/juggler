package juggler

type InstanceState int

const (
	Stopped InstanceState = iota
	Starting
	Running
	Stopping
)
