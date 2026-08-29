package proc

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CPUPressure reads cpu.pressure for the cgroup of pid.
// Returns empty string if the cgroup or the file doesn't exist.
func CPUPressure(pid int) (string, error) {
	cgPath, err := cgroupV2Path(pid)
	if err != nil || cgPath == "" {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join("/sys/fs/cgroup", cgPath, "cpu.pressure"))
	if err != nil {
		return "", nil // pressure file simply absent — not an error
	}
	return strings.TrimSpace(string(data)), nil
}

// cgroupV2Path returns the unified-hierarchy (v2) cgroup path for pid,
// e.g. "user.slice/user-1000.slice/session-1.scope".
func cgroupV2Path(pid int) (string, error) {
	f, err := os.Open(filepath.Join("/proc", strconv.Itoa(pid), "cgroup"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// cgroup v2 line: "0::<path>"
		parts := strings.SplitN(scanner.Text(), ":", 3)
		if len(parts) == 3 && parts[0] == "0" && parts[1] == "" {
			return strings.TrimLeft(parts[2], "/"), nil
		}
	}
	return "", scanner.Err()
}
