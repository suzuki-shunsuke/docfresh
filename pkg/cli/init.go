package cli

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/docfresh/pkg/controller/initcmd"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

// InitArgs holds the flag and argument values for the init command.
type InitArgs struct {
	*Flags

	ConfigFilePath string // positional argument
}

// NewInit creates a new init command instance with the provided logger.
// It returns a CLI command that can be registered with the main CLI application.
func NewInit(logger *slogutil.Logger, gFlags *Flags) *cli.Command {
	args := &InitArgs{
		Flags: gFlags,
	}
	return &cli.Command{
		Name:  "init",
		Usage: "Create docfresh.yaml if it doesn't exist",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return action(ctx, logger, args)
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "config-file-path",
				Destination: &args.ConfigFilePath,
			},
		},
	}
}

func action(_ context.Context, logger *slogutil.Logger, args *InitArgs) error {
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	configFilePath := args.ConfigFilePath
	if configFilePath == "" {
		configFilePath = args.Config
	}
	if configFilePath == "" {
		configFilePath = "docfresh.yaml"
	}
	fs := afero.NewOsFs()
	ctrl := initcmd.New(fs)
	return ctrl.Init(logger.Logger, configFilePath) //nolint:wrapcheck
}
