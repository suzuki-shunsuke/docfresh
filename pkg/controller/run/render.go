package run

import (
	"bytes"
	"context"
	"fmt"
)

func (c *Controller) renderBlock(ctx context.Context, tpls *Templates, file string, block *Block) (string, error) {
	if block.Type == "text" {
		return block.Content, nil
	}
	content := block.BeginComment
	result, err := c.exec(ctx, file, block.Input)
	if err != nil {
		return "", fmt.Errorf("execute a command: %w", err)
	}
	s, err := c.render(tpls, result)
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	content += "\n" + s + block.EndComment
	return content, nil
}

func (c *Controller) render(tpls *Templates, result *TemplateInput) (string, error) {
	switch result.Type {
	case "local-file":
		return result.Content, nil
	case "command":
		buf := &bytes.Buffer{}
		if err := tpls.Command.Execute(buf, result); err != nil {
			return "", fmt.Errorf("execute a template: %w", err)
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("unknown type: %s", result.Type)
	}
}
