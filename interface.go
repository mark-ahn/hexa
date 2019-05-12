package hexa

// The StoppableOne interface represents a stoppable service.
// It offers the way to stop parallel routine such as go rutine.
//
// Close should try to stop the parallel routine.
//
// DoneNotify should return a channel that is closed when the parallel routine is closed.
type StoppableOne interface {
	Close()
	DoneNotify() <-chan struct{}
	Err() error
}
