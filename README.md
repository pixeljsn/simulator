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

### Custom run options

```bash
go run . -tick=200ms -max-ticks=200
```

Flags:
- `-tick duration` tick duration (example: `100ms`, `1s`)
- `-max-ticks int` number of ticks to run (`0` means run until manually stopped)
- `-interactive` if true, press ENTER to stop

Interactive example:

```bash
go run . -interactive -max-ticks=0
```

## Test

```bash
go test ./...
```

## Notes

- This is a conceptual simulator only.
- It does **not** use real containers, real Kubernetes APIs, or real GPUs.
