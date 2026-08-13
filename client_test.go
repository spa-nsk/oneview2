package oneview

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFlexIntAndString(t *testing.T) {
	var n FlexInt
	if err := json.Unmarshal([]byte(`"3800"`), &n); err != nil || n.Int() != 3800 {
		t.Fatalf("string flexint: %v %d", err, n)
	}
	if err := json.Unmarshal([]byte(`8800`), &n); err != nil || n.Int() != 8800 {
		t.Fatalf("number flexint: %v %d", err, n)
	}
	var s FlexString
	if err := json.Unmarshal([]byte(`16`), &s); err != nil || s.String() != "16" {
		t.Fatalf("number flexstring: %v %q", err, s)
	}
	if err := json.Unmarshal([]byte(`"abc"`), &s); err != nil || s.String() != "abc" {
		t.Fatalf("string flexstring: %v %q", err, s)
	}
	var b FlexBool
	if err := json.Unmarshal([]byte(`true`), &b); err != nil || !b.Bool() {
		t.Fatalf("bool flexbool: %v %v", err, b)
	}
	if err := json.Unmarshal([]byte(`"false"`), &b); err != nil || b.Bool() {
		t.Fatalf("string flexbool: %v %v", err, b)
	}
	var hw ServerHardware
	if err := json.Unmarshal([]byte(`{"name":"bay1","powerLock":true,"hostOsType":"12"}`), &hw); err != nil {
		t.Fatal(err)
	}
	if !hw.PowerLock.Bool() || hw.HostOsType.Int() != 12 {
		t.Fatalf("hw %+v", hw)
	}
}

func TestListOptions(t *testing.T) {
	q := ListOptions{
		Start:  10,
		Count:  50,
		Filter: []string{"status EQ 'OK'", "name EQ 's1'"},
		Sort:   "name:asc",
		Query:  "name EQ 'x'",
	}.values()
	if q.Get("start") != "10" || q.Get("count") != "50" {
		t.Fatalf("paging: %v", q)
	}
	if len(q["filter"]) != 2 {
		t.Fatalf("filters: %v", q["filter"])
	}
}

func TestIDFromURI(t *testing.T) {
	if got := IDFromURI("/rest/server-hardware/abc-123/"); got != "abc-123" {
		t.Fatalf("got %q", got)
	}
}

func TestSessionAuthHeader(t *testing.T) {
	ov := Session{SessionID: "sid", Token: "tok"}
	if ov.AuthHeader() != "sid" {
		t.Fatal("appliance sessionID should win")
	}
	gd := Session{Token: "tok"}
	if gd.AuthHeader() != "tok" {
		t.Fatal("GD token")
	}
}

func TestIsSupportedAPI(t *testing.T) {
	if !isSupportedAPI(300) || !isSupportedAPI(3800) || !isSupportedAPI(8800) {
		t.Fatal("expected 300, 3800, 8800 to be supported")
	}
	if isSupportedAPI(1000) || isSupportedAPI(9000) {
		t.Fatal("unexpected supported version")
	}
}

func TestClientLoginSendsAckAndRetriesDomain(t *testing.T) {
	var bodies []string
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, `{"minimumVersion":3800,"currentVersion":3800}`)
		case r.URL.Path == "/rest/login-sessions" && r.Method == http.MethodPost:
			raw, _ := io.ReadAll(r.Body)
			bodies = append(bodies, string(raw))
			if strings.Contains(string(raw), `"loginMsgAck":true`) && strings.Contains(string(raw), `"authLoginDomain":"Local"`) {
				w.Write([]byte(`{"sessionID":"ok"}`))
				return
			}
			w.WriteHeader(400)
			io.WriteString(w, `{"errorCode":"AUTHN_LOGIN_INVALID_MESSAGE_ACK","message":"Login message must be acknowledged.","details":"Set loginMsgAck to true."}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c, _ := New(Config{Endpoint: srv.URL, Username: "admin", Password: "p", Domain: "LOCAL", InsecureTLS: true, HTTPClient: srv.Client()})
	if err := c.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	if c.AuthToken() != "ok" {
		t.Fatalf("token %q bodies=%v", c.AuthToken(), bodies)
	}
}

func TestClientLoginAppliance(t *testing.T) {
	var gotAuth, gotAPI string
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, `{"minimumVersion":3800,"currentVersion":8800}`)
		case r.URL.Path == "/rest/login-sessions" && r.Method == http.MethodPost:
			gotAPI = r.Header.Get("X-API-Version")
			w.Write([]byte(`{"sessionID":"ov-session"}`))
		case r.URL.Path == "/rest/server-hardware":
			gotAuth = r.Header.Get("Auth")
			io.WriteString(w, `{"category":"server-hardware","count":1,"total":1,"members":[{"name":"bay1","uri":"/rest/server-hardware/1","powerState":"On","eTag":12}]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := New(Config{Endpoint: srv.URL, Username: "admin", Password: "p", InsecureTLS: true, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	if c.AuthToken() != "ov-session" {
		t.Fatalf("token %q", c.AuthToken())
	}
	if c.Product() != ProductAppliance {
		t.Fatalf("product %s", c.Product())
	}
	if gotAPI != "8800" {
		t.Fatalf("X-API-Version %q", gotAPI)
	}
	hw, err := c.ListServerHardware(context.Background(), ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "ov-session" || len(hw.Members) != 1 || hw.Members[0].Name != "bay1" {
		t.Fatalf("list: auth=%q hw=%+v", gotAuth, hw)
	}
	if hw.Members[0].ETag.String() != "12" {
		t.Fatalf("etag %q", hw.Members[0].ETag)
	}
}

func TestClientLoginGlobalDashboard(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, `{"minimumVersion":2,"currentVersion":300}`)
		case strings.HasPrefix(r.URL.Path, "/rest/login-sessions"):
			io.WriteString(w, `{"token":"gd-token","user":{"userName":"sarah","domain":"local"}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := New(Config{Endpoint: srv.URL, Username: "sarah", Password: "p", InsecureTLS: true, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	if c.Product() != ProductGlobalDashboard {
		t.Fatalf("product %s", c.Product())
	}
	if c.APIVersion() != 300 {
		t.Fatalf("api %d", c.APIVersion())
	}
	if c.AuthToken() != "gd-token" {
		t.Fatalf("token %q", c.AuthToken())
	}
}

func TestAPIError(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, `{"errorCode":"INVALID_RESOURCE_URI","message":"Resource not found."}`)
	}))
	defer srv.Close()
	c, _ := New(Config{Endpoint: srv.URL, InsecureTLS: true, HTTPClient: srv.Client(), APIVersion: 3800})
	c.setVersion(3800, 8800, 3800)
	_, err := c.GetServerHardware(context.Background(), "missing")
	if !IsNotFound(err) {
		t.Fatalf("expected 404, got %v", err)
	}
}

func TestWaitTask(t *testing.T) {
	n := 0
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n++
		if n == 1 {
			io.WriteString(w, `{"uri":"/rest/tasks/1","taskState":"Running","percentComplete":10}`)
			return
		}
		io.WriteString(w, `{"uri":"/rest/tasks/1","taskState":"Completed","percentComplete":100}`)
	}))
	defer srv.Close()
	c, _ := New(Config{Endpoint: srv.URL, InsecureTLS: true, HTTPClient: srv.Client(), APIVersion: 3800})
	c.setVersion(3800, 8800, 3800)
	task, err := c.WaitTaskInterval(context.Background(), "/rest/tasks/1", 1)
	if err != nil {
		t.Fatal(err)
	}
	if task.TaskState != "Completed" {
		t.Fatalf("%+v", task)
	}
}
