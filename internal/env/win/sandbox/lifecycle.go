package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (sandbox *Sandbox) start() error {
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

func (sandbox *Sandbox) connect() error {
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
		return fmt.Errorf("failed to connect to the development windows sandbox: %w", err)
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

	return fmt.Errorf("failed to identify windows sandbox remote session process")
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
			"failed to query windows sandbox remote session processes: %w",
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

func (sandbox *Sandbox) stop() error {
	fmt.Println("Stopping the development Windows Sandbox...")
	cmd := exec.Command("wsb", "stop",
		"--id", sandbox.id,
		"--raw",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop development windows sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) isRunning() (bool, error) {
	if sandbox.remoteSessionPID == 0 {
		return false, nil
	}
	pids, err := getRemoteSessionPIDs()
	if err != nil {
		return false, err
	}

	_, exists := pids[sandbox.remoteSessionPID]

	return exists, nil
}

func (sandbox *Sandbox) installVCRuntime(dependenciesDir string) error {
	installer := filepath.Join(dependenciesDir, "VC_redist.x64.exe")

	command := fmt.Sprintf(
		`"%s" /install /quiet /norestart`,
		installer,
	)

	if err := sandbox.RunCommand(command); err != nil {
		return fmt.Errorf("failed to install Visual C++ Runtime in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) installWebView2(dependenciesDir string) error {
	installer := filepath.Join(
		dependenciesDir,
		"MicrosoftEdgeWebView2RuntimeInstallerX64.exe",
	)

	command := fmt.Sprintf(
		`"%s" /silent /install`,
		installer,
	)

	if err := sandbox.RunCommand(command); err != nil {
		return fmt.Errorf("failed to install WebView2 Runtime in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) installNotepadPlusPlus(dependenciesDir string) error {
	installer := filepath.Join(
		dependenciesDir,
		"npp.8.9.8.Installer.x64.exe",
	)

	command := fmt.Sprintf(
		`"%s" /S`,
		installer,
	)

	if err := sandbox.RunCommand(command); err != nil {
		return fmt.Errorf("failed to install Notepad++ in Windows Sandbox: %w", err)
	}

	return nil
}

func (sandbox *Sandbox) setupEnv() error {
	if err := sandbox.ShareFolder(
		sandbox.Config.HostDependenciesPath,
		sandbox.Config.DependenciesPath,
	); err != nil {
		return fmt.Errorf("failed to share Sandbox dependencies: %w", err)
	}

	if err := sandbox.installVCRuntime(sandbox.Config.DependenciesPath); err != nil {
		return err
	}

	if err := sandbox.installWebView2(sandbox.Config.DependenciesPath); err != nil {
		return err
	}

	if err := sandbox.installNotepadPlusPlus(sandbox.Config.DependenciesPath); err != nil {
		return err
	}

	return nil
}

