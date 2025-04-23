// Package parser contains code for parsing the `.http` file format
package parser

import (
	"bufio"
	"net/http"
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
