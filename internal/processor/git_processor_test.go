package processor

import (
	"testing"

	"github.com/carlosarraes/shush/internal/types"
)

func TestRemoveCommentsFromLine(t *testing.T) {

	cli := types.CLI{}
	p := &Processor{cli: cli}


	jsLanguage := types.Language{
LineComment: "
		BlockComment: &types.BlockComment{
Start: "
			End:   "*/",
		},
	}


	pyLanguage := types.Language{
		LineComment: "#",
	}

	tests := []struct {
		name     string
		line     string
		language types.Language
		cli      types.CLI
		expected string
	}{
		{
			name:     "line comment removal - JavaScript",
line:     "console.log('hello');
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "console.log('hello');",
		},
		{
			name:     "line comment removal - Python",
			line:     "print('hello')  # This is a comment",
			language: pyLanguage,
			cli:      types.CLI{},
			expected: "print('hello')",
		},
		{
			name:     "block comment removal - single line",
line:     "var x = 5;  var y = 10;",
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "var x = 5;  var y = 10;",
		},
		{
			name:     "only line comment when inline flag set",
line:     "code();
			language: jsLanguage,
			cli:      types.CLI{Inline: true},
			expected: "code();",
		},
		{
			name:     "only block comment when block flag set",
line:     "code();
			language: jsLanguage,
			cli:      types.CLI{Block: true},
expected: "code();
		},
		{
			name:     "no comments to remove",
			line:     "var x = 5;",
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "var x = 5;",
		},
		{
			name:     "comment-only line",
line:     "
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "",
		},
		{
			name:     "multiple block comments",
line:     " code ",
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.cli = tt.cli
			result := p.removeCommentsFromLine(tt.line, tt.language)
			if result != tt.expected {
				t.Errorf("removeCommentsFromLine() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRemoveCommentsFromLineEdgeCases(t *testing.T) {
	cli := types.CLI{}
	p := &Processor{cli: cli}

	jsLanguage := types.Language{
LineComment: "
		BlockComment: &types.BlockComment{
Start: "
			End:   "*/",
		},
	}

	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "empty line",
			line:     "",
			expected: "",
		},
		{
			name:     "whitespace only",
			line:     "   \t  ",
			expected: "   \t  ",
		},
		{
			name:     "comment at start",
line:     "
			expected: "",
		},
		{
			name:     "comment with special characters",
line:     "code();
			expected: "code();",
		},
		{
			name:     "comment after code",
line:     `console.log("hello");
			expected: `console.log("hello");`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.removeCommentsFromLine(tt.line, jsLanguage)
			if result != tt.expected {
				t.Errorf("removeCommentsFromLine() = %q, want %q", result, tt.expected)
			}
		})
	}
}
