package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func (sandbox *Sandbox) Start() error {
	fmt.Println("Starting the development Windows Sandbox...")
	cmd := exec.Command("wsb", "start",
		"--id", sandbox.id,
		"--raw",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start development windows sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) Connect() error {
	fmt.Println("Connecting to the development Windows Sandbox...")

	before, err := getRemoteSessionPIDs()
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"wsb", "connect",
		"--id", sandbox.id,
		"--raw",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to connect to the development Windows Sandbox: %w", err)
	}

	for range 10 {
		after, err := getRemoteSessionPIDs()
		if err != nil {
			return err
		}

		for pid := range after {
			if _, existed := before[pid]; !existed {
				sandbox.remoteSessionPID = pid
				return nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("failed to identify Windows Sandbox remote session process")
}

func getRemoteSessionPIDs() (map[int]struct{}, error) {
	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-Command",
		`
$p = Get-Process WindowsSandboxRemoteSession -ErrorAction SilentlyContinue
if ($p) {
    $p | Select-Object -ExpandProperty Id
}
exit 0
`,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to query Windows Sandbox remote session processes: %w",
			err,
		)
	}

	pids := make(map[int]struct{})

	for field := range strings.FieldsSeq(string(output)) {
		pid, err := strconv.Atoi(field)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse remote session PID %q: %w",
				field,
				err,
			)
		}

		pids[pid] = struct{}{}
	}

	return pids, nil
}
