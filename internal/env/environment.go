// Package env provides an interface for managing the environment in which the application runs.
package env

type Environment interface {
    EnsureAvailable() error
    Prepare() error
    Monitor() error
    Cleanup() error
}

