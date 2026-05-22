package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"simulator/engine"
	"simulator/renderer"
)

func main() {
	tickDuration := flag.Duration("tick", 300*time.Millisecond, "simulation tick duration")
	maxTicks := flag.Int("max-ticks", 120, "maximum ticks to run (0 = unlimited)")
	interactive := flag.Bool("interactive", false, "wait for ENTER to stop")
	flag.Parse()

	eng := engine.New(engine.Config{TickDuration: *tickDuration, MaxTicks: *maxTicks})
	eng.SeedDemoScenario()

	r := renderer.New(eng)
	if *interactive {
		go r.RunLoop()
		fmt.Println("Mini Infrastructure Simulator running. Press ENTER to stop.")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
		r.Stop()
		return
	}

	r.RunLoop()
}
