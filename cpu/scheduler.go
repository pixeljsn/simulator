package cpu

import (
	"simulator/models"
	"simulator/process"
)

type Core struct {
	ID               int
	Running          *models.Process
	SliceUsed        int
	ContextSwitchTax int
	SwitchPenalty    int
}

type Scheduler struct {
	Cores      []Core
	TimeSlice  int
	ReadyQueue *process.Queue
}

func NewScheduler(coreCount, timeSlice, switchPenalty int) *Scheduler {
	cores := make([]Core, coreCount)
	for i := range cores {
		cores[i] = Core{ID: i + 1, SwitchPenalty: switchPenalty}
	}
	return &Scheduler{Cores: cores, TimeSlice: timeSlice, ReadyQueue: &process.Queue{}}
}

func (s *Scheduler) Tick(onEvent func(string)) {
	for i := range s.Cores {
		c := &s.Cores[i]
		if c.ContextSwitchTax > 0 {
			c.ContextSwitchTax--
			continue
		}

		if c.Running == nil {
			next := s.ReadyQueue.Dequeue()
			if next != nil {
				next.State = models.ProcessStateRunning
				c.Running = next
				c.SliceUsed = 0
				onEvent("CPU core scheduled " + next.Name)
			}
			continue
		}

		c.Running.RemainingTicks--
		c.SliceUsed++
		if c.Running.RemainingTicks <= 0 {
			onEvent("Process completed: " + c.Running.Name)
			c.Running.State = models.ProcessStateCompleted
			c.Running = nil
			c.SliceUsed = 0
			continue
		}

		if c.SliceUsed >= s.TimeSlice {
			old := c.Running
			old.State = models.ProcessStateReady
			s.ReadyQueue.Enqueue(old)
			c.Running = nil
			c.SliceUsed = 0
			c.ContextSwitchTax = c.SwitchPenalty
			onEvent("Context switch on core")
		}
	}
}
