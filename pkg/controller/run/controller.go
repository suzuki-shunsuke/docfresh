// Package run implements the business logic for the 'docfresh run' command.
package run

import (
	"github.com/spf13/afero"
)

// Controller manages the initialization of docfresh configuration.
// It provides methods to create configuration files with appropriate permissions.
type Controller struct {
	fs afero.Fs
}

// New creates a new Controller instance with the provided filesystem and environment.
// The filesystem is used for all file operations, allowing for easy testing with mock filesystems.
func New(fs afero.Fs) *Controller {
	return &Controller{
		fs: fs,
	}
}
