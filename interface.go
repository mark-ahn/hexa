package hexa

type StepControllableOne interface {
	StepTrigger() chan<- interface{}
	StepDoneNotify() <-chan interface{}
}

type StoppableOne interface {
	Close()
	DoneNotify() <-chan struct{}
}

type StepControlStoppableOne interface {
	StepControllableOne
	StoppableOne
}
