package server

type SubscriptionManager struct {
	Command <-chan TelemetryCommand
}

func NewSubscriptionManager() SubscriptionManager {
	return SubscriptionManager{make(chan TelemetryCommand)}
}

func (sm *SubscriptionManager) Listen() {

	var tc TelemetryCommand

	for {
		tc = <-sm.Command

		switch tc.Cmd {
		case Subscribe:

		case Unsubscribe:
		}
	}
}

type pubsub interface {
	update(key string)
}
