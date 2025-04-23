// Package parser contains code for parsing the `.http` file format
package parser

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var _ Parser = (*HttpFileParser)(nil)

var (
	Methods = map[string]string{
		"GET":    http.MethodGet,
		"PUT":    http.MethodPut,
		"POST":   http.MethodPost,
		"DELETE": http.MethodDelete,
	}
)

type Parser interface {
	Parse(filepath string) (bool, error)
}

type HTTPRequest struct {
	Description string
	*http.Request
}

type HttpFileParser struct {
	Requests     []*HTTPRequest
	ScannedLines [][]string
}

// Parse scans a given `.http` file and appends valid results to `HttpFileParser.requests`
// If the list of requests was built, Parse returns `true`, otherwise `false` and a resulting error
func (h *HttpFileParser) Parse(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if err != nil {
		return false, err
	}

	var scannedLine []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else {
			if prefix, reset := strings.CutPrefix(line, "###"); !reset {
			} else {
				if len(scannedLine) > 0 {
					h.ScannedLines = append(h.ScannedLines, scannedLine)
				}
				scannedLine = []string{strings.TrimSpace(prefix)}
				continue
			}
			scannedLine = append(scannedLine, strings.TrimSpace(line))
		}
	}

	if len(scannedLine) > 0 {
		h.ScannedLines = append(h.ScannedLines, scannedLine)
	}

	if len(h.ScannedLines) > 0 {
		return true, nil
	} else {
		return false, err
	}
}

func (h *HttpFileParser) BuildRequests() error {
	for _, lines := range h.ScannedLines {
		if len(lines) < 2 {
			continue
		}

		desc := lines[0]
		line := strings.TrimSpace(lines[1])
		parts := strings.Fields(line)

		if len(parts) < 2 {
			return fmt.Errorf("malformed request line: %q", line)
		}

		method := parts[0]
		rawURL := parts[1]
		parsedURL, err := url.ParseRequestURI(rawURL)
		if err != nil {
			return fmt.Errorf("failed to parse UR: %q: %w", rawURL, err)
		}

		req := &http.Request{
			Method: method,
			URL:    parsedURL,
			Header: make(http.Header),
		}

		for _, headerLine := range lines[2:] {
			headerLine = strings.TrimSpace(headerLine)
			if headerLine == "" {
				continue
			}

			kv := strings.SplitN(headerLine, ":", 2)
			if len(kv) != 2 {
				return fmt.Errorf("malformed header line: %q", headerLine)
			}

			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			req.Header.Add(key, value)
		}

		h.Requests = append(h.Requests, &HTTPRequest{
			Description: desc,
			Request:     req,
		})
	}

	return nil
}
