// GraphQL requests

package graphql

import (
	"io"
	"net/http"
)

// Request is a GraphQL request
type Request struct {
	query     string
	variables map[string]interface{}
	files     []file // File uploads (multipart only)
	header    http.Header
}

// file represents a file to be uploaded
type file struct {
	fieldname string
	filename  string
	reader    io.Reader
}

// NewRequest creates a new GraphQL request with the given query
func NewRequest(query string) *Request {
	return &Request{
		query:     query,
		variables: make(map[string]interface{}),
		header:    make(http.Header),
	}
}

// ***************************************************** //
//                     HELPER METHODS                    //
// ***************************************************** //

// Var sets a variable for the request.
func (r *Request) Var(fieldname string, value interface{}) {
	r.variables[fieldname] = value
}

// Header sets a header for the request.
func (r *Request) Header(key, value string) {
	r.header.Set(key, value)
}

// File sets a file upload for the request.
// A request with files should be sent using a multipart content type.
func (r *Request) File(fieldname, filename string, reader io.Reader) {
	r.files = append(r.files, file{fieldname, filename, reader})
}
