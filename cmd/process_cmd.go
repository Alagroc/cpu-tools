package cmd

import (
	"fmt"
	"time"

	"cpu-tools/internal/proc"
)

func runProcess() error {
	pid := flagProcess

	name, _ := proc.GetName(pid)
	if name == "" {
		name = "-"
	}

	cpus, err := proc.GetAffinity(pid)
	if err != nil {
		return fmt.Errorf("cannot get affinity for PID %d: %w", pid, err)
	}

	fmt.Printf("PID:      %d\n", pid)
	fmt.Printf("Name:     %s\n", name)
	fmt.Printf("Affinity: %v\n", cpus)

	if flagShow {
		pct, err := proc.CPUPercent(pid, 500*time.Millisecond)
		if err != nil {
			return fmt.Errorf("cannot get CPU usage for PID %d: %w", pid, err)
		}
		fmt.Printf("CPU%%:     %.1f%%\n", pct)
	}
	return nil
}
