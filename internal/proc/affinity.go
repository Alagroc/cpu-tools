package proc

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// SystemCPUCount reads /sys/devices/system/cpu/possible to get the true
// number of CPUs on the system, independent of Go runtime or cgroup limits.
func SystemCPUCount() int {
	data, err := os.ReadFile("/sys/devices/system/cpu/possible")
	if err != nil {
		return runtime.NumCPU()
	}
	count := 0
	for _, part := range strings.Split(strings.TrimSpace(string(data)), ",") {
		if idx := strings.Index(part, "-"); idx >= 0 {
			lo, e1 := strconv.Atoi(part[:idx])
			hi, e2 := strconv.Atoi(part[idx+1:])
			if e1 == nil && e2 == nil {
				count += hi - lo + 1
			}
		} else if _, e := strconv.Atoi(strings.TrimSpace(part)); e == nil {
			count++
		}
	}
	if count == 0 {
		return runtime.NumCPU()
	}
	return count
}

func GetAffinity(pid int) ([]int, error) {
	var set unix.CPUSet
	if err := unix.SchedGetaffinity(pid, &set); err != nil {
		return nil, err
	}
	cpus := make([]int, 0, set.Count())
	for i := range 1024 {
		if set.IsSet(i) {
			cpus = append(cpus, i)
		}
	}
	return cpus, nil
}
