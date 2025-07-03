package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/carlosarraes/shush/internal/types"
)

var languageMap = map[string]types.Language{
	"lua":  {LineComment: "--"},
	"py":   {LineComment: "#"},
	"sh":   {LineComment: "#"},
	"js":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"ts":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"go":   {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"c":    {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"cpp":  {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"java": {LineComment: "//", BlockComment: &types.BlockComment{Start: "/*", End: "*/"}},
	"rb":   {LineComment: "#"},
	"pl":   {LineComment: "#"},
	"yml":  {LineComment: "#"},
	"yaml": {LineComment: "#"},
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
		"js":   "JavaScript",
		"ts":   "TypeScript",
		"go":   "Go",
		"c":    "C",
		"cpp":  "C++",
		"java": "Java",
		"rb":   "Ruby",
		"pl":   "Perl",
		"yml":  "YAML",
		"yaml": "YAML",
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
	return ok
}