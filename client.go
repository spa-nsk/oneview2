package oneview

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultUserAgent = "oneview-go/1.0"

// Config configures a Client.
type Config struct {
	// Endpoint is the appliance base URL, e.g. https://oneview.example.com
	Endpoint string
	Username string
	Password string
	// Domain is the login domain (LOCAL for local users, or an AD/LDAP directory name).
	Domain string
	// APIVersion is X-API-Version. Zero means use the appliance currentVersion.
	APIVersion int
	// InsecureTLS skips TLS certificate verification (lab use only).
	InsecureTLS bool
	Timeout     time.Duration
	IfMatch     string
	UserAgent   string
	HTTPClient  *http.Client
}

// ConfigFromEnv reads ONEVIEW_OV_ENDPOINT, ONEVIEW_OV_USER, ONEVIEW_OV_PASSWORD,
// ONEVIEW_OV_DOMAIN, ONEVIEW_APIVERSION, ONEVIEW_SSLVERIFY — the same names as
// the official HPE oneview-golang SDK.
func ConfigFromEnv() Config {
	insecure := true
	if v := strings.ToLower(os.Getenv("ONEVIEW_SSLVERIFY")); v == "true" || v == "1" {
		insecure = false
	}
	api, _ := strconv.Atoi(os.Getenv("ONEVIEW_APIVERSION"))
	return Config{
		Endpoint:    os.Getenv("ONEVIEW_OV_ENDPOINT"),
		Username:    os.Getenv("ONEVIEW_OV_USER"),
		Password:    os.Getenv("ONEVIEW_OV_PASSWORD"),
		Domain:      os.Getenv("ONEVIEW_OV_DOMAIN"),
		APIVersion:  api,
		InsecureTLS: insecure,
	}
}

// Client is an HPE OneView / Global Dashboard REST client.
type Client struct {
	base       *url.URL
	http       *http.Client
	apiVersion int
	username   string
	password   string
	domain     string
	ifMatch    string
	userAgent  string

	mu             sync.RWMutex
	authToken      string
	minVersion     int
	currentVersion int
	product        Product
}

// New builds a client. It does not contact the appliance until Login or GetVersion.
func New(cfg Config) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("oneview: Endpoint is required")
	}
	raw := strings.TrimRight(cfg.Endpoint, "/")
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("oneview: invalid Endpoint: %w", err)
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		if cfg.InsecureTLS {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
		}
		httpClient = &http.Client{Timeout: timeout, Transport: tr}
	}
	ua := cfg.UserAgent
	if ua == "" {
		ua = defaultUserAgent
	}
	domain := cfg.Domain
	if domain == "" {
		domain = "LOCAL"
	}
	ifMatch := cfg.IfMatch
	if ifMatch == "" {
		ifMatch = "*"
	}
	return &Client{
		base:       u,
		http:       httpClient,
		apiVersion: cfg.APIVersion,
		username:   cfg.Username,
		password:   cfg.Password,
		domain:     domain,
		ifMatch:    ifMatch,
		userAgent:  ua,
	}, nil
}

// APIVersion returns the X-API-Version currently sent on requests.
func (c *Client) APIVersion() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiVersion
}

// Product returns the detected appliance family after GetVersion or Login.
func (c *Client) Product() Product {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.product
}

// IsGlobalDashboard reports whether the target is OneView Global Dashboard.
func (c *Client) IsGlobalDashboard() bool {
	return c.Product() == ProductGlobalDashboard
}

// ApplianceVersion returns (minimum, current) after GetVersion or Login.
func (c *Client) ApplianceVersion() (min, current int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.minVersion, c.currentVersion
}

// AuthToken returns the current session token, if any.
func (c *Client) AuthToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authToken
}

// SetAuthToken sets a previously obtained session token (skips Login).
func (c *Client) SetAuthToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.authToken = token
}

// SetAPIVersion overrides X-API-Version for subsequent requests.
func (c *Client) SetAPIVersion(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiVersion = v
}

// BaseURL returns the appliance origin.
func (c *Client) BaseURL() string {
	return strings.TrimRight(c.base.String(), "/")
}

func (c *Client) setVersion(min, current, requested int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.minVersion = min
	c.currentVersion = current
	c.product = detectProduct(current)
	if requested > 0 {
		c.apiVersion = requested
	} else if c.apiVersion == 0 {
		c.apiVersion = current
	}
}
