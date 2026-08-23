package cmd

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"cpu-tools/internal/proc"
)

func runActive() error {
	numCPU := runtime.NumCPU()

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
		disp string
	}

	var matching []entry
	for _, p := range procs {
		cpus, err := proc.GetAffinity(p.PID)
		if err != nil {
			continue
		}
		if coreSet != nil {
			if !flagShowAllAffinities && len(cpus) == numCPU {
				continue
			}
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
		disp := p.Name
		if flagPidResolution {
			if cmdline, err := proc.GetCmdline(p.PID); err == nil && cmdline != "" {
				disp = truncate(cmdline, 40)
			}
		}
		matching = append(matching, entry{p: p, cpus: cpus, disp: disp})
	}

	nameW, nameHdr := 20, "NAME"
	if flagPidResolution {
		nameW, nameHdr = 40, "CMDLINE"
	}

	if !flagShow {
		fmt.Printf("%-8s %-*s %s\n", "PID", nameW, nameHdr, "CPUS")
		for _, e := range matching {
			fmt.Printf("%-8d %-*s %v\n", e.p.PID, nameW, e.disp, e.cpus)
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

	fmt.Printf("%-8s %-*s %-8s %s\n", "PID", nameW, nameHdr, "CPU%", "CPUS")
	for _, e := range matching {
		fmt.Printf("%-8d %-*s %-8.1f %v\n", e.p.PID, nameW, e.disp, pcts[e.p.PID], e.cpus)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
