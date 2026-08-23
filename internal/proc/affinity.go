package proc

import "golang.org/x/sys/unix"

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
