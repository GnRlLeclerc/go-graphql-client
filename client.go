// The GraphQL client

package graphql

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

type ContentType string

const (
	ContentTypeJSON      = "application/json"
	ContentTypeMultipart = "multipart/form-data"
)

type Client struct {
	endpoint    string
	contentType ContentType
	httpClient  *http.Client
}

// NewClient creates a new client with the given endpoint and options.
func NewClient(endpoint string, opts ...ClientOption) *Client {
	// Create default client
	client := &Client{
		endpoint:    endpoint,
		contentType: ContentTypeJSON,
		httpClient:  http.DefaultClient,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// ***************************************************** //
//                 CONFIGURATION OPTIONS                 //
// ***************************************************** //

// ClientOption is a function that updates a client
type ClientOption func(*Client)

// UseMultipart sets the client to use multipart content type instead of the default JSON one.
func UseMultipart() ClientOption {
	return func(c *Client) {
		c.contentType = ContentTypeMultipart
	}
}

func UseCookies() ClientOption {
	jar, _ := cookiejar.New(nil)
	return func(c *Client) {
		c.httpClient.Jar = jar
	}
}

// ***************************************************** //
//                     HELPER METHODS                    //
// ***************************************************** //

// ClearCookies clears all cookies from the client's cookie jar.
// This method panics on clients that have no cookie jar.
func (c *Client) ClearCookies() {
	if c.httpClient.Jar == nil {
		panic("Called ClearCookies on a client that does not have a cookie jar")
	}

	jar, _ := cookiejar.New(nil)
	c.httpClient.Jar = jar
}

// AddCookie adds a simple cookie to the client's cookie jar.
func (c *Client) AddCookie(name, value string) {
	if c.httpClient.Jar == nil {
		panic("Called AddCookie on a client that does not have a cookie jar")
	}

	c.httpClient.Jar.SetCookies(nil, []*http.Cookie{
		{
			Name:  name,
			Value: value,
		},
	})
}

// Run runs a graphql request and attempts to unmarshal the response into the given interface.
func (c *Client) Run(ctx context.Context, request *Request, response interface{}) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if len(request.files) > 0 && c.contentType != ContentTypeMultipart {
		return fmt.Errorf("Cannot process file uploads with content type %s. Use a multipart content type client instead.", c.contentType)
	}

	// TODO: handle both content types

	return nil
}
