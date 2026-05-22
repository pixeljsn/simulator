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
	lines := []string{"+--------------------------------------------------+", fmt.Sprintf("| Tick %-45d|", r.eng.TickCount), "|                                                  |"}
	for _, c := range r.eng.CPU.Cores {
		state := "IDLE"
		if c.Running != nil {
			state = fmt.Sprintf("%s (%dt)", c.Running.Name, c.Running.RemainingTicks)
		}
		lines = append(lines, fmt.Sprintf("| CPU%-2d | %-42s|", c.ID, state))
	}
	gpu := "IDLE"
	if r.eng.GPU.RunningJob != nil {
		gpu = fmt.Sprintf("%s (LOCKED %dt)", r.eng.GPU.RunningJob.JobID, r.eng.GPU.RunningJob.RemainingTicks)
	}
	lines = append(lines, fmt.Sprintf("| GPU  | %-42s|", gpu), fmt.Sprintf("| GPU MEM: %-42s|", fmt.Sprintf("%d/%d GB", r.eng.GPU.UsedMemory, r.eng.GPU.TotalMemory)), "|                                                  |")
	queue := make([]string, 0)
	for _, p := range r.eng.CPU.ReadyQueue.Snapshot() { queue = append(queue, p.Name) }
	lines = append(lines, fmt.Sprintf("| Ready Queue: %-37s|", strings.Join(queue, " ")), "|                                                  |")
	for _, n := range r.eng.Nodes { lines = append(lines, fmt.Sprintf("| %-48s|", fmt.Sprintf("%s CPU:%d/%d MEM:%d/%d GPU:%d/%d", n.Name, n.CPUAllocated, n.CPUCapacity, n.MemAllocated, n.MemCapacity, n.GPUAllocated, n.GPUCapacity))) }
	lines = append(lines, "|                                                  |", "| Events:                                           |")
	for i, e := range r.eng.Events.Entries() {
		if i == 6 { break }
		lines = append(lines, fmt.Sprintf("| - %-46s|", e))
	}
	lines = append(lines, "+--------------------------------------------------+")
	return strings.Join(lines, "\n")
}
