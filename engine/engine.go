package engine

import (
	"time"

	"simulator/cpu"
	"simulator/events"
	"simulator/gpu"
	"simulator/models"
)

type Config struct {
	TickDuration time.Duration
	MaxTicks     int
}

type Engine struct {
	Config      Config
	TickCount   int
	CPU         *cpu.Scheduler
	GPU         *gpu.Scheduler
	Events      *events.Log
	Processes   []*models.Process
	Nodes       []models.Node
	Pods        []models.Pod
	PendingPods []models.Pod
}

func New(cfg Config) *Engine {
	return &Engine{Config: cfg, CPU: cpu.NewScheduler(2, 2, 1), GPU: gpu.NewScheduler(24), Events: events.NewLog(24), Nodes: []models.Node{{Name: "node-a", CPUCapacity: 8, MemCapacity: 32000, GPUCapacity: 1}, {Name: "node-b", CPUCapacity: 8, MemCapacity: 32000, GPUCapacity: 0}}}
}

func (e *Engine) SeedDemoScenario() {
	e.Processes = []*models.Process{{PID: 1, Name: "ingest", State: models.ProcessStateReady, RemainingTicks: 8, CPURequest: 1, MemoryRequest: 512, Namespace: "ns-ai", ContainerID: "ctr-1", CgroupID: "cg-fast"}, {PID: 2, Name: "preprocess", State: models.ProcessStateReady, RemainingTicks: 12, CPURequest: 1, MemoryRequest: 1024, Namespace: "ns-ai", ContainerID: "ctr-2", CgroupID: "cg-fast"}, {PID: 3, Name: "sidecar", State: models.ProcessStateReady, RemainingTicks: 20, CPURequest: 1, MemoryRequest: 256, Namespace: "ns-ops", ContainerID: "ctr-3", CgroupID: "cg-throttle"}}
	for _, p := range e.Processes {
		e.CPU.ReadyQueue.Enqueue(p)
	}
	e.GPU.Enqueue(&models.GPUJob{JobID: "train-a", ProcessPID: 2, RemainingTicks: 18, MemoryRequested: 20})
	e.GPU.Enqueue(&models.GPUJob{JobID: "embed-b", ProcessPID: 1, RemainingTicks: 5, MemoryRequested: 8})
	e.Pods = []models.Pod{{Name: "trainer", Namespace: "ml", CPU: 2, Memory: 4096, GPU: 1, Status: "Running", NodeName: "node-a"}}
	e.PendingPods = []models.Pod{{Name: "batch-infer", Namespace: "ml", CPU: 2, Memory: 2048, GPU: 1, Status: "Pending"}}
	e.Nodes[0].CPUAllocated, e.Nodes[0].MemAllocated, e.Nodes[0].GPUAllocated = 2, 4096, 1
}

func (e *Engine) Tick() {
	e.TickCount++
	e.CPU.Tick(func(msg string) { e.Events.Add(e.TickCount, msg) })
	e.GPU.Tick(func(msg string) { e.Events.Add(e.TickCount, msg) })
	for _, p := range e.Processes {
		if p.MemoryRequest > 900 && p.State == models.ProcessStateRunning {
			e.Events.Add(e.TickCount, "CPU throttled by cgroup for process %s", p.Name)
		}
	}
}
