package cmd

import (
	"fmt"
	"strings"
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

	if !flagStats {
		return nil
	}

	pct, err := proc.CPUPercent(pid, 500*time.Millisecond)
	if err != nil {
		return fmt.Errorf("cannot get CPU usage for PID %d: %w", pid, err)
	}
	fmt.Printf("CPU%%:     %.1f%%\n", pct)

	if st, err := proc.ReadStatus(pid); err == nil {
		fmt.Printf("Vol ctx switches:     %d\n", st.VoluntaryCtxtSwitches)
		fmt.Printf("Nonvol ctx switches:  %d\n", st.NonVoluntaryCtxtSwitches)
	}

	if pressure, err := proc.CPUPressure(pid); err == nil && pressure != "" {
		fmt.Println("CPU pressure:")
		for _, line := range strings.Split(pressure, "\n") {
			fmt.Printf("  %s\n", line)
		}
	}

	return nil
}
