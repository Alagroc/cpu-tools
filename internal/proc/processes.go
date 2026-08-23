package proc

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Process struct {
	PID  int
	Name string
}

func AllProcesses() ([]Process, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var procs []Process
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		name, _ := readComm(pid)
		procs = append(procs, Process{PID: pid, Name: name})
	}
	return procs, nil
}

func GetName(pid int) (string, error) {
	return readComm(pid)
}

func GetCmdline(pid int) (string, error) {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	if err != nil {
		return "", err
	}
	// cmdline is null-byte delimited; last byte is also null
	return strings.ReplaceAll(strings.TrimRight(string(data), "\x00"), "\x00", " "), nil
}

func readComm(pid int) (string, error) {
	data, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "comm"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
