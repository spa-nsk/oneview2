package oneview

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIError is an HPE OneView / Global Dashboard error payload.
type APIError struct {
	ErrorCode          string          `json:"errorCode,omitempty"`
	Message            string          `json:"message,omitempty"`
	Details            string          `json:"details,omitempty"`
	ErrorSource        string          `json:"errorSource,omitempty"`
	RecommendedActions []string        `json:"recommendedActions,omitempty"`
	NestedErrors       []APIError      `json:"nestedErrors,omitempty"`
	Data               json.RawMessage `json:"data,omitempty"`
	CanForce           bool            `json:"canForce,omitempty"`
	StatusCode         int             `json:"-"`
	URI                string          `json:"-"`
	Method             string          `json:"-"`
	Body               []byte          `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return "oneview: unknown error"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "oneview: %s %s: HTTP %d", e.Method, e.URI, e.StatusCode)
	if e.ErrorCode != "" {
		fmt.Fprintf(&b, " [%s]", e.ErrorCode)
	}
	if e.Message != "" {
		fmt.Fprintf(&b, ": %s", e.Message)
	}
	if e.Details != "" && e.Details != e.Message {
		fmt.Fprintf(&b, " (%s)", e.Details)
	}
	return b.String()
}

func parseAPIError(method, uri string, status int, body []byte) error {
	ae := &APIError{
		StatusCode: status,
		URI:        uri,
		Method:     method,
		Body:       body,
		Message:    http.StatusText(status),
	}
	if len(body) > 0 {
		_ = json.Unmarshal(body, ae)
		if ae.Message == "" || ae.Message == http.StatusText(status) {
			ae.Message = strings.TrimSpace(string(body))
		}
	}
	return ae
}

// IsNotFound reports whether err is an HTTP 404 from the appliance.
func IsNotFound(err error) bool {
	var ae *APIError
	if asAPIError(err, &ae) {
		return ae.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorized reports whether err is an HTTP 401 from the appliance.
func IsUnauthorized(err error) bool {
	var ae *APIError
	if asAPIError(err, &ae) {
		return ae.StatusCode == http.StatusUnauthorized
	}
	return false
}

func asAPIError(err error, target **APIError) bool {
	if err == nil {
		return false
	}
	ae, ok := err.(*APIError)
	if !ok {
		return false
	}
	*target = ae
	return true
}
