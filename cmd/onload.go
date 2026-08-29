package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func runOnload() error {
	bin, err := exec.LookPath("onload_stackdump")
	if err != nil {
		return fmt.Errorf("onload_stackdump not found in PATH — is OpenOnload installed?")
	}
	if _, err := os.Stat("/dev/onload"); os.IsNotExist(err) {
		return fmt.Errorf("OpenOnload is installed but not running (no /dev/onload)")
	}
	if os.Geteuid() != 0 {
		return fmt.Errorf("onload_stackdump requires root — re-run with: sudo cpu-tools --onload")
	}

	if flagOnload == "all" {
		return runOnloadFull(bin)
	}

	stack, err := strconv.Atoi(flagOnload)
	if err != nil || stack < 0 {
		return fmt.Errorf("invalid stack ID %q — must be a non-negative integer", flagOnload)
	}

	switch {
	case flagInterrupts:
		return runOnloadInterrupts(bin, stack)
	case flagConnections:
		return runOnloadConnections(bin, stack)
	case flagErrors:
		return runOnloadErrors(bin, stack)
	default:
		fmt.Fprintf(os.Stderr, "Stack %d selected — what do you want to check?\n\n", stack)
		fmt.Fprintf(os.Stderr, "  cpu-tools --onload=%d --interrupts    analyse interrupt rate (10s sample)\n", stack)
		fmt.Fprintf(os.Stderr, "  cpu-tools --onload=%d --connections   list ESTABLISHED connections\n", stack)
		fmt.Fprintf(os.Stderr, "  cpu-tools --onload=%d --errors        show non-zero error counters\n", stack)
		return nil
	}
}

func runOnloadFull(bin string) error {
	cmd := exec.Command(bin)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runOnloadInterrupts(bin string, stack int) error {
	sample := func() (int64, string, error) {
		out, err := exec.Command("sh", "-c",
			fmt.Sprintf("%s %d lots | grep -i interrupt", bin, stack)).Output()
		if err != nil {
			return 0, "", fmt.Errorf("onload_stackdump failed: %w", err)
		}
		return sumInts(string(out)), string(out), nil
	}

	fmt.Printf("Sampling interrupts for stack %d...\n\n", stack)

	count1, raw1, err := sample()
	if err != nil {
		return err
	}
	t1 := time.Now()

	if raw1 == "" {
		return fmt.Errorf("no interrupt counters found for stack %d", stack)
	}

	fmt.Printf("=== t0 ===\n%s\n", raw1)
	fmt.Printf("Waiting 10 seconds...\n\n")
	time.Sleep(10 * time.Second)

	count2, raw2, err := sample()
	if err != nil {
		return err
	}
	elapsed := time.Since(t1).Seconds()

	fmt.Printf("=== t1 ===\n%s\n", raw2)

	delta := count2 - count1
	rate := float64(delta) / elapsed

	fmt.Printf("─────────────────────────────────\n")
	fmt.Printf("Delta:    %d interrupts in %.1fs\n", delta, elapsed)
	fmt.Printf("Rate:     %.1f interrupts/sec\n", rate)

	switch {
	case rate > 50000:
		fmt.Printf("Verdict:  HIGH — stack is receiving heavy interrupt load\n")
	case rate > 5000:
		fmt.Printf("Verdict:  MODERATE — interrupt load is elevated\n")
	default:
		fmt.Printf("Verdict:  LOW — interrupt rate is normal\n")
	}

	return nil
}

func runOnloadConnections(bin string, stack int) error {
	out, err := exec.Command(bin, strconv.Itoa(stack), "netstat").Output()
	if err != nil {
		return fmt.Errorf("onload_stackdump netstat failed: %w", err)
	}

	var established []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "ESTABLISHED") {
			established = append(established, line)
		}
	}

	if len(established) == 0 {
		fmt.Printf("No ESTABLISHED connections on stack %d\n", stack)
		return nil
	}

	fmt.Printf("ESTABLISHED connections on stack %d (%d total):\n\n", stack, len(established))
	for _, line := range established {
		fmt.Println(line)
	}
	return nil
}

func runOnloadErrors(bin string, stack int) error {
	type source struct {
		label string
		arg   string
	}
	sources := []source{
		{"Stack stats", "stats"},
		{"VI stats", "vi_stats"},
		{"TCP stats", "tcp_stats"},
		{"UDP stats", "udp_stats"},
	}

	totalErrors := 0
	for _, src := range sources {
		out, err := exec.Command(bin, strconv.Itoa(stack), src.arg).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: onload_stackdump %d %s failed: %v\n", stack, src.arg, err)
			continue
		}

		var active []string
		for _, line := range strings.Split(string(out), "\n") {
			lower := strings.ToLower(line)
			if !reErrField.MatchString(lower) {
				continue
			}
			// Skip lines where the counter value is zero.
			if reZeroValue.MatchString(line) {
				continue
			}
			active = append(active, strings.TrimSpace(line))
		}

		fmt.Printf("=== %s (stack %d) ===\n", src.label, stack)
		if len(active) == 0 {
			fmt.Println("  no errors")
		} else {
			for _, l := range active {
				fmt.Printf("  %s\n", l)
			}
		}
		fmt.Println()
		totalErrors += len(active)
	}

	if totalErrors == 0 {
		fmt.Printf("✓ No errors detected on stack %d\n", stack)
	} else {
		fmt.Printf("⚠ %d error counter(s) with non-zero values on stack %d\n", totalErrors, stack)
	}
	return nil
}

var (
	reErrField  = regexp.MustCompile(`err|drop|fail|bad`)
	// sumInts sums all integers found in s, used to aggregate interrupt counters.
	reInt = regexp.MustCompile(`\d+`)
	// matches lines where the counter value is 0, e.g. "field : 0" or "field=0"
	reZeroValue = regexp.MustCompile(`[=:]\s*0\s*$`)
)

func sumInts(s string) int64 {
	var total int64
	for _, line := range strings.Split(s, "\n") {
		matches := reInt.FindAllString(line, -1)
		if len(matches) == 0 {
			continue
		}
		// Take the last number on the line — typically the counter value.
		if n, err := strconv.ParseInt(matches[len(matches)-1], 10, 64); err == nil {
			total += n
		}
	}
	return total
}
