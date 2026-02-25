package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

func (c *Controller) renderBlock(ctx context.Context, tpls *Templates, file string, block *Block) (string, error) {
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
	result, err := c.exec(ctx, file, block.Input)
	if err != nil {
		return "", fmt.Errorf("execute a command: %w", err)
	}
	s, err := c.render(tpl, result)
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	if strings.HasSuffix(s, "\n") {
		return content + "\n" + s + block.EndComment, nil
	}
	return content + "\n" + s + "\n" + block.EndComment, nil
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
