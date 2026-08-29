# cpu-tools

Linux CLI for inspecting CPU affinity, usage, and OpenOnload stack state.

## Usage

```
cpu-tools [flags]

Flags:
  -a, --all                    list processes on selected cores
  -c, --cores string           core range, e.g. 0-3 or 0,1,2,3
  -h, --help                   help for cpu-tools
  -o, --onload [STACK]         dump OpenOnload stacks (requires root)
        --connections          list ESTABLISHED connections for --onload STACK
        --errors               show non-zero error counters for --onload STACK
                                 (stats, more_stats, vi_stats, ip_stats, tcp_stats, udp_stats)
        --interrupts           analyse interrupt rate for --onload STACK (10s sample)
      --pid-resolution         show full cmdline instead of process name
  -p, --process int            PID to inspect
  -s, --show                   show CPU usage (-a: percentage, -p: current usage)
      --show-all-affinities    include processes with affinity spanning all cores
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

Dump all OpenOnload stacks:
```
sudo cpu-tools -o
```

List ESTABLISHED connections on onload stack 1:
```
sudo cpu-tools -o 1 --connections
```

Analyse interrupt rate on onload stack 1 (samples over 10s, shows rate and verdict):
```
sudo cpu-tools -o 1 --interrupts
```

Show non-zero error counters on onload stack 1:
```
sudo cpu-tools -o 1 --errors
```

## Build

```
go build -o cpu-tools .
```

Requires Go 1.25+ and Linux.
