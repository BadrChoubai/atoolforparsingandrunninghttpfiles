package parser

import (
	"os"
	"testing"
)

func TestHttpFileParser_Parse(t *testing.T) {
	for _, tc := range []struct {
		name                    string
		content                 string
		requestsParsedFromInput int
	}{
		{
			name: "single request",
			content: `
### Get user
GET /users/1
`,
			requestsParsedFromInput: 1,
		},
		{
			name: "multiple requests",
			content: `
### Get users
GET /users

### Create user
POST /users
`,
			requestsParsedFromInput: 2,
		},
		{
			name:                    "empty file",
			content:                 ``,
			requestsParsedFromInput: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Write content to a temporary file
			tmpFile, err := os.CreateTemp("", "*.http")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tc.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			parser := &HttpFileParser{}
			ok, err := parser.Parse(tmpFile.Name())
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if !ok && len(parser.ScannedLines) > 0 {
				t.Errorf("Expected parse to succeed, got false")
			}
			if len(parser.ScannedLines) != tc.requestsParsedFromInput {
				t.Fatalf("Expected %d requests, got %d", tc.requestsParsedFromInput, len(parser.Requests))
			}
		})
	}
}
