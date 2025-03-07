// The GraphQL client

package graphql

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient creates a new client with the given endpoint and options.
func NewClient(endpoint string) *Client {
	// Create default client
	client := &Client{
		endpoint:   endpoint,
		httpClient: http.DefaultClient,
	}

	// Add cookies support
	jar, _ := cookiejar.New(nil)
	client.httpClient.Jar = jar

	return client
}

// ***************************************************** //
//                     HELPER METHODS                    //
// ***************************************************** //

// ClearCookies clears all cookies from the client's cookie jar.
// This method panics on clients that have no cookie jar.
func (c *Client) ClearCookies() {
	jar, _ := cookiejar.New(nil)
	c.httpClient.Jar = jar
}

// AddCookie adds a simple cookie to the client's cookie jar.
func (c *Client) AddCookie(endpoint, name, value string) {
	endpointUrl, _ := url.Parse(endpoint)
	c.httpClient.Jar.SetCookies(endpointUrl, []*http.Cookie{
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

	// Create the request
	var httpRequest *http.Request
	var err error
	if len(request.files) > 0 {
		httpRequest, err = c.requestMultipart(request)
	} else {
		httpRequest, err = c.requestJson(request)
	}

	if err != nil {
		return err
	}

	// Do the request
	httpRequest = httpRequest.WithContext(ctx)
	res, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Process the body
	data, err := unmarshal(res.Body, response)

	if err != nil {
		return err
	}

	if len(data.Errors) > 0 {
		return data.Errors[0]
	}

	return nil
}
