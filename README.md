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
- tick duration (`cycle duration`): `300ms`
- max ticks (`max cycles`): `120`
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
- `-tick duration` tick duration / cycle duration (example: `100ms`, `1s`, `2s`)
- `-max-ticks int` number of ticks to run (`0` means run until manually stopped)
- `-interactive` if true, press ENTER to stop

Interactive example:

```bash
go run . -interactive -max-ticks=0 -tick=1s
```


### Reading the screen quickly

- **Cycle N** (formerly shown as tick): the simulator has advanced N simulation steps.
- **CPU row**:
  - `RUNNING <name> (X ticks left in process)` means that process needs X more ticks on CPU.
  - `IDLE (context switch overhead: 1 tick left)` means the core is intentionally paused for switch cost, not necessarily out of work.
- **GPU row**:
  - `LOCKED by <job> (X ticks left in GPU job)` means GPU is monopolized by one non-preemptive job.
- **Ready Queue** shows next-dispatch order from left to right.
- **Recent Events** explains *why* visible changes happened.
- **Color cues** (ANSI terminals):
  - **Red**: heavily engaged (many cycles left before release).
  - **Yellow**: in-progress, mid-flight.
  - **Green**: about to release (few cycles left).
  - **Cyan**: idle.
- **This Cycle Summary** is a one-line real-time explanation of what changed on the current cycle.
- **Status terms** are intentionally technical and map to scheduler behavior:
  - `RUNNABLE`: waiting in ready queue, eligible for CPU.
  - `DISPATCH`: selected from ready queue to run on a core.
  - `RUNNING`: currently consuming CPU/GPU ticks.
  - `CONTEXT-SWITCH-OVERHEAD`: temporary core stall during switch penalty.
  - `LOCKED`: GPU monopolized by one non-preemptive job.

## Understanding cycles

Think of a **cycle** as one simulation "frame" or "clock step".

On each cycle, the engine does three things in order:
1. CPU scheduler advances process execution by one step.
2. GPU scheduler advances the running GPU job by one step.
3. New events are added to the event log and rendered.

So if `-tick=1s`, then **every 1 second** the simulator advances by exactly one scheduling cycle.

### How to read `RemainingTicks` (cycles left)

When you see a process or GPU job like `(8t)`, it means:
- it needs 8 more simulation ticks to complete
- after each tick while running, this number decreases
- when it reaches `0`, the process/job completes

### Why it can feel hard to follow

With the default `300ms` cycle duration and screen-clearing render loop, updates are very fast.
Use a slower value (for example `1s` or `2s`) to observe transitions clearly.

## Events explained

The **Recent Events** panel is the simulator's reasoning trail.

- Events are shown **newest first**.
- Each event includes a tick stamp like `[T003]`, meaning it happened during simulation cycle 3.
- Use events to understand *why* a visible state changed between frames.

Common event meanings in the current demo:

- `CPU core scheduled <process>`: a core pulled a process from the ready queue and started running it.
- `Context switch on core`: a process used up its time slice and the core incurred context-switch overhead.
- `Process completed: <process>`: process reached zero remaining ticks and exited.
- `GPU locked by job <job>`: GPU accepted a job and became non-preemptively occupied.
- `GPU job finished: <job>`: running GPU job completed and GPU memory was released.
- `CPU throttled by cgroup for process <process>`: simulated cgroup throttle signal triggered for that running process.

### How to use events while watching

1. See a change in CPU/GPU rows.
2. Check the top-most event line for the cause on that cycle.
3. Correlate with queue order and cycles-left values to predict what should happen next.

## Test

```bash
go test ./...
```

## Notes

- This is a conceptual simulator only.
- It does **not** use real containers, real Kubernetes APIs, or real GPUs.

## Crazy visualization mode

The default renderer now uses a high-contrast "crazy" dashboard style:

- Neon-like ANSI colors and boxed panels
- CPU/GPU heat bars that shrink as jobs approach completion
- GPU memory pressure bar
- Node capacity radar bars for CPU/MEM/GPU
- Cycle-level summary plus event storm feed

If your terminal does not support ANSI colors, output still works but appears less vivid.

Recommended command for the best effect:

```bash
go run . -tick=500ms -max-ticks=120
```

