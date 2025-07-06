package processor

import (
	"testing"

	"github.com/carlosarraes/shush/internal/config"
	"github.com/carlosarraes/shush/internal/types"
)

func TestRemoveCommentsFromLine(t *testing.T) {
	cli := types.CLI{}
	p := &Processor{cli: cli}

	jsLanguage := types.Language{
		LineComment: "//",
		BlockComment: &types.BlockComment{
			Start: "/*",
			End:   "*/",
		},
	}

	pyLanguage := types.Language{
		LineComment: "#",
	}

	cfg := config.Default()

	tests := []struct {
		name     string
		line     string
		language types.Language
		cli      types.CLI
		expected string
	}{
		{
			name:     "line comment removal - JavaScript",
			line:     "console.log('hello'); // This is a comment",
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
			line:     "var x = 5; /* comment */ var y = 10;",
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "var x = 5;  var y = 10;",
		},
		{
			name:     "only line comment when inline flag set",
			line:     "code(); /* block */ // line comment",
			language: jsLanguage,
			cli:      types.CLI{Inline: true},
			expected: "code(); /* block */",
		},
		{
			name:     "only block comment when block flag set",
			line:     "code(); /* block */ // line comment",
			language: jsLanguage,
			cli:      types.CLI{Block: true},
			expected: "code();  // line comment",
		},
		{
			name:     "preserve line when preserve lines flag set",
			line:     "// just a comment",
			language: jsLanguage,
			cli:      types.CLI{PreserveLines: true},
			expected: "",
		},
		{
			name:     "no comment",
			line:     "regular code line",
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "regular code line",
		},
		{
			name:     "string literal with comment markers preserved",
			line:     "console.log(\"/* not a comment */\"); // real comment",
			language: jsLanguage,
			cli:      types.CLI{},
			expected: "console.log(\"/* not a comment */\");",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.cli = tt.cli
			result := p.removeCommentsFromLine(tt.line, tt.language, cfg)
			if result != tt.expected {
				t.Errorf("removeCommentsFromLine() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFindCommentIndex(t *testing.T) {
	p := &Processor{}

	tests := []struct {
		name          string
		line          string
		commentMarker string
		expected      int
	}{
		{
			name:          "simple comment",
			line:          "code(); // comment",
			commentMarker: "//",
			expected:      8,
		},
		{
			name:          "comment in string should be ignored",
			line:          "console.log(\"url: http://example.com\"); // real comment",
			commentMarker: "//",
			expected:      40,
		},
		{
			name:          "no comment",
			line:          "regular code",
			commentMarker: "//",
			expected:      -1,
		},
		{
			name:          "comment marker in single quotes",
			line:          "code('//'); // comment",
			commentMarker: "//",
			expected:      12,
		},
		{
			name:          "block comment marker in string",
			line:          "console.log(\"/* not comment */\"); /* real comment */",
			commentMarker: "/*",
			expected:      34,
		},
		{
			name:          "hash comment marker in string - Python",
			line:          "print(\"# not a comment\"); # real comment",
			commentMarker: "#",
			expected:      26,
		},
		{
			name:          "lua comment marker in string",
			line:          "print(\"-- not a comment\"); -- real comment",
			commentMarker: "--",
			expected:      27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.findCommentIndex(tt.line, tt.commentMarker)
			if result != tt.expected {
				t.Errorf("findCommentIndex() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestLineHasComment(t *testing.T) {
	p := &Processor{}

	jsLanguage := types.Language{
		LineComment: "//",
		BlockComment: &types.BlockComment{
			Start: "/*",
			End:   "*/",
		},
	}

	tests := []struct {
		name     string
		line     string
		language types.Language
		expected bool
	}{
		{
			name:     "has line comment",
			line:     "code(); // comment",
			language: jsLanguage,
			expected: true,
		},
		{
			name:     "has block comment start",
			line:     "code(); /* comment",
			language: jsLanguage,
			expected: true,
		},
		{
			name:     "has block comment end",
			line:     "comment */ code();",
			language: jsLanguage,
			expected: true,
		},
		{
			name:     "no comments",
			line:     "regular code line",
			language: jsLanguage,
			expected: false,
		},
		{
			name:     "comment markers in string should not count",
			line:     "console.log(\"/* fake */ and // fake\");",
			language: jsLanguage,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.lineHasComment(tt.line, tt.language)
			if result != tt.expected {
				t.Errorf("lineHasComment() = %t, want %t", result, tt.expected)
			}
		})
	}
}
