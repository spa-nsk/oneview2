package oneview

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Response is a raw HTTP response from the appliance.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// Location returns the Location header (task URI on 202 Accepted).
func (r *Response) Location() string {
	if r == nil {
		return ""
	}
	return r.Header.Get("Location")
}

// TaskURI returns a task URI from Location or the JSON body.
func (r *Response) TaskURI() string {
	if loc := r.Location(); loc != "" {
		if i := strings.Index(loc, "/rest/"); i >= 0 {
			return loc[i:]
		}
		return loc
	}
	if r == nil || len(r.Body) == 0 {
		return ""
	}
	var t struct {
		URI string `json:"uri"`
	}
	if json.Unmarshal(r.Body, &t) == nil && strings.Contains(t.URI, "/rest/tasks") {
		return t.URI
	}
	return ""
}

func (c *Client) newRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		path = c.BaseURL() + path
	}
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, path, rdr)
	if err != nil {
		return nil, err
	}
	c.mu.RLock()
	api := c.apiVersion
	token := c.authToken
	ifMatch := c.ifMatch
	ua := c.userAgent
	c.mu.RUnlock()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", ua)
	if api > 0 {
		req.Header.Set("X-API-Version", strconv.Itoa(api))
	}
	if token != "" {
		req.Header.Set("Auth", token)
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}
	if method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete {
		req.Header.Set("If-Match", ifMatch)
	}
	return req, nil
}

// Do sends a request. Path may be absolute or appliance-relative (/rest/...).
func (c *Client) Do(ctx context.Context, method, path string, body []byte) (*Response, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oneview: %s %s: %w", method, path, err)
	}
	defer httpResp.Body.Close()
	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("oneview: read body: %w", err)
	}
	resp := &Response{StatusCode: httpResp.StatusCode, Header: httpResp.Header.Clone(), Body: raw}
	if httpResp.StatusCode >= 400 {
		return resp, parseAPIError(method, path, httpResp.StatusCode, raw)
	}
	return resp, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, in, out any) (*Response, error) {
	var body []byte
	var err error
	if in != nil {
		body, err = json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("oneview: marshal request: %w", err)
		}
	}
	resp, err := c.Do(ctx, method, path, body)
	if err != nil {
		return resp, err
	}
	if out != nil && len(resp.Body) > 0 && resp.StatusCode != http.StatusNoContent {
		if err := decodeJSON(resp.Body, out); err != nil {
			return resp, fmt.Errorf("oneview: decode %s %s: %w", method, path, err)
		}
	}
	return resp, nil
}

// GetJSON performs GET and unmarshals JSON into out.
func (c *Client) GetJSON(ctx context.Context, path string, out any) error {
	_, err := c.doJSON(ctx, http.MethodGet, path, nil, out)
	return err
}

// PostJSON performs POST and unmarshals JSON into out.
func (c *Client) PostJSON(ctx context.Context, path string, in, out any) (*Response, error) {
	return c.doJSON(ctx, http.MethodPost, path, in, out)
}

// PutJSON performs PUT and unmarshals JSON into out.
func (c *Client) PutJSON(ctx context.Context, path string, in, out any) (*Response, error) {
	return c.doJSON(ctx, http.MethodPut, path, in, out)
}

// PatchJSON performs PATCH (JSON Patch or OneView patch array) and unmarshals into out.
func (c *Client) PatchJSON(ctx context.Context, path string, in, out any) (*Response, error) {
	return c.doJSON(ctx, http.MethodPatch, path, in, out)
}

// DeleteJSON performs DELETE.
func (c *Client) DeleteJSON(ctx context.Context, path string, out any) (*Response, error) {
	return c.doJSON(ctx, http.MethodDelete, path, nil, out)
}

// GetAll follows nextPageUri until the collection is exhausted.
func GetAll[T any](ctx context.Context, c *Client, path string, opts ListOptions) (*Collection[T], error) {
	if opts.Count == 0 {
		opts.Count = -1
	}
	first := &Collection[T]{}
	if err := c.GetJSON(ctx, withQuery(path, opts), first); err != nil {
		return nil, err
	}
	for first.NextPageURI != "" {
		next := &Collection[T]{}
		if err := c.GetJSON(ctx, first.NextPageURI, next); err != nil {
			return nil, err
		}
		first.Members = append(first.Members, next.Members...)
		first.Count = len(first.Members)
		first.NextPageURI = next.NextPageURI
		if next.Total > 0 {
			first.Total = next.Total
		}
	}
	return first, nil
}

// GetResource loads any REST resource by URI or id-relative path into out.
func (c *Client) GetResource(ctx context.Context, uri string, out any) error {
	return c.GetJSON(ctx, uri, out)
}

// PatchOp is a JSON Patch / OneView PATCH operation.
type PatchOp struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}
