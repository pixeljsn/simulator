# Mini Kubernetes + GPU Infrastructure Simulator

A fully local, deterministic, visualization-first simulator written in Go.

This project simulates CPU scheduling, GPU monopolization behavior, cgroup-like throttling signals, and Kubernetes-style node/pod state for educational purposes.

## Requirements

- Go 1.22+
- Terminal with ANSI escape support

## Build

```bash
go build -o simulator .
```

## Run

### Default run

```bash
go run .
```

Defaults:
- tick duration: `300ms`
- max ticks: `120`
- interactive mode: `false` (auto-exits at `max-ticks`)

### Slower run (recommended for learning)

```bash
go run . -tick=2s -max-ticks=40
```

This makes each tick stay on screen for 2 seconds so state changes are easier to follow.

### Custom run options

```bash
go run . -tick=200ms -max-ticks=200
```

Flags:
- `-tick duration` tick duration (example: `100ms`, `1s`, `2s`)
- `-max-ticks int` number of ticks to run (`0` means run until manually stopped)
- `-interactive` if true, press ENTER to stop

Interactive example:

```bash
go run . -interactive -max-ticks=0 -tick=1s
```


### Reading the screen quickly

- **Tick N**: the simulator has advanced N steps.
- **CPU row**:
  - `RUNNING <name> (X ticks left in process)` means that process needs X more ticks on CPU.
  - `IDLE (context switch overhead: 1 tick left)` means the core is intentionally paused for switch cost, not necessarily out of work.
- **GPU row**:
  - `LOCKED by <job> (X ticks left in GPU job)` means GPU is monopolized by one non-preemptive job.
- **Ready Queue** shows next-dispatch order from left to right.
- **Recent Events** explains *why* visible changes happened.

## Understanding ticks

Think of a **tick** as one simulation "frame" or "clock step".

On each tick, the engine does three things in order:
1. CPU scheduler advances process execution by one step.
2. GPU scheduler advances the running GPU job by one step.
3. New events are added to the event log and rendered.

So if `-tick=1s`, then **every 1 second** the simulator advances by exactly one scheduling step.

### How to read `RemainingTicks`

When you see a process or GPU job like `(8t)`, it means:
- it needs 8 more simulation ticks to complete
- after each tick while running, this number decreases
- when it reaches `0`, the process/job completes

### Why it can feel hard to follow

With the default `300ms` tick duration and screen-clearing render loop, updates are very fast.
Use a slower value (for example `1s` or `2s`) to observe transitions clearly.

## Test

```bash
go test ./...
```

## Notes

- This is a conceptual simulator only.
- It does **not** use real containers, real Kubernetes APIs, or real GPUs.
