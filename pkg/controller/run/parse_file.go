package run

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	beginMarker = "<!-- docfresh begin"
	endMarker   = "<!-- docfresh end -->"
)

// parseFile parses a file and returns a list of blocks.
func parseFile(content string) ([]*Block, error) { //nolint:cyclop
	var blocks []*Block
	pos := 0
	for pos < len(content) {
		beginIdx := strings.Index(content[pos:], beginMarker)
		endIdx := strings.Index(content[pos:], endMarker)

		// No more markers — emit remaining text and break.
		if beginIdx == -1 && endIdx == -1 {
			blocks = appendText(blocks, content[pos:])
			break
		}

		// end before begin (or end without begin).
		if beginIdx == -1 || (endIdx != -1 && endIdx < beginIdx) {
			return nil, errors.New("found <!-- docfresh end --> without a matching <!-- docfresh begin")
		}

		// Emit text before the begin marker.
		if beginIdx > 0 {
			blocks = appendText(blocks, content[pos:pos+beginIdx])
		}

		// Find closing --> of the begin comment.
		beginStart := pos + beginIdx
		closeIdx := strings.Index(content[beginStart+len(beginMarker):], "-->")
		if closeIdx == -1 {
			return nil, errors.New("unclosed <!-- docfresh begin comment: missing -->")
		}
		beginCommentEnd := beginStart + len(beginMarker) + closeIdx + len("-->")
		beginComment := content[beginStart:beginCommentEnd]

		// Extract YAML from inside the begin comment.
		yamlStr := content[beginStart+len(beginMarker) : beginStart+len(beginMarker)+closeIdx]
		yamlStr = strings.TrimSpace(yamlStr)
		var input BlockInput
		if err := yaml.Unmarshal([]byte(yamlStr), &input); err != nil {
			return nil, fmt.Errorf("failed to parse YAML in begin comment: %w", err)
		}

		// Find matching end marker after the begin comment.
		rest := content[beginCommentEnd:]
		endIdx = strings.Index(rest, endMarker)
		if endIdx == -1 {
			return nil, fmt.Errorf("missing %s for begin comment", endMarker)
		}

		// Check for nested begin markers between this begin and the end.
		between := rest[:endIdx]
		if strings.Contains(between, beginMarker) {
			return nil, errors.New("nested <!-- docfresh begin found before <!-- docfresh end -->")
		}

		endCommentEnd := beginCommentEnd + endIdx + len(endMarker)
		endComment := content[beginCommentEnd+endIdx : endCommentEnd]

		blocks = append(blocks, &Block{
			Type:         "block",
			Input:        &input,
			BeginComment: beginComment,
			EndComment:   endComment,
		})

		pos = endCommentEnd
	}
	return blocks, nil
}

func appendText(blocks []*Block, text string) []*Block {
	if text == "" {
		return blocks
	}
	return append(blocks, &Block{
		Type:    "text",
		Content: text,
	})
}
