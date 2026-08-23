package cmd

import (
	"fmt"
	"sort"
	"time"

	"cpu-tools/internal/proc"
)

func runActive() error {
	var coreSet map[int]bool
	if flagCores != "" {
		cores, err := parseCores(flagCores)
		if err != nil {
			return err
		}
		coreSet = make(map[int]bool, len(cores))
		for _, c := range cores {
			coreSet[c] = true
		}
	}

	procs, err := proc.AllProcesses()
	if err != nil {
		return err
	}

	type entry struct {
		p    proc.Process
		cpus []int
	}

	var matching []entry
	for _, p := range procs {
		cpus, err := proc.GetAffinity(p.PID)
		if err != nil {
			continue
		}
		if coreSet != nil {
			var found bool
			for _, c := range cpus {
				if coreSet[c] {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		matching = append(matching, entry{p: p, cpus: cpus})
	}

	if !flagShow {
		fmt.Printf("%-8s %-20s %s\n", "PID", "NAME", "CPUS")
		for _, e := range matching {
			fmt.Printf("%-8d %-20s %v\n", e.p.PID, e.p.Name, e.cpus)
		}
		return nil
	}

	pids := make([]int, len(matching))
	for i, e := range matching {
		pids[i] = e.p.PID
	}

	pcts, err := proc.CPUPercentBatch(pids, 500*time.Millisecond)
	if err != nil {
		return err
	}

	sort.Slice(matching, func(i, j int) bool {
		return pcts[matching[i].p.PID] > pcts[matching[j].p.PID]
	})

	fmt.Printf("%-8s %-20s %-8s %s\n", "PID", "NAME", "CPU%", "CPUS")
	for _, e := range matching {
		fmt.Printf("%-8d %-20s %-8.1f %v\n", e.p.PID, e.p.Name, pcts[e.p.PID], e.cpus)
	}
	return nil
}
