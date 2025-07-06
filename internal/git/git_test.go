package git

import (
	"testing"
)

func TestParseGitDiffUnified(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []LineRange
	}{
		{
			name: "single line change",
			input: `@@ -10,1 +10,1 @@
-old line
+new line`,
			expected: []LineRange{{Start: 10, End: 10}},
		},
		{
			name: "multiple line addition",
			input: `@@ -5,0 +5,3 @@
+line 1
+line 2
+line 3`,
			expected: []LineRange{{Start: 5, End: 7}},
		},
		{
			name: "multiple hunks",
			input: `@@ -10,1 +10,1 @@
-old line
+new line
@@ -20,0 +20,2 @@
+added line 1
+added line 2`,
			expected: []LineRange{
				{Start: 10, End: 10},
				{Start: 20, End: 21},
			},
		},
		{
			name: "deletion only",
			input: `@@ -10,2 +10,0 @@
-deleted line 1
-deleted line 2`,
			expected: []LineRange{},
		},
		{
			name: "no count specified (single line)",
			input: `@@ -10 +10 @@
-old line
+new line`,
			expected: []LineRange{{Start: 10, End: 10}},
		},
		{
			name:     "empty diff",
			input:    "",
			expected: []LineRange{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDiffUnified(tt.input)
			if err != nil {
				t.Fatalf("ParseDiffUnified() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("ParseDiffUnified() = %v, want %v", result, tt.expected)
			}

			for i, r := range result {
				if r.Start != tt.expected[i].Start || r.End != tt.expected[i].End {
					t.Errorf("ParseDiffUnified()[%d] = %v, want %v", i, r, tt.expected[i])
				}
			}
		})
	}
}

func TestIsInLineRanges(t *testing.T) {
	ranges := []LineRange{
		{Start: 5, End: 7},
		{Start: 10, End: 12},
		{Start: 20, End: 20},
	}

	tests := []struct {
		lineNum  int
		expected bool
	}{
		{1, false},
		{4, false},
		{5, true},
		{6, true},
		{7, true},
		{8, false},
		{10, true},
		{11, true},
		{12, true},
		{15, false},
		{20, true},
		{21, false},
		{100, false},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.lineNum)), func(t *testing.T) {
			result := IsInLineRanges(tt.lineNum, ranges)
			if result != tt.expected {
				t.Errorf("IsInLineRanges(%d) = %v, want %v", tt.lineNum, result, tt.expected)
			}
		})
	}
}

func TestIsInLineRangesEmptyRanges(t *testing.T) {
	var emptyRanges []LineRange

	testLines := []int{1, 5, 10, 100}
	for _, lineNum := range testLines {
		result := IsInLineRanges(lineNum, emptyRanges)
		if result != false {
			t.Errorf("IsInLineRanges(%d, emptyRanges) = %v, want false", lineNum, result)
		}
	}
}
