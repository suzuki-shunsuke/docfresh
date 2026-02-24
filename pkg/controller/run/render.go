package run

import (
	"bytes"
	"context"
	"text/template"
)

func (c *Controller) renderBlock(ctx context.Context, block *Block) (string, error) {
	if block.Type == "text" {
		return block.Content, nil
	}
	tpl, err := template.New("_").Funcs(txtFuncMap()).Parse(commandTemplate)
	if err != nil {
		return "", err
	}
	content := block.BeginComment
	result, err := c.exec(ctx, block)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, TemplateInput{
		Command:        block.Input.Command.Command,
		Stdout:         result.Stdout,
		Stderr:         result.Stderr,
		CombinedOutput: result.CombinedOutput,
		ExitCode:       result.ExitCode,
	}); err != nil {
		return "", err
	}
	content += "\n" + buf.String() + block.EndComment
	return content, nil
}
