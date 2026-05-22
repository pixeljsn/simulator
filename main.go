package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"simulator/engine"
	"simulator/renderer"
)

func main() {
	eng := engine.New(engine.Config{TickDuration: 300 * time.Millisecond, MaxTicks: 120})
	eng.SeedDemoScenario()

	r := renderer.New(eng)
	go r.RunLoop()

	fmt.Println("Mini Infrastructure Simulator running. Press ENTER to stop.")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	r.Stop()
}
