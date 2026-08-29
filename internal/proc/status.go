package proc

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcStatus struct {
	VoluntaryCtxtSwitches    int64
	NonVoluntaryCtxtSwitches int64
}

func ReadStatus(pid int) (ProcStatus, error) {
	f, err := os.Open(filepath.Join("/proc", strconv.Itoa(pid), "status"))
	if err != nil {
		return ProcStatus{}, err
	}
	defer f.Close()

	var s ProcStatus
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}
		val, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}
		switch strings.TrimSpace(parts[0]) {
		case "voluntary_ctxt_switches":
			s.VoluntaryCtxtSwitches = val
		case "nonvoluntary_ctxt_switches":
			s.NonVoluntaryCtxtSwitches = val
		}
	}
	return s, scanner.Err()
}
