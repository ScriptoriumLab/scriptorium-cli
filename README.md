# Scriptorium CLI

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

> A cross-platform developer toolchain for building, running, testing, and maintaining Scriptorium development environments.

Scriptorium CLI is designed to provide a single, consistent entry point for common development workflows.

Instead of manually switching between repositories, running build commands, starting processes, registering components, and cleaning up the environment, the CLI orchestrates these tasks through a small set of developer-focused commands.

The goal is simple:

> **Make local development reproducible, observable, and easy to operate.**

The command-line tool is exposed as **`orium`**, a short and distinctive name derived from Scriptorium.

---

## Goals

Scriptorium CLI aims to:

- Provide a single entry point for common development workflows
- Reduce repetitive manual setup and cleanup
- Coordinate builds and tests across multiple components
- Manage development process lifecycles
- Provide clear aggregated feedback when something fails
- Keep development workflows reproducible for contributors
- Evolve into a long-lived, cross-platform developer toolchain

---

## Technology

Scriptorium CLI is implemented in **Go**.

Command-line parsing and command organization are built with [Cobra](https://github.com/spf13/cobra).

The project intentionally starts with a small structure and will grow incrementally as real development workflows introduce new responsibilities and boundaries.

---

## Building

Development is currently focused on **Windows**.

To build the CLI on Windows:

```powershell
go build -o orium.exe
```

Then run:

```powershell
.\orium.exe --help
```

For example:

```powershell
.\orium.exe greet
```

Support for additional platforms will be added as the toolchain evolves.

---

## Planned Commands

The command set will evolve as development needs grow.

Initial ideas include:

```text
orium dev
orium test
orium doctor
orium health
orium perf
orium clean
```

### `orium dev`

Prepare and start a complete local development environment.

The development workflow runs Scriptorium inside an isolated, disposable environment so that platform components such as the Windows TSF integration can be registered and exercised without modifying the host development environment unnecessarily.

The environment can be selected with:

```shell
orium dev --env <environment>
```

The current default is:

```shell
orium dev --env vm
```

Running:

```shell
orium dev
```

is therefore currently equivalent to:

```shell
orium dev --env vm
```

The first implementation uses a VMware-based Windows virtual machine as the disposable development environment.

Windows Sandbox and Hyper-V were explored during the initial implementation. Both provide useful isolation, but they are more tightly coupled to the host Windows edition and Windows virtualization stack. The current development machine cannot reliably use Windows Sandbox, and requiring it would also exclude contributors whose Windows environments do not support it.

VMware was therefore selected as the initial implementation because it provides a more broadly available and predictable environment for contributors.

This is an implementation choice rather than a permanent architectural constraint.

The `--env` option intentionally leaves room for additional development environments in the future, for example:

```shell
orium dev --env sandbox
orium dev --env hyperv
```

These environments are not currently implemented.

Future contributors who prefer Windows Sandbox, Hyper-V, or another suitable environment can add support without changing the user-facing purpose of `orium dev`.

The initial workflow is expected to:
- prepare or reset the disposable development environment
- make locally built Scriptorium artifacts available to it
- register required platform components
- start the Scriptorium runtime
- wait while manual testing is performed
- discard or reset the environment when the session ends

The environment strategy is expected to continue evolving as the development workflow and contributor needs become clearer.

### `orium test`

Build relevant components and run their test suites, presenting a single aggregated result.

### `orium doctor`

Inspect the local development environment and report missing or invalid prerequisites.

### `orium health`

Inspect a running development environment and report the current health of its components.

### `orium perf`

Run performance-related workflows and benchmarks.

### `orium clean`

Remove temporary development state and restore the environment to a clean state.

---

## Design Direction

Scriptorium CLI is intended to be more than a collection of shell scripts.

It is treated as a long-lived engineering project with explicit responsibilities, predictable behavior, testable workflows, and room for future evolution.

The implementation will grow incrementally from real development needs rather than attempting to predict every future use case upfront.

Cross-platform support is part of the long-term direction. Development currently focuses on Windows, with support for additional platforms planned as the broader toolchain evolves.

---

## Status

This project is in its very early stages.

The first milestone is to establish the basic CLI and implement the initial Windows development workflow.

Commands, structure, and behavior are expected to evolve rapidly during the early versions.

---

## License

Licensed under the **Apache License 2.0**.

See `LICENSE` for details.

---

*Copyright © 2026 ScriptoriumLab.*
