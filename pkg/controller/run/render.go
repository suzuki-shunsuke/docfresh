package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Controller) renderBlock(ctx context.Context, logger *slog.Logger, tpls *Templates, file string, block *Block) (gS string, gErr error) {
	if block.Type == "text" {
		return block.Content, nil
	}
	if block.Input == nil {
		return "", errors.New("block input is nil")
	}
	tpl, err := c.getTemplate(tpls, block)
	if err != nil {
		return "", err
	}
	content := block.BeginComment
	if block.Input.PostCommand != nil {
		defer func() {
			if _, err := c.execCommand(ctx, file, block.Input.PostCommand); err != nil {
				if gErr == nil {
					gErr = fmt.Errorf("execute post_command: %w", err)
					return
				}
				slogerr.WithError(logger, err).Error("execute post_command")
			}
		}()
	}
	if err := c.runPreCommand(ctx, file, block); err != nil {
		return "", fmt.Errorf("execute pre_command: %w", err)
	}
	result, err := c.exec(ctx, file, block.Input)
	if err != nil {
		return "", fmt.Errorf("execute a command: %w", err)
	}
	s, err := c.render(tpl, result)
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return appendEndComment(content, s, block.EndComment), nil
}

func appendEndComment(content, s, endComment string) string {
	if strings.HasSuffix(s, "\n") {
		return content + "\n" + s + endComment
	}
	return content + "\n" + s + "\n" + endComment
}

func (c *Controller) runPreCommand(ctx context.Context, file string, block *Block) error {
	if block.Input.PreCommand == nil {
		return nil
	}
	if _, err := c.execCommand(ctx, file, block.Input.PreCommand); err != nil {
		return err
	}
	return nil
}

func (c *Controller) render(tpl *template.Template, result *TemplateInput) (string, error) {
	switch result.Type {
	case "local-file", "http":
		return result.Content, nil
	case "command":
		buf := &bytes.Buffer{}
		if err := tpl.Execute(buf, result); err != nil {
			return "", fmt.Errorf("execute a template: %w", err)
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("unknown type: %s", result.Type)
	}
}

func (c *Controller) getTemplate(tpls *Templates, block *Block) (*template.Template, error) {
	if block.Input.Template == nil {
		if block.Input.Command != nil {
			return tpls.Command, nil
		}
		return nil, nil //nolint:nilnil
	}
	tpl, err := template.New("_").Funcs(tpls.Funcs).Parse(block.Input.Template.Content)
	if err != nil {
		return nil, fmt.Errorf("parse block template: %w", err)
	}
	return tpl, nil
}
