package gpu

import "simulator/models"

type Scheduler struct {
	TotalMemory int
	UsedMemory  int
	RunningJob  *models.GPUJob
	Queue       []*models.GPUJob
}

func NewScheduler(totalMemory int) *Scheduler {
	return &Scheduler{TotalMemory: totalMemory, Queue: make([]*models.GPUJob, 0)}
}

func (s *Scheduler) Enqueue(job *models.GPUJob) {
	s.Queue = append(s.Queue, job)
}

func (s *Scheduler) Tick(onEvent func(string)) {
	if s.RunningJob == nil {
		if len(s.Queue) > 0 {
			n := s.Queue[0]
			s.Queue = s.Queue[1:]
			s.RunningJob = n
			s.UsedMemory += n.MemoryRequested
			onEvent("GPU locked by job " + n.JobID)
		}
		return
	}

	s.RunningJob.RemainingTicks--
	if s.RunningJob.RemainingTicks <= 0 {
		onEvent("GPU job finished: " + s.RunningJob.JobID)
		s.UsedMemory -= s.RunningJob.MemoryRequested
		s.RunningJob = nil
	}
}
