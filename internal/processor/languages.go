package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/carlosarraes/shush/internal/ignore"
	"github.com/carlosarraes/shush/internal/types"
)

var languageMap = map[string]types.Language{
	"lua":  {LineComment: "--"},
	"py":   {LineComment: "#"},
	"sh":   {LineComment: "#"},
	"bash": {LineComment: "#"},
	"zsh":  {LineComment: "#"},
	"fish": {LineComment: "#"},
	"ps1":  {LineComment: "#"},
	"r":    {LineComment: "#"},

	"js":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"ts":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"jsx":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"tsx":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"go":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"c":     {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"cpp":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"cc":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"cxx":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"h":     {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"hpp":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"java":  {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"cs":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"rs":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"swift": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"kt":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"kts":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"dart":  {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"scala": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"php":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},

	"css":  {BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"scss": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"sass": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"less": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},

	"html": {BlockComment: &types.BlockComment{Start: "<!--", End: "-->"}},
	"htm":  {BlockComment: &types.BlockComment{Start: "<!--", End: "-->"}},
	"xml":  {BlockComment: &types.BlockComment{Start: "<!--", End: "-->"}},
	"svg":  {BlockComment: &types.BlockComment{Start: "<!--", End: "-->"}},

	"rb":   {LineComment: "#"},
	"pl":   {LineComment: "#"},
	"yml":  {LineComment: "#"},
	"yaml": {LineComment: "#"},
	"toml": {LineComment: "#"},
	"ini":  {LineComment: "#", AlternateLineComment: ";"},
	"conf": {LineComment: "#"},
	"cfg":  {LineComment: "#"},

	"sql": {LineComment: "--", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},

	"dockerfile": {LineComment: "#"},
	"makefile":   {LineComment: "#"},
}

func DetectLanguage(filename string) (types.Language, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return types.Language{}, fmt.Errorf("no file extension found")
	}

	ext = strings.TrimPrefix(ext, ".")

	if language, ok := languageMap[ext]; ok {
		return language, nil
	}

	return types.Language{}, fmt.Errorf("unsupported file extension: %s", ext)
}

func GetLanguageName(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".")

	names := map[string]string{
		"lua":  "Lua",
		"py":   "Python",
		"sh":   "Shell",
		"bash": "Bash",
		"zsh":  "Zsh",
		"fish": "Fish",
		"ps1":  "PowerShell",
		"r":    "R",

		"js":    "JavaScript",
		"ts":    "TypeScript",
		"jsx":   "JSX",
		"tsx":   "TSX",
		"go":    "Go",
		"c":     "C",
		"cpp":   "C++",
		"cc":    "C++",
		"cxx":   "C++",
		"h":     "C Header",
		"hpp":   "C++ Header",
		"java":  "Java",
		"cs":    "C#",
		"rs":    "Rust",
		"swift": "Swift",
		"kt":    "Kotlin",
		"kts":   "Kotlin Script",
		"dart":  "Dart",
		"scala": "Scala",
		"php":   "PHP",

		"css":  "CSS",
		"scss": "SCSS",
		"sass": "Sass",
		"less": "Less",

		"html": "HTML",
		"htm":  "HTML",
		"xml":  "XML",
		"svg":  "SVG",

		"rb":   "Ruby",
		"pl":   "Perl",
		"yml":  "YAML",
		"yaml": "YAML",
		"toml": "TOML",
		"ini":  "INI",
		"conf": "Config",
		"cfg":  "Config",

		"sql": "SQL",

		"dockerfile": "Dockerfile",
		"makefile":   "Makefile",
	}

	if name, ok := names[ext]; ok {
		return name
	}
	return ext
}

func IsSupportedFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return false
	}
	ext = strings.TrimPrefix(ext, ".")
	_, ok := languageMap[ext]
	if !ok {
		return false
	}

	return !IsIgnored(filename)
}

func IsIgnored(filename string) bool {
	ignoreChecker, err := ignore.Load()
	if err != nil {
		return false
	}
	return ignoreChecker.IsIgnored(filename)
}
