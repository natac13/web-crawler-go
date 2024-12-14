package main

import "testing"

func TestNormalizeUrl(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "remove scheme",
			inputURL:    "http://example.com",
			expectedURL: "example.com",
		},
		{
			name:        "remove scheme and www preserve path",
			inputURL:    "http://www.example.com/path",
			expectedURL: "example.com/path",
		},
		{
			name:        "remove https scheme",
			inputURL:    "https://example.com",
			expectedURL: "example.com",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := normalizeURL(tt.inputURL)
			if err != nil {
				t.Errorf("text %v - '%s' FAIL: unexpected error: %v", i, tt.name, err)
			}
			if actual != tt.expectedURL {
				t.Errorf("test %v - '%s' FAIL: expected '%s' but got '%s'", i, tt.name, tt.expectedURL, actual)
			}
		})
	}
}
