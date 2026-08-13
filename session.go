package oneview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	uriVersion       = "/rest/version"
	uriLoginSessions = "/rest/login-sessions"
)

// Credentials is POST /rest/login-sessions.
type Credentials struct {
	UserName        string `json:"userName"`
	Password        string `json:"password"`
	AuthLoginDomain string `json:"authLoginDomain,omitempty"`
	LoginMsgAck     *bool  `json:"loginMsgAck,omitempty"`
}

// SessionUser is the Global Dashboard login user block.
type SessionUser struct {
	Domain    string   `json:"domain,omitempty"`
	UserName  string   `json:"userName,omitempty"`
	UserRoles []string `json:"user_roles,omitempty"`
}

// Session is the login response. OneView appliances return sessionID;
// Global Dashboard (swagger 300.json) returns token + user.
type Session struct {
	SessionID string      `json:"sessionID,omitempty"`
	Token     string      `json:"token,omitempty"`
	User      SessionUser `json:"user,omitempty"`
}

// AuthHeader returns the value for the Auth header.
func (s Session) AuthHeader() string {
	if s.SessionID != "" {
		return s.SessionID
	}
	return s.Token
}

// GetVersion calls GET /rest/version (no Auth required) and records the
// appliance current/minimum API versions. If Config.APIVersion was 0, the
// client switches to currentVersion.
func (c *Client) GetVersion(ctx context.Context) (*VersionInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL()+uriVersion, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oneview: GET /rest/version: %w", err)
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("oneview: GET /rest/version: %w", err)
	}
	if httpResp.StatusCode >= 400 {
		return nil, parseAPIError(http.MethodGet, uriVersion, httpResp.StatusCode, body)
	}
	var v VersionInfo
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("oneview: decode /rest/version: %w", err)
	}
	min, current := int(v.MinimumVersion), int(v.CurrentVersion)
	requested := c.apiVersion
	if requested == 0 {
		requested = current
	}
	if current > 0 && requested > current {
		requested = current
	}
	if current > 0 && requested < min {
		return &v, fmt.Errorf("oneview: requested API %d is below appliance minimum %d", requested, min)
	}
	c.setVersion(min, current, requested)
	return &v, nil
}

// Login negotiates X-API-Version (if needed) and creates a session.
func (c *Client) Login(ctx context.Context) error {
	if c.username == "" {
		return fmt.Errorf("oneview: Username is required")
	}
	c.mu.RLock()
	needVersion := c.currentVersion == 0
	c.mu.RUnlock()
	if needVersion {
		if _, err := c.GetVersion(ctx); err != nil {
			return err
		}
	}

	paths := []string{uriLoginSessions, uriLoginSessions + "/"}
	ackTrue := true
	acks := []*bool{&ackTrue, nil}
	var last error
	var sess Session
	for _, domain := range loginDomainCandidates(c.domain) {
		for _, ack := range acks {
			body := Credentials{
				UserName:        c.username,
				Password:        c.password,
				AuthLoginDomain: domain,
				LoginMsgAck:     ack,
			}
			for _, path := range paths {
				sess = Session{}
				if _, err := c.PostJSON(ctx, path, body, &sess); err != nil {
					last = err
					if !isRetryableLogin(err) {
						return err
					}
					continue
				}
				if tok := sess.AuthHeader(); tok != "" {
					c.SetAuthToken(tok)
					if domain != "" {
						c.domain = domain
					}
					return nil
				}
				last = fmt.Errorf("oneview: login succeeded but session token is empty")
			}
		}
	}
	if last == nil {
		last = fmt.Errorf("oneview: login failed")
	}
	return last
}

// Logout deletes the current session. It is safe to call more than once.
func (c *Client) Logout(ctx context.Context) error {
	if c.AuthToken() == "" {
		return nil
	}
	_, err := c.DeleteJSON(ctx, uriLoginSessions, nil)
	c.SetAuthToken("")
	if err != nil && !IsUnauthorized(err) {
		_, err2 := c.DeleteJSON(ctx, uriLoginSessions+"/", nil)
		if err2 == nil || IsUnauthorized(err2) {
			return nil
		}
		return err
	}
	return nil
}

func loginDomainCandidates(d string) []string {
	d = strings.TrimSpace(d)
	if d == "" || strings.EqualFold(d, "local") {
		// Appliance/Postman use "Local"; some docs and env vars use "LOCAL";
		// Global Dashboard swagger uses "local".
		return []string{"Local", "LOCAL", "local", ""}
	}
	return []string{d}
}

func isRetryableLogin(err error) bool {
	var ae *APIError
	if !asAPIError(err, &ae) {
		return false
	}
	switch ae.StatusCode {
	case http.StatusBadRequest, http.StatusNotFound:
		return true
	}
	return false
}
