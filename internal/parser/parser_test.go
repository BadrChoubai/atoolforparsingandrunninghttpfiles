package parser

import (
	"os"
	"testing"
)

func TestHttpFileParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []*HTTPRequest
	}{
		{
			name: "single request",
			content: `
### Get user
GET /users/1
`,
			expected: []*HTTPRequest{
				{
					Description: "Get user",
					Method:      "GET",
					URL:         "/users/1",
				},
			},
		},
		{
			name: "multiple requests",
			content: `
### Get users
GET /users

### Create user
POST /users
`,
			expected: []*HTTPRequest{
				{
					Description: "Get users",
					Method:      "GET",
					URL:         "/users",
				},
				{
					Description: "Create user",
					Method:      "POST",
					URL:         "/users",
				},
			},
		},
		{
			name:     "empty file",
			content:  ``,
			expected: []*HTTPRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write content to a temporary file
			tmpFile, err := os.CreateTemp("", "*.http")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			parser := &HttpFileParser{}
			ok, err := parser.Parse(tmpFile.Name())
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if !ok && len(tt.expected) > 0 {
				t.Errorf("Expected parse to succeed, got false")
			}
			if len(parser.Requests) != len(tt.expected) {
				t.Fatalf("Expected %d requests, got %d", len(tt.expected), len(parser.Requests))
			}
			for i, req := range parser.Requests {
				exp := tt.expected[i]
				if req.Description != exp.Description || req.Method != exp.Method || req.URL != exp.URL {
					t.Errorf("Request %d mismatch: got %+v, want %+v", i, req, exp)
				}
			}
		})
	}
}
