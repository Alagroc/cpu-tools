# cpu-tools

Linux CLI for inspecting CPU affinity and usage of running processes.

## Usage

```
cpu-tools [flags]

Flags:
  -a, --all                  list processes on selected cores
  -c, --cores string         core range, e.g. 0-3 or 0,1,2,3
  -p, --process int          PID to inspect
  -s, --show                 show CPU usage (-a: percentage, -p: current usage)
      --pid-resolution       show full cmdline instead of process name
      --show-all-affinities  include processes with affinity spanning all cores
  -h, --help                 help for cpu-tools
```

## Examples

List processes exclusively pinned to cores 0–3 (excludes full-affinity processes):
```
cpu-tools -a -c 0-3
```

Same, sorted by CPU% (samples for 500ms):
```
cpu-tools -a -c 0-3 -s
```

Include processes that can run on any core:
```
cpu-tools -a -c 0-3 --show-all-affinities
```

Show full command line instead of process name:
```
cpu-tools -a -c 0-3 --pid-resolution
```

Show affinity and current CPU usage for a process:
```
cpu-tools -p <PID> -s
```

## Build

```
go build -o cpu-tools .
```

Requires Go 1.25+ and Linux.
