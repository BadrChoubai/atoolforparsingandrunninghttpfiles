package parser

import (
	"net/http"
	"net/url"
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

func TestHttpFileParser_BuildRequests(t *testing.T) {
	for _, tc := range []struct {
		name         string
		scannedLines [][]string
		expected     []*HTTPRequest
	}{
		{
			name: "single request",
			scannedLines: [][]string{
				{"Get user", "GET http://example.com/users/1"},
			},
			expected: []*HTTPRequest{
				{
					Description: "Get user",
					Request: &http.Request{
						Method: "GET",
						URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/users/1"},
					},
				},
			},
		},
		{
			name: "multiple requests",
			scannedLines: [][]string{
				{"List users", "GET http://example.com/users"},
				{"Create user", "POST http://example.com/users"},
			},
			expected: []*HTTPRequest{
				{
					Description: "List users",
					Request: &http.Request{
						Method: "GET",
						URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/users"},
					},
				},
				{
					Description: "Create user",
					Request: &http.Request{
						Method: "POST",
						URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/users"},
					},
				},
			},
		}} {
		t.Run(tc.name, func(t *testing.T) {
			parser := &HttpFileParser{
				ScannedLines: tc.scannedLines,
			}

			err := parser.BuildRequests()
			if err != nil {
				t.Fatalf("BuildRequests() returned unexpected error: %v", err)
			}

			if len(parser.Requests) != len(tc.expected) {
				t.Fatalf("Expected %d requests, got %d", len(tc.expected), len(parser.Requests))
			}

			for i, req := range parser.Requests {
				exp := tc.expected[i]
				if req.Description != exp.Description {
					t.Errorf("Request %d: expected description %q, got %q", i, exp.Description, req.Description)
				}
				if req.Request.Method != exp.Request.Method {
					t.Errorf("Request %d: expected method %q, got %q", i, exp.Request.Method, req.Request.Method)
				}
				if req.Request.URL.String() != exp.Request.URL.String() {
					t.Errorf("Request %d: expected URL %q, got %q", i, exp.Request.URL.String(), req.Request.URL.String())
				}
			}
		})
	}
}
