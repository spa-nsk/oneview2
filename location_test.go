package oneview

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectServersLocation(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest/version":
			io.WriteString(w, `{"minimumVersion":3800,"currentVersion":8800}`)
		case "/rest/login-sessions":
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			io.WriteString(w, `{"sessionID":"t"}`)
		case "/rest/server-hardware":
			io.WriteString(w, `{
				"count":2,"total":2,"members":[
					{"name":"Encl1, bay 5","uri":"/rest/server-hardware/abc","uuid":"11111111-1111-1111-1111-111111111111","serialNumber":"CZ123"},
					{"name":"rack-01","uri":"/rest/server-hardware/dl1","uuid":"22222222-2222-2222-2222-222222222222","serialNumber":"DL001"}
				]
			}`)
		case "/rest/server-hardware/abc":
			io.WriteString(w, `{
				"name":"Encl1, bay 5",
				"uri":"/rest/server-hardware/abc",
				"uuid":"11111111-1111-1111-1111-111111111111",
				"serialNumber":"CZ123",
				"locationUri":"/rest/enclosures/enc1",
				"position":5,
				"serverProfileUri":"/rest/server-profiles/p1",
				"powerState":"On"
			}`)
		case "/rest/server-hardware/dl1":
			io.WriteString(w, `{
				"name":"rack-01",
				"uri":"/rest/server-hardware/dl1",
				"uuid":"22222222-2222-2222-2222-222222222222",
				"serialNumber":"DL001",
				"locationUri":null,
				"position":null,
				"powerState":"Off"
			}`)
		case "/rest/enclosures/enc1":
			io.WriteString(w, `{
				"name":"Encl1",
				"uri":"/rest/enclosures/enc1",
				"uuid":"eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee",
				"serialNumber":"SGH100X6J1",
				"enclosureModel":"BladeSystem c7000 Enclosure G3",
				"enclosureType":"C7000",
				"rackName":"Rack-221",
				"deviceBayCount":16,
				"status":"OK",
				"state":"Configured",
				"activeOaPreferredIP":"172.18.1.11",
				"deviceBays":[
					{"type":"DeviceBayV200","bayNumber":5,"devicePresence":"Present","deviceUri":"/rest/server-hardware/abc","deviceBayType":"C7000DeviceBay","bayPowerState":"On"}
				]
			}`)
		case "/rest/racks":
			io.WriteString(w, `{
				"count":1,"total":1,"members":[{
					"name":"Rack-221",
					"uri":"/rest/racks/rk1",
					"uuid":"rrrrrrrr-rrrr-rrrr-rrrr-rrrrrrrrrrrr",
					"model":"10642 G2",
					"uHeight":42,
					"serialNumber":"RACKSN",
					"rackMounts":[
						{"mountUri":"/rest/enclosures/enc1","topUSlot":20,"uHeight":10,"location":"CenterFront"},
						{"mountUri":"/rest/server-hardware/dl1","topUSlot":4,"uHeight":2,"location":"CenterFront"}
					]
				}]
			}`)
		case "/rest/datacenters":
			io.WriteString(w, `{
				"count":1,"total":1,"members":[{
					"name":"Datacenter 1",
					"uri":"/rest/datacenters/dc1",
					"uuid":"dddddddd-dddd-dddd-dddd-dddddddddddd",
					"status":"OK",
					"coolingCapacity":100,
					"width":10,
					"height":20,
					"contents":[{"resourceUri":"/rest/racks/rk1","x":1,"y":2}]
				}]
			}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	servers, err := CollectServers(context.Background(), []Config{
		{Endpoint: srv.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: srv.Client()},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 {
		t.Fatalf("got %d servers", len(servers))
	}

	var blade, rack *Server
	for i := range servers {
		switch servers[i].Identity.SerialNumber {
		case "CZ123":
			blade = &servers[i]
		case "DL001":
			rack = &servers[i]
		}
	}
	if blade == nil || rack == nil {
		t.Fatalf("missing: %+v", servers)
	}

	if blade.EnclosureInfo.Name != "Encl1" || blade.EnclosureInfo.EnclosureType != "C7000" {
		t.Fatalf("enclosure %+v", blade.EnclosureInfo)
	}
	if blade.BayInfo.BayNumber != 5 || blade.BayInfo.DevicePresence != "Present" || blade.BayInfo.DeviceBayType != "C7000DeviceBay" {
		t.Fatalf("bay %+v", blade.BayInfo)
	}
	if blade.Rack.Name != "Rack-221" || blade.Rack.TopUSlot != 20 || blade.Rack.UHeight != 42 {
		t.Fatalf("blade rack %+v", blade.Rack)
	}
	if blade.Datacenter.Name != "Datacenter 1" || blade.Datacenter.URI != "/rest/datacenters/dc1" {
		t.Fatalf("blade dc %+v", blade.Datacenter)
	}

	if rack.EnclosureInfo.URI != "" || rack.BayInfo.BayNumber != 0 {
		t.Fatalf("rack server should have empty enclosure/bay: enc=%+v bay=%+v", rack.EnclosureInfo, rack.BayInfo)
	}
	if rack.Rack.Name != "Rack-221" || rack.Rack.TopUSlot != 4 || rack.Rack.MountUHeight != 2 {
		t.Fatalf("rack-mount placement %+v", rack.Rack)
	}
	if rack.Datacenter.Name != "Datacenter 1" {
		t.Fatalf("rack dc %+v", rack.Datacenter)
	}
}

func TestCollectServersLocationFromDashboardInventory(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest/version":
			io.WriteString(w, `{"minimumVersion":2,"currentVersion":300}`)
		case "/rest/login-sessions":
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			io.WriteString(w, `{"token":"t"}`)
		case "/rest/server-hardware":
			io.WriteString(w, `{"count":1,"total":1,"members":[{"name":"Encl1, bay 1","uri":"/rest/server-hardware/abc","uuid":"11111111-1111-1111-1111-111111111111","serialNumber":"CZ123"}]}`)
		case "/rest/server-hardware/abc":
			io.WriteString(w, `{
				"name":"Encl1, bay 1",
				"uri":"/rest/server-hardware/abc",
				"originalUri":"/rest/server-hardware/abc",
				"uuid":"11111111-1111-1111-1111-111111111111",
				"serialNumber":"CZ123",
				"locationUri":"/rest/enclosures/enc1",
				"position":1
			}`)
		case "/rest/enclosures/enc1":
			io.WriteString(w, `{"name":"Encl1","uri":"/rest/enclosures/enc1","rackName":"Rack-241","deviceBayCount":16}`)
		case "/rest/datacenters":
			io.WriteString(w, `{
				"count":1,"total":1,"members":[{
					"name":"Datacenter 1",
					"uri":"/rest/datacenters/dc1",
					"rackInventory":[{"name":"Rack-241","originalUri":"/rest/racks/Rack-241"}]
				}]
			}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	servers, err := CollectServers(context.Background(), []Config{
		{Endpoint: srv.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: srv.Client(), APIVersion: 300},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 1 {
		t.Fatalf("got %d", len(servers))
	}
	s := servers[0]
	if s.EnclosureInfo.Name != "Encl1" || s.BayInfo.BayNumber != 1 {
		t.Fatalf("enc/bay %+v %+v", s.EnclosureInfo, s.BayInfo)
	}
	if s.Rack.Name != "Rack-241" {
		t.Fatalf("rack name from enclosure: %+v", s.Rack)
	}
	if s.Datacenter.Name != "Datacenter 1" {
		t.Fatalf("dc from rackInventory: %+v", s.Datacenter)
	}
}
