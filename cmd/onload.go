package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func runOnload() error {
	// Check binary is available
	bin, err := exec.LookPath("onload_stackdump")
	if err != nil {
		return fmt.Errorf("onload_stackdump not found in PATH — is OpenOnload installed?")
	}

	// Check the kernel module is loaded (/dev/onload only exists when it is)
	if _, err := os.Stat("/dev/onload"); os.IsNotExist(err) {
		return fmt.Errorf("OpenOnload is installed but not running (no /dev/onload)")
	}

	// Require root
	if os.Geteuid() != 0 {
		return fmt.Errorf("onload_stackdump requires root — re-run with: sudo cpu-tools --onload")
	}

	cmd := exec.Command(bin)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
