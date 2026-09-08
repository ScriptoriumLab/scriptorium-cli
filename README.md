# Scriptorium CLI

[![Scriptorium CLI Build CI](https://github.com/ScriptoriumLab/scriptorium-cli/actions/workflows/scriptorium-cli-build.yml/badge.svg)](https://github.com/ScriptoriumLab/scriptorium-cli/actions/workflows/scriptorium-cli-build.yml)

[![Scriptorium CLI Release CI](https://github.com/ScriptoriumLab/scriptorium-cli/actions/workflows/scriptorium-cli-release.yml/badge.svg)](https://github.com/ScriptoriumLab/scriptorium-cli/actions/workflows/scriptorium-cli-release.yml)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

> A cross-platform developer toolchain for building, running, testing, and maintaining Scriptorium development environments.

Scriptorium CLI provides a single entry point for common Scriptorium development workflows.

The command-line tool is exposed as **`orium`**.

---

## Goals

- Provide a consistent entry point for development workflows
- Reduce repetitive setup and cleanup
- Coordinate builds and tests across Scriptorium components
- Run Scriptorium inside isolated development environments
- Keep workflows reproducible for contributors
- Evolve into a long-lived, cross-platform developer toolchain

---

## Technology

Scriptorium CLI is implemented in **Go** and uses [Cobra](https://github.com/spf13/cobra) for command-line parsing.

The project grows incrementally from real development needs rather than introducing abstractions upfront.

---

## Building

Development currently focuses on **Windows**.

```powershell
go build .\cmd\orium\
```

Then:

```powershell
.\orium.exe --help
```

---

## `orium dev`

`orium dev` builds and tests the Scriptorium workspace, prepares a development environment, deploys the product, starts the runtime, and waits for manual testing.

```powershell
.\orium.exe dev
```

Select an environment explicitly with:

```powershell
.\orium.exe dev --env sandbox
.\orium.exe dev --env vm
```

The current default is:

```text
sandbox
```

### Windows Sandbox

Windows Sandbox provides a clean and disposable development environment.

The current workflow prepares runtime dependencies, deploys Brush, Inkstone, Ink, and the dictionary, starts Scriptorium, opens a manual test environment, monitors the Sandbox session, and cleans temporary host-side state afterwards.

### VMware VM

The VMware environment restores a Windows development VM to a known baseline snapshot, deploys the latest locally built artifacts, starts the development workflow, monitors the session, and restores the VM afterwards.

---

## Architecture

The current `dev` workflow is split into three layers:

```text
dev command orchestration
        ↓
environment-specific development workflow
        ↓
environment infrastructure
```

The top-level `dev` command owns workspace and product orchestration.

Environment-specific commands live under:

```text
internal/command/dev/env/win/sandbox
internal/command/dev/env/win/vm
```

Low-level infrastructure lives under:

```text
internal/env/win/sandbox
internal/env/win/vm
```

This keeps Scriptorium-specific workflow logic separate from VM and Sandbox mechanics.

---

## Planned Commands

```text
orium test
orium doctor
orium health
orium perf
orium clean
```

These commands will be added as real development needs require them.

---

## Cross-Platform Direction

Windows is the first supported platform because Scriptorium currently integrates with the Windows Text Services Framework (TSF).

The CLI is intended to remain the common developer entry point as Scriptorium expands to macOS and potentially Linux.

---

## Status

`orium dev` currently supports end-to-end development workflows using:

- Windows Sandbox
- VMware Workstation

The CLI is under active development. Current areas of focus include lifecycle reliability, developer feedback, packaging, and additional developer workflows.

---

## License

Licensed under the **Apache License 2.0**.

See `LICENSE` for details.

---

*Copyright © 2026 ScriptoriumLab.*
