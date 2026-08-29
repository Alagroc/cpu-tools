package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const usageCol = 34 // column at which flag usage descriptions start

func init() {
	rootCmd.SetHelpFunc(printHelp)
}

func printHelp(cmd *cobra.Command, _ []string) {
	get := func(name string) string {
		if f := cmd.Flags().Lookup(name); f != nil {
			return f.Usage
		}
		return ""
	}

	fmt.Printf("%s\n\nUsage:\n  %s [flags]\n", cmd.Short, cmd.Name())

	if cmd.Example != "" {
		fmt.Printf("\nExamples:\n%s\n", cmd.Example)
	}

	type entry struct {
		left   string
		usage  string
		sub    bool // indented sub-option of the flag above
	}

	flags := []entry{
		{"-a, --all", get("all"), false},
		{"-c, --cores string", get("cores"), false},
		{"-h, --help", "help for " + cmd.Name(), false},
		{"-o, --onload [STACK]", get("onload"), false},
		{"--connections", get("connections"), true},
		{"--errors", get("errors"), true},
		{"--interrupts", get("interrupts"), true},
		{"    --pid-resolution", get("pid-resolution"), false},
		{"-p, --process int", get("process"), false},
		{"-s, --stats", get("stats"), false},
		{"    --show-all-affinities", get("show-all-affinities"), false},
	}

	fmt.Println("\nFlags:")
	for _, e := range flags {
		var prefix string
		if e.sub {
			prefix = "        " // 8 spaces — visually nested under --onload
		} else {
			prefix = "  "
		}
		left := prefix + e.left
		pad := usageCol - len(left)
		if pad < 2 {
			pad = 2
		}
		fmt.Printf("%s%s%s\n", left, strings.Repeat(" ", pad), e.usage)
	}
	fmt.Println()
}
