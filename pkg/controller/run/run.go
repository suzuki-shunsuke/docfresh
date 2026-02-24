package run

import (
	"context"
	_ "embed"
	"log/slog"
	"os"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

// File and directory permissions for created configuration files
const (
	filePermission os.FileMode = 0o644 // Standard file permissions (rw-r--r--)
	dirPermission  os.FileMode = 0o755 // Standard directory permissions (rwxr-xr-x)
)

var (
	//go:embed command_template.md
	commandTemplate string
)

type Input struct {
	ConfigFilePath string
	Files          map[string]struct{}
}

func (c *Controller) Run(ctx context.Context, logger *slog.Logger, input *Input) error {
	for file := range input.Files {
		logger := logger.With("file", file)
		if err := c.run(ctx, logger, file); err != nil {
			return slogerr.With(err, "file", file)
		}
	}
	return nil
}

func (c *Controller) run(ctx context.Context, logger *slog.Logger, file string) error {
	b, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return err
	}
	bs := string(b)
	blocks, err := parseFile(string(b))
	if err != nil {
		return err
	}
	content := ""
	for _, block := range blocks {
		s, err := c.renderBlock(ctx, block)
		if err != nil {
			return err
		}
		content += s
	}
	if content != bs {
		if err := afero.WriteFile(c.fs, file, []byte(content), filePermission); err != nil {
			return err
		}
	}

	return nil
}

type Block struct {
	// text, code block
	Type         string
	Content      string
	Input        *BlockInput
	BeginComment string
	EndComment   string
}

type BlockInput struct {
	Command *Command
}

type Command struct {
	Command string
	Shell   []string
}

type TemplateInput struct {
	Command        string
	Stdout         string
	Stderr         string
	CombinedOutput string
	ExitCode       int
}
