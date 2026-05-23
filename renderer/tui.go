package renderer

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"simulator/engine"
)

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
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
		"| Status terms: RUNNING, RUNNABLE, DISPATCH, CONTEXT-SWITCH-OVERHEAD, LOCKED.   |",
		fmt.Sprintf("| This Tick Summary: %-60s|", r.tickSummary()),
		"|                                                                                |",
	}

	for _, c := range r.eng.CPU.Cores {
		state := "IDLE"
		if c.ContextSwitchTax > 0 {
			state = colorize(fmt.Sprintf("CONTEXT-SWITCH-OVERHEAD (%d tick left)", c.ContextSwitchTax), ansiYellow)
		} else if c.Running != nil {
			state = colorize(fmt.Sprintf("RUNNING %-12s (%d ticks left)", c.Running.Name, c.Running.RemainingTicks), cpuStateColor(c.Running.RemainingTicks))
		} else {
			state = colorize(state, ansiCyan)
		}
		lines = append(lines, fmt.Sprintf("| CPU%-2d | %-72s|", c.ID, state))
	}

	gpu := "IDLE"
	if r.eng.GPU.RunningJob != nil {
		gpu = colorize(fmt.Sprintf("LOCKED by %-12s (%d ticks left)", r.eng.GPU.RunningJob.JobID, r.eng.GPU.RunningJob.RemainingTicks), gpuStateColor(r.eng.GPU.RunningJob.RemainingTicks))
	} else {
		gpu = colorize(gpu, ansiCyan)
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
	lines = append(lines, fmt.Sprintf("| Ready Queue (RUNNABLE; next DISPATCH order): %-36s|", queueView))
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

func colorize(s, color string) string {
	return color + s + ansiReset
}

func cpuStateColor(remaining int) string {
	if remaining <= 2 {
		return ansiGreen
	}
	if remaining <= 5 {
		return ansiYellow
	}
	return ansiRed
}

func gpuStateColor(remaining int) string {
	if remaining <= 2 {
		return ansiGreen
	}
	if remaining <= 6 {
		return ansiYellow
	}
	return ansiRed
}

func (r *Renderer) tickSummary() string {
	events := r.eng.Events.Entries()
	if len(events) == 0 {
		return "No events recorded yet"
	}
	prefix := fmt.Sprintf("[T%03d]", r.eng.TickCount)
	parts := make([]string, 0, 3)
	for _, e := range events {
		if strings.HasPrefix(e, prefix) {
			parts = append(parts, strings.TrimSpace(strings.TrimPrefix(e, prefix)))
			if len(parts) == 3 {
				break
			}
		}
	}
	if len(parts) == 0 {
		return "No state transition event this tick"
	}
	return strings.Join(parts, " | ")
}
