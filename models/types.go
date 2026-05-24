package models

type ProcessState string

const (
	ProcessStateReady     ProcessState = "Ready"
	ProcessStateRunning   ProcessState = "Running"
	ProcessStateWaiting   ProcessState = "Waiting"
	ProcessStateBlocked   ProcessState = "Blocked"
	ProcessStateCompleted ProcessState = "Completed"
)

type Process struct {
	PID            int
	Name           string
	State          ProcessState
	RemainingTicks int
	CPURequest     int
	MemoryRequest  int
	Namespace      string
	ContainerID    string
	CgroupID       string
}

type GPUJob struct {
	JobID           string
	ProcessPID      int
	RemainingTicks  int
	MemoryRequested int
}

type Pod struct {
	Name      string
	Namespace string
	CPU       int
	Memory    int
	GPU       int
	Status    string
	NodeName  string
	Processes []Process
}

type Node struct {
	Name         string
	CPUCapacity  int
	MemCapacity  int
	GPUCapacity  int
	CPUAllocated int
	MemAllocated int
	GPUAllocated int
}
