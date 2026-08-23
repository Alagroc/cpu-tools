package proc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ProcStat struct {
	PID       int
	Utime     uint64
	Stime     uint64
	Processor int
}

func ReadStat(pid int) (ProcStat, error) {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return ProcStat{}, err
	}
	closeIdx := strings.LastIndex(string(data), ")")
	if closeIdx < 0 {
		return ProcStat{}, fmt.Errorf("invalid stat format for PID %d", pid)
	}
	fields := strings.Fields(string(data)[closeIdx+1:])
	if len(fields) < 37 {
		return ProcStat{}, fmt.Errorf("stat has too few fields for PID %d", pid)
	}
	utime, _ := strconv.ParseUint(fields[11], 10, 64)
	stime, _ := strconv.ParseUint(fields[12], 10, 64)
	processor, _ := strconv.Atoi(fields[36])
	return ProcStat{PID: pid, Utime: utime, Stime: stime, Processor: processor}, nil
}

func CPUPercent(pid int, d time.Duration) (float64, error) {
	clkTck := clkTicks()
	s1, err := ReadStat(pid)
	if err != nil {
		return 0, err
	}
	t1 := time.Now()
	time.Sleep(d)
	s2, err := ReadStat(pid)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(t1).Seconds()
	delta := float64((s2.Utime + s2.Stime) - (s1.Utime + s1.Stime))
	return delta / (elapsed * float64(clkTck)) * 100, nil
}

func CPUPercentBatch(pids []int, d time.Duration) (map[int]float64, error) {
	clkTck := clkTicks()

	t1Stats := make(map[int]ProcStat, len(pids))
	for _, pid := range pids {
		s, err := ReadStat(pid)
		if err != nil {
			continue
		}
		t1Stats[pid] = s
	}

	t1 := time.Now()
	time.Sleep(d)
	elapsed := time.Since(t1).Seconds()

	result := make(map[int]float64, len(pids))
	for _, pid := range pids {
		s1, ok := t1Stats[pid]
		if !ok {
			continue
		}
		s2, err := ReadStat(pid)
		if err != nil {
			continue
		}
		delta := float64((s2.Utime + s2.Stime) - (s1.Utime + s1.Stime))
		result[pid] = delta / (elapsed * float64(clkTck)) * 100
	}
	return result, nil
}

func clkTicks() int64 {
	return 100 // USER_HZ is always 100 on Linux
}
