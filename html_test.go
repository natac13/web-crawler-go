package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{
				"https://example.com/path/one",
				"https://other.com/path/one",
			},
		},
		{
			name:     "no URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<span>Boot.dev</span>
	</body>
</html>
`,
			expected: []string{},
		},
		{
			name:     "many URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/two">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/three">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{
				"https://example.com/path/one",
				"https://other.com/path/one",
				"https://other.com/path/two",
				"https://other.com/path/three",
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tt.inputBody, tt.inputURL)
			if err != nil {
				t.Errorf("Test %d: getURLsFromHTML() error = %v", i, err)
			}
			if len(actual) != len(tt.expected) {
				t.Errorf("Test %d: getURLsFromHTML() = %v, want %v", i, actual, tt.expected)
			}
			if len(actual) == 0 {
				return
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Test %d: getURLsFromHTML() = %v, want %v", i, actual, tt.expected)
			}
		})
	}
}
