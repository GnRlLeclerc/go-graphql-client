// Multipat content type requests

package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
)

func (c *Client) requestMultipart(request *Request) (*http.Request, error) {
	var requestBody bytes.Buffer
	fileMap := make(map[string][]string)
	writer := multipart.NewWriter(&requestBody)

	// Add file null variables to the request for each file to be uploaded
	filecount := 0
	for fieldname, files := range request.files {
		if len(files) > 1 {
			// Slice of nil
			request.Var(fieldname, make([]interface{}, len(files)))
			for i := range files {
				fileMap[strconv.Itoa(filecount)] = []string{fmt.Sprintf("variables.%s.%d", fieldname, i)}
				filecount += 1
			}
		} else {
			request.Var(fieldname, nil)
			fileMap[strconv.Itoa(filecount)] = []string{fmt.Sprintf("variables.%s", fieldname)}
			filecount += 1
		}
	}

	query, err := request.toJSON()
	if err != nil {
		return nil, fmt.Errorf("Error marshaling the graphql request to json: %v", err)
	}

	// Write the GraphQL query
	writer.WriteField("operations", string(query))

	// Write the file map
	res, err := json.Marshal(fileMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling the file map: %v", err)
	}
	writer.WriteField("map", string(res))

	// Write the files
	filecount = 0
	for _, files := range request.files {
		for _, file := range files {
			// Get the file mime content type
			ext := filepath.Ext(file.filename)
			contentType := mime.TypeByExtension(ext)
			if contentType == "" {
				contentType = "application/octet-stream" // Fallback for unknown types
			}

			// File multipart header
			partHeader := make(textproto.MIMEHeader)
			partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%d"; filename="%s"`, filecount, file.filename))
			partHeader.Set("Content-Type", contentType)

			w, err := writer.CreatePart(partHeader)
			if err != nil {
				return nil, fmt.Errorf("Error creating form file: %v", err)
			}

			fileBytes, err := io.ReadAll(file.reader)
			if err != nil {
				return nil, fmt.Errorf("Error reading file %s: %v", file.filename, err)
			}

			w.Write(fileBytes)

			filecount += 1
		}
	}

	// Write the closing boundary
	writer.Close()

	// Form the http request
	r, err := http.NewRequest(http.MethodPost, c.endpoint, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating graphql request: %v", err)
	}

	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Header.Set("Accept", "application/json; charset=utf-8")

	for key, values := range request.header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	return r, nil
}
