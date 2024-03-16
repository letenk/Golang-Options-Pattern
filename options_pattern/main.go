package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client represents out HTTP client.
type Client struct {
	client          *http.Client
	timeout         time.Duration
	userAgent       string
	followRedirects bool
}

// Option is a functional option type that allows us to configure
type Option func(*Client)

// NewClient creates a new HTTP client with default options.
func NewClient(options ...Option) *Client {
	client := &Client{
		client:          &http.Client{},
		timeout:         30 * time.Second, // Default timeout
		userAgent:       "My HTTP Client", // Default user agent
		followRedirects: true,             // Default follows redirects
	}

	// Apply all the functional options to configure the client
	for _, opt := range options {
		opt(client)
	}

	return client
}

// WithTimeout is a functional option to set the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithUserAgent is a functional option to set the HTTP client user agent.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithoutRedirects is a functional option to disable following redirects.
func WithoutRedirects() Option {
	return func(c *Client) {
		c.followRedirects = false
	}
}

// UseInsecureTransport is a functional option to use an insecure HTTP transport.
func UseInsecureTransport() Option {
	return func(c *Client) {
		c.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

// Get performs an HTTP GET request.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)

	c.client.Timeout = c.timeout
	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if !c.followRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}

	return c.client.Do(req)
}

// Example usage:
func main() {
	// Create a new HTTP client with custom options.
	client := NewClient(
		WithTimeout(10*time.Second),
		WithUserAgent("My Custom User Agent"),
		UseInsecureTransport(),
	)

	// Use the client to make an HTTP request.
	response, err := client.Get("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Print the response body
	fmt.Println(string(body))
}
