package hexa

type ChanToStepControllableOne struct {
	trigger chan interface{}
	done    chan interface{}
}

func NewStepControllerFromChan(trigger chan interface{}, done chan interface{}) *ChanToStepControllableOne {
	return &ChanToStepControllableOne{
		trigger: trigger,
		done:    done,
	}
}

func (__ *ChanToStepControllableOne) StepTrigger() chan<- interface{} {
	return __.trigger
}

func (__ *ChanToStepControllableOne) StepDoneNotify() <-chan interface{} {
	return __.done
}
