// Package parser contains code for parsing the `.http_file` file format
package parser

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var _ Parser = (*HTTPFile)(nil)

type Parser interface {
	Parse(filepath string) (bool, error)
}

type HTTPFile struct {
	Requests     []*HTTPRequest
	ScannedLines [][]string
}

type HTTPRequest struct {
	Description string
	*http.Request
}

func NewHttpFileParser() *HTTPFile {
	return &HTTPFile{}
}

// Parse scans a given `.http_file` file and appends valid results to `HttpFileParser.requests`
// If the list of requests was built, Parse returns `true`, otherwise `false` and a resulting error
func (h *HTTPFile) Parse(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	var scannedLine []string
	scanning := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Start of a new request block
		if strings.HasPrefix(line, "###") {
			if scanning && len(scannedLine) > 0 {
				h.ScannedLines = append(h.ScannedLines, scannedLine)
			}

			desc := strings.TrimSpace(strings.TrimPrefix(line, "###"))
			scannedLine = []string{desc}
			scanning = true
			continue
		}

		if scanning {
			scannedLine = append(scannedLine, line)
		}
	}

	// Add final block if any
	if scanning && len(scannedLine) > 0 {
		h.ScannedLines = append(h.ScannedLines, scannedLine)
	}

	return len(h.ScannedLines) > 0, nil
}

func (h *HTTPFile) BuildRequests() error {
	for _, block := range h.ScannedLines {
		req, err := parseRequestBlock(block)
		if err != nil {
			return err
		}
		h.Requests = append(h.Requests, req)
	}
	return nil
}

func parseRequestBlock(lines []string) (*HTTPRequest, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("incomplete request block")
	}

	desc := lines[0]
	requestLine := strings.TrimSpace(lines[1])
	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		return nil, fmt.Errorf("malformed request line: %q", requestLine)
	}

	method := parts[0]
	rawURL := parts[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %q: %w", rawURL, err)
	}

	req := &http.Request{
		Method: method,
		URL:    parsedURL,
		Header: make(http.Header),
	}

	var bodyLines []string
	foundBlank := false
	for _, line := range lines[2:] {
		line = strings.TrimSpace(line)
		if line == "" {
			foundBlank = true
			continue
		}

		if !foundBlank {
			kv := strings.SplitN(line, ":", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("malformed header line: %q", line)
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			req.Header.Add(key, value)
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	if len(bodyLines) > 0 {
		body := strings.Join(bodyLines, "\n")
		req.Body = io.NopCloser(strings.NewReader(body))
	}

	return &HTTPRequest{
		Description: desc,
		Request:     req,
	}, nil
}
