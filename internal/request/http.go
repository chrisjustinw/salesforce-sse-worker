package request

import "io"

type HTTPRequest struct {
	Path    string
	Headers map[string]string
	Queries map[string]string
	Body    io.Reader
}
