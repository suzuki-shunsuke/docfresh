package run

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
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
	blocks, err := c.parseFile(string(b))
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

// parseFile parses a file and returns a list of blocks.
func (c *Controller) parseFile(_ string) ([]*Block, error) {
	return []*Block{
		{
			Type: "text",
			Content: `# Hello

`,
		},
		{
			Type: "block",
			BeginComment: `<!-- docfresh begin
command:
  command: echo "Hello"
-->`,
			EndComment: `<!-- docfresh end -->`,
			Input: &BlockInput{
				Command: &Command{
					Command: `echo "Hello"`,
				},
			},
		},
	}, nil
}

type TemplateInput struct {
	Command        string
	Stdout         string
	Stderr         string
	CombinedOutput string
	ExitCode       int
}

func (c *Controller) renderBlock(ctx context.Context, block *Block) (string, error) {
	if block.Type == "text" {
		return block.Content, nil
	}
	fncs := sprig.TxtFuncMap()
	delete(fncs, "env")
	delete(fncs, "expandenv")
	delete(fncs, "getHostByName")
	tpl, err := template.New("_").Funcs(fncs).Parse(commandTemplate)
	if err != nil {
		return "", err
	}
	content := block.BeginComment
	shell := block.Input.Command.Shell
	if shell == nil {
		shell = []string{"bash", "-c"}
	}
	cmd := exec.CommandContext(ctx, shell[0], append(shell[1:], block.Input.Command.Command)...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	setCancel(cmd)
	fmt.Fprintln(os.Stderr, "+", block.Input.Command.Command)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, TemplateInput{
		Command:        block.Input.Command.Command,
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       cmd.ProcessState.ExitCode(),
	}); err != nil {
		return "", err
	}
	content += "\n" + buf.String() + block.EndComment
	return content, nil
}

const waitDelay = 1000 * time.Hour

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = waitDelay
}
