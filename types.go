package oneview

import (
	"encoding/json"
	"strconv"
	"strings"
)

// FlexString unmarshals JSON strings or numbers into a string.
// OneView/Global Dashboard sometimes return eTag and version fields as either.
type FlexString string

func (s FlexString) String() string { return string(s) }

func (s *FlexString) UnmarshalJSON(b []byte) error {
	b = bytesTrim(b)
	if len(b) == 0 || string(b) == "null" {
		*s = ""
		return nil
	}
	if b[0] == '"' {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		*s = FlexString(v)
		return nil
	}
	*s = FlexString(b)
	return nil
}

func (s FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func bytesTrim(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

// FlexInt unmarshals JSON numbers or numeric strings into an int.
type FlexInt int

func (n FlexInt) Int() int { return int(n) }

func (n *FlexInt) UnmarshalJSON(b []byte) error {
	b = bytesTrim(b)
	if len(b) == 0 || string(b) == "null" {
		*n = 0
		return nil
	}
	if b[0] == '"' {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		if v == "" {
			*n = 0
			return nil
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			f, err2 := strconv.ParseFloat(v, 64)
			if err2 != nil {
				return err
			}
			i = int(f)
		}
		*n = FlexInt(i)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(b, &num); err != nil {
		return err
	}
	i, err := num.Int64()
	if err != nil {
		f, err2 := num.Float64()
		if err2 != nil {
			return err
		}
		i = int64(f)
	}
	*n = FlexInt(i)
	return nil
}

func (n FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(n))
}

// FlexBool unmarshals JSON booleans or "true"/"false" strings.
type FlexBool bool

func (b FlexBool) Bool() bool { return bool(b) }

func (b *FlexBool) UnmarshalJSON(raw []byte) error {
	raw = bytesTrim(raw)
	if len(raw) == 0 || string(raw) == "null" {
		*b = false
		return nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return err
		}
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "", "false", "0", "no", "off":
			*b = false
		default:
			*b = true
		}
		return nil
	}
	var v bool
	if err := json.Unmarshal(raw, &v); err == nil {
		*b = FlexBool(v)
		return nil
	}
	var n FlexInt
	if err := n.UnmarshalJSON(raw); err != nil {
		return err
	}
	*b = n.Int() != 0
	return nil
}

func (b FlexBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}

// VersionInfo is GET /rest/version.
type VersionInfo struct {
	MinimumVersion FlexInt `json:"minimumVersion"`
	CurrentVersion FlexInt `json:"currentVersion"`
}

// Resource is the common OneView resource envelope (type, uri, eTag, status, …)
// plus Global Dashboard aggregation fields (id, applianceLocation, originalUri).
type Resource struct {
	Type              string     `json:"type,omitempty"`
	URI               string     `json:"uri,omitempty"`
	Category          string     `json:"category,omitempty"`
	Name              string     `json:"name,omitempty"`
	Description       string     `json:"description,omitempty"`
	Status            string     `json:"status,omitempty"`
	State             string     `json:"state,omitempty"`
	StateReason       string     `json:"stateReason,omitempty"`
	ETag              FlexString `json:"eTag,omitempty"`
	Created           string     `json:"created,omitempty"`
	Modified          string     `json:"modified,omitempty"`
	ID                string     `json:"id,omitempty"`
	UUID              string     `json:"uuid,omitempty"`
	ApplianceLocation string     `json:"applianceLocation,omitempty"`
	ApplianceName     string     `json:"applianceName,omitempty"`
	ApplianceURI      string     `json:"appluri,omitempty"`
	OriginalURI       string     `json:"originalUri,omitempty"`
	Groups            []GroupRef `json:"groups,omitempty"`
}

// GroupRef is a Global Dashboard group membership pointer.
type GroupRef struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// Collection is a paginated OneView list response.
type Collection[T any] struct {
	Type        string `json:"type,omitempty"`
	Category    string `json:"category,omitempty"`
	URI         string `json:"uri,omitempty"`
	Start       int    `json:"start"`
	Count       int    `json:"count"`
	Total       int    `json:"total"`
	NextPageURI string `json:"nextPageUri,omitempty"`
	PrevPageURI string `json:"prevPageUri,omitempty"`
	Members     []T    `json:"members"`
}

// AssociatedResource is used by alerts and tasks.
type AssociatedResource struct {
	AssociationType  string `json:"associationType,omitempty"`
	ResourceCategory string `json:"resourceCategory,omitempty"`
	ResourceName     string `json:"resourceName,omitempty"`
	ResourceURI      string `json:"resourceUri,omitempty"`
}

// NamedURI is a {name, uri} pair used throughout the API.
type NamedURI struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// IDFromURI returns the last path segment of a resource URI.
func IDFromURI(uri string) string {
	uri = strings.TrimRight(uri, "/")
	if i := strings.LastIndex(uri, "/"); i >= 0 && i+1 < len(uri) {
		return uri[i+1:]
	}
	return uri
}

// JSONList unmarshals a JSON array or a single object into a slice.
// OneView sometimes returns LogicalDrives/PhysicalDrives as an object, sometimes as an array.
type JSONList[T any] []T

func (l *JSONList[T]) UnmarshalJSON(b []byte) error {
	b = bytesTrim(b)
	if len(b) == 0 || string(b) == "null" {
		*l = nil
		return nil
	}
	if b[0] == '[' {
		var s []T
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*l = s
		return nil
	}
	if b[0] == '{' && string(b) == "{}" {
		*l = nil
		return nil
	}
	var one T
	if err := json.Unmarshal(b, &one); err != nil {
		return err
	}
	*l = []T{one}
	return nil
}
