package device

type Commands int

const (
	// PING sends a message to the device telling it to immedietly respond on a response topic
	// in order to determine at time of request if the device is reachable and responsive
	PING Commands = iota
	OTA
)

// BrokerCommands are the commands that the Device can send to the broker
type BrokerCommands int

const (
	BIND BrokerCommands = iota
	UNBIND
)
