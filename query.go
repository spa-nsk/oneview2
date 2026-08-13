package oneview

import (
	"net/url"
	"strconv"
	"strings"
)

// ListOptions are the common collection query parameters from the REST API
// (start, count, filter, query, sort, fields, view).
type ListOptions struct {
	Start  int
	Count  int // 0 = appliance default (usually 25); -1 = all items up to page size
	Filter []string
	Query  string
	Sort   string
	Fields string
	View   string
	// GroupURIs is a Global Dashboard-only filter (comma-separated group URIs).
	GroupURIs string
	// Extra is appended as-is for resource-specific query keys.
	Extra url.Values
}

func (o ListOptions) values() url.Values {
	q := url.Values{}
	if o.Start > 0 {
		q.Set("start", strconv.Itoa(o.Start))
	}
	if o.Count != 0 {
		q.Set("count", strconv.Itoa(o.Count))
	}
	for _, f := range o.Filter {
		if f != "" {
			q.Add("filter", f)
		}
	}
	if o.Query != "" {
		q.Set("query", o.Query)
	}
	if o.Sort != "" {
		q.Set("sort", o.Sort)
	}
	if o.Fields != "" {
		q.Set("fields", o.Fields)
	}
	if o.View != "" {
		q.Set("view", o.View)
	}
	if o.GroupURIs != "" {
		q.Set("groupUris", o.GroupURIs)
	}
	for k, vs := range o.Extra {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	return q
}

func joinPath(base string, extra ...string) string {
	p := strings.TrimRight(base, "/")
	for _, e := range extra {
		e = strings.Trim(e, "/")
		if e == "" {
			continue
		}
		p += "/" + e
	}
	return p
}

func withQuery(path string, opts ListOptions) string {
	q := opts.values()
	if len(q) == 0 {
		return path
	}
	return path + "?" + q.Encode()
}
