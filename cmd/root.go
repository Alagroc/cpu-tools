package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	flagAll              bool
	flagCores            string
	flagProcess          int
	flagShow             bool
	flagPidResolution    bool
	flagShowAllAffinities bool
)

var rootCmd = &cobra.Command{
	Use:   "cpu-tools",
	Short: "CPU inspection utilities",
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&flagAll, "all", "a", false, "list processes on selected cores")
	rootCmd.Flags().StringVarP(&flagCores, "cores", "c", "", "core range, e.g. 0-3 or 0,1,2,3")
	rootCmd.Flags().IntVarP(&flagProcess, "process", "p", 0, "PID to inspect")
	rootCmd.Flags().BoolVarP(&flagShow, "show", "s", false, "show CPU usage (-a: percentage, -p: current usage)")
	rootCmd.Flags().BoolVar(&flagPidResolution, "pid-resolution", false, "show full cmdline instead of process name")
	rootCmd.Flags().BoolVar(&flagShowAllAffinities, "show-all-affinities", false, "include processes with affinity spanning all cores")
}

func runRoot(cmd *cobra.Command, args []string) error {
	switch {
	case flagAll:
		return runActive()
	case flagProcess > 0:
		return runProcess()
	default:
		return cmd.Help()
	}
}

func parseCores(s string) ([]int, error) {
	var result []int
	seen := make(map[int]bool)
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(strings.TrimSpace(bounds[0]))
			hi, err2 := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err1 != nil || err2 != nil || lo > hi {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := lo; i <= hi; i++ {
				if !seen[i] {
					seen[i] = true
					result = append(result, i)
				}
			}
		} else {
			n, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid core: %s", part)
			}
			if !seen[n] {
				seen[n] = true
				result = append(result, n)
			}
		}
	}
	sort.Ints(result)
	return result, nil
}
