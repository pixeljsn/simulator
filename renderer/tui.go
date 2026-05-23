package renderer

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"simulator/engine"
)

type Renderer struct {
	eng     *engine.Engine
	stopped atomic.Bool
}

func New(e *engine.Engine) *Renderer { return &Renderer{eng: e} }

func (r *Renderer) Stop() { r.stopped.Store(true) }

func (r *Renderer) RunLoop() {
	ticker := time.NewTicker(r.eng.Config.TickDuration)
	defer ticker.Stop()
	for !r.stopped.Load() {
		r.eng.Tick()
		fmt.Print("\033[H\033[2J")
		fmt.Println(r.View())
		if r.eng.Config.MaxTicks > 0 && r.eng.TickCount >= r.eng.Config.MaxTicks {
			r.Stop()
		}
		<-ticker.C
	}
}

func (r *Renderer) View() string {
	lines := []string{
		"+--------------------------------------------------------------------------------+",
		fmt.Sprintf("| Tick %-74d|", r.eng.TickCount),
		"| Legend: tick = one simulation step; 'ticks left' = remaining work units.       |",
		"|                                                                                |",
	}

	for _, c := range r.eng.CPU.Cores {
		state := "IDLE"
		if c.ContextSwitchTax > 0 {
			state = fmt.Sprintf("IDLE (context switch overhead: %d tick left)", c.ContextSwitchTax)
		} else if c.Running != nil {
			state = fmt.Sprintf("RUNNING %-12s (%d ticks left in process)", c.Running.Name, c.Running.RemainingTicks)
		}
		lines = append(lines, fmt.Sprintf("| CPU%-2d | %-72s|", c.ID, state))
	}

	gpu := "IDLE"
	if r.eng.GPU.RunningJob != nil {
		gpu = fmt.Sprintf("LOCKED by %-12s (%d ticks left in GPU job)", r.eng.GPU.RunningJob.JobID, r.eng.GPU.RunningJob.RemainingTicks)
	}

	lines = append(lines,
		fmt.Sprintf("| GPU  | %-72s|", gpu),
		fmt.Sprintf("| GPU Memory: %-65s|", fmt.Sprintf("%d/%d GB reserved", r.eng.GPU.UsedMemory, r.eng.GPU.TotalMemory)),
		"|                                                                                |",
	)

	queue := make([]string, 0)
	for _, p := range r.eng.CPU.ReadyQueue.Snapshot() {
		queue = append(queue, p.Name)
	}
	queueView := strings.Join(queue, " -> ")
	if queueView == "" {
		queueView = "(empty)"
	}
	lines = append(lines, fmt.Sprintf("| Ready Queue (next dispatch order): %-43s|", queueView))
	lines = append(lines, "|                                                                                |")

	for _, n := range r.eng.Nodes {
		lines = append(lines, fmt.Sprintf("| Node %-72s|", fmt.Sprintf("%s CPU:%d/%d MEM:%d/%d GPU:%d/%d", n.Name, n.CPUAllocated, n.CPUCapacity, n.MemAllocated, n.MemCapacity, n.GPUAllocated, n.GPUCapacity)))
	}

	lines = append(lines, "|                                                                                |", "| Recent Events (newest first):                                                 |")
	for i, e := range r.eng.Events.Entries() {
		if i == 6 {
			break
		}
		lines = append(lines, fmt.Sprintf("| - %-78s|", e))
	}
	lines = append(lines, "+--------------------------------------------------------------------------------+")
	return strings.Join(lines, "\n")
}
