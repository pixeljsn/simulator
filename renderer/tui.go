package renderer

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"simulator/engine"
)

const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
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
		colorize("╔════════════════════════════════════════════════════════════════════════════════╗", ansiBlue),
		fmt.Sprintf("║ %s %-67s ║", colorize("⚡ MINI K8S + GPU CHAOS-VIZ", ansiBold+ansiMagenta), colorize(fmt.Sprintf("Tick %d", r.eng.TickCount), ansiBold+ansiCyan)),
		fmt.Sprintf("║ %s ║", pad(colorize("Legend: RED=hot busy  YELLOW=warming  GREEN=about-to-release  CYAN=idle", ansiDim), 78)),
		fmt.Sprintf("║ %s ║", pad(colorize("States: RUNNABLE -> DISPATCH -> RUNNING -> (CONTEXT-SWITCH-OVERHEAD) -> RUNNABLE", ansiDim), 78)),
		fmt.Sprintf("║ %s ║", pad(colorize("This Tick ► "+r.tickSummary(), ansiBold), 78)),
		colorize("╠════════════════════════════════════════════════════════════════════════════════╣", ansiBlue),
		fmt.Sprintf("║ %s ║", pad(colorize("CPU CORE MATRIX", ansiBold+ansiYellow), 78)),
	}

	for _, c := range r.eng.CPU.Cores {
		state := colorize("IDLE", ansiCyan)
		meter := colorize("░░░░░░░░░░", ansiCyan)
		if c.ContextSwitchTax > 0 {
			state = colorize(fmt.Sprintf("CONTEXT-SWITCH-OVERHEAD (%d tick left)", c.ContextSwitchTax), ansiYellow)
			meter = colorize("▓▓░░░░░░░░", ansiYellow)
		} else if c.Running != nil {
			state = colorize(fmt.Sprintf("RUNNING %-12s (%d ticks left)", c.Running.Name, c.Running.RemainingTicks), cpuStateColor(c.Running.RemainingTicks))
			meter = heatBar(c.Running.RemainingTicks, 12)
		}
		row := fmt.Sprintf("CPU%-2d %s  %s", c.ID, meter, state)
		lines = append(lines, fmt.Sprintf("║ %s ║", pad(row, 78)))
	}

	lines = append(lines,
		colorize("╠════════════════════════════════════════════════════════════════════════════════╣", ansiBlue),
		fmt.Sprintf("║ %s ║", pad(colorize("GPU CHAMBER", ansiBold+ansiMagenta), 78)),
	)

	gpuState := colorize("IDLE", ansiCyan)
	gpuMeter := colorize("░░░░░░░░░░░░", ansiCyan)
	if r.eng.GPU.RunningJob != nil {
		rem := r.eng.GPU.RunningJob.RemainingTicks
		gpuState = colorize(fmt.Sprintf("LOCKED by %-12s (%d ticks left)", r.eng.GPU.RunningJob.JobID, rem), gpuStateColor(rem))
		gpuMeter = heatBar(rem, 18)
	}
	lines = append(lines,
		fmt.Sprintf("║ %s ║", pad(fmt.Sprintf("GPU    %s  %s", gpuMeter, gpuState), 78)),
		fmt.Sprintf("║ %s ║", pad(fmt.Sprintf("MEMORY %s", memoryBar(r.eng.GPU.UsedMemory, r.eng.GPU.TotalMemory)), 78)),
	)

	queue := make([]string, 0)
	for _, p := range r.eng.CPU.ReadyQueue.Snapshot() {
		queue = append(queue, p.Name)
	}
	queueView := strings.Join(queue, " -> ")
	if queueView == "" {
		queueView = "(empty)"
	}
	lines = append(lines,
		colorize("╠════════════════════════════════════════════════════════════════════════════════╣", ansiBlue),
		fmt.Sprintf("║ %s ║", pad(colorize("DISPATCH PIPELINE", ansiBold+ansiGreen), 78)),
		fmt.Sprintf("║ %s ║", pad("RUNNABLE QUEUE ► "+queueView, 78)),
	)

	lines = append(lines,
		colorize("╠════════════════════════════════════════════════════════════════════════════════╣", ansiBlue),
		fmt.Sprintf("║ %s ║", pad(colorize("CLUSTER NODE RADAR", ansiBold+ansiCyan), 78)),
	)
	for _, n := range r.eng.Nodes {
		nodeRow := fmt.Sprintf("%s  CPU %s  MEM %s  GPU %s", n.Name, ratioBar(n.CPUAllocated, n.CPUCapacity, 12), ratioBar(n.MemAllocated, n.MemCapacity, 12), ratioBar(n.GPUAllocated, max(1, n.GPUCapacity), 6))
		lines = append(lines, fmt.Sprintf("║ %s ║", pad(nodeRow, 78)))
	}

	lines = append(lines,
		colorize("╠════════════════════════════════════════════════════════════════════════════════╣", ansiBlue),
		fmt.Sprintf("║ %s ║", pad(colorize("EVENT STORM (newest first)", ansiBold+ansiRed), 78)),
	)
	for i, e := range r.eng.Events.Entries() {
		if i == 6 {
			break
		}
		lines = append(lines, fmt.Sprintf("║ %s ║", pad("• "+e, 78)))
	}

	lines = append(lines, colorize("╚════════════════════════════════════════════════════════════════════════════════╝", ansiBlue))
	return strings.Join(lines, "\n")
}

func pad(s string, width int) string {
	plain := stripANSI(s)
	if len(plain) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(plain))
}

func stripANSI(s string) string {
	res := strings.Builder{}
	esc := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == 0x1b {
			esc = true
			continue
		}
		if esc {
			if ch == 'm' {
				esc = false
			}
			continue
		}
		res.WriteByte(ch)
	}
	return res.String()
}

func colorize(s, color string) string { return color + s + ansiReset }

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

func heatBar(remaining, scale int) string {
	if scale <= 0 {
		scale = 1
	}
	filled := remaining
	if filled > scale {
		filled = scale
	}
	if filled < 0 {
		filled = 0
	}
	segments := 12
	n := filled * segments / scale
	if n < 1 && remaining > 0 {
		n = 1
	}
	bar := strings.Repeat("█", n) + strings.Repeat("░", segments-n)
	return colorize(bar, cpuStateColor(remaining))
}

func memoryBar(used, total int) string {
	if total <= 0 {
		total = 1
	}
	segments := 18
	n := used * segments / total
	if n > segments {
		n = segments
	}
	bar := strings.Repeat("■", n) + strings.Repeat("·", segments-n)
	color := ansiGreen
	pct := used * 100 / total
	if pct >= 80 {
		color = ansiRed
	} else if pct >= 50 {
		color = ansiYellow
	}
	return fmt.Sprintf("%s %d/%d GB", colorize(bar, color), used, total)
}

func ratioBar(used, total, segments int) string {
	if total <= 0 {
		return colorize("n/a", ansiDim)
	}
	n := used * segments / total
	if n > segments {
		n = segments
	}
	bar := strings.Repeat("▮", n) + strings.Repeat("▯", segments-n)
	pct := used * 100 / total
	color := ansiGreen
	if pct >= 80 {
		color = ansiRed
	} else if pct >= 50 {
		color = ansiYellow
	}
	return colorize(fmt.Sprintf("%s %d/%d", bar, used, total), color)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
