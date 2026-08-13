package oneview

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCollectServersDedupAndMerge(t *testing.T) {
	gd := mockAppliance(t, 300, `{
		"count":1,"total":1,"members":[{
			"name":"Encl1, bay 5",
			"uri":"/rest/server-hardware/abc",
			"uuid":"11111111-1111-1111-1111-111111111111",
			"serialNumber":"CZ123",
			"processorCount":2,
			"processorCoreCount":16,
			"processorType":"Intel Xeon",
			"memoryMb":65536,
			"powerState":"On"
		}]
	}`, `{
		"name":"Encl1, bay 5",
		"uri":"/rest/server-hardware/abc",
		"uuid":"11111111-1111-1111-1111-111111111111",
		"serialNumber":"CZ123",
		"shortModel":"BL460c Gen10",
		"processorCount":2,
		"processorCoreCount":16,
		"processorSpeedMhz":2600,
		"processorType":"Intel Xeon",
		"memoryMb":65536,
		"powerState":"On",
		"status":"OK"
	}`)
	defer gd.Close()

	ov := mockAppliance(t, 8800, `{
		"count":1,"total":1,"members":[{
			"name":"Encl1, bay 5",
			"uri":"/rest/server-hardware/abc",
			"uuid":"11111111-1111-1111-1111-111111111111",
			"serialNumber":"CZ123"
		}]
	}`, `{
		"name":"Encl1, bay 5",
		"uri":"/rest/server-hardware/abc",
		"uuid":"11111111-1111-1111-1111-111111111111",
		"serialNumber":"CZ123",
		"shortModel":"BL460c Gen10",
		"processorCount":2,
		"processorCoreCount":16,
		"processorType":"Intel Xeon Gold",
		"memoryMb":65536,
		"subResources":{
			"Processors":{"data":[
				{"Id":"1","Socket":"Cpu1","Model":"Intel Xeon Gold","Manufacturer":"Intel","TotalCores":16,"TotalThreads":32,"MaxSpeedMHz":2600},
				{"Id":"2","Socket":"Cpu2","Model":"Intel Xeon Gold","Manufacturer":"Intel","TotalCores":16,"TotalThreads":32,"MaxSpeedMHz":2600}
			]},
			"Memory":{"data":[
				{"DeviceLocator":"PROC 1 DIMM 1","CapacityMiB":32768,"SerialNumber":"DIMM1","MemoryDeviceType":"DDR4"},
				{"DeviceLocator":"PROC 1 DIMM 2","CapacityMiB":32768,"SerialNumber":"DIMM2","MemoryDeviceType":"DDR4"}
			]},
			"LocalStorage":{"data":[{
				"Name":"Smart Array","Model":"P408i-a","SerialNumber":"CTRL1","AdapterType":"SmartArray",
				"PhysicalDrives":[{"Location":"1:1:1","MediaType":"SSD","CapacityMiB":228936,"Model":"VK000240GWJPH"}],
				"LogicalDrives":[{"LogicalDriveName":"RAID1","Raid":1,"CapacityMiB":228936}]
			}]}
		}
	}`)
	defer ov.Close()

	other := mockAppliance(t, 7400, `{
		"count":1,"total":1,"members":[{
			"name":"rack-01",
			"uri":"/rest/server-hardware/zzz",
			"uuid":"22222222-2222-2222-2222-222222222222",
			"serialNumber":"CZ999"
		}]
	}`, `{
		"name":"rack-01",
		"uri":"/rest/server-hardware/zzz",
		"uuid":"22222222-2222-2222-2222-222222222222",
		"serialNumber":"CZ999",
		"processorCount":null,
		"memoryMb":null,
		"position":null,
		"powerLock":false
	}`)
	defer other.Close()

	servers, err := CollectServers(context.Background(), []Config{
		{Endpoint: gd.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: gd.Client(), APIVersion: 300},
		{Endpoint: ov.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: ov.Client(), APIVersion: 8800},
		{Endpoint: other.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: other.Client(), APIVersion: 7400},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 {
		t.Fatalf("want 2 unique servers, got %d", len(servers))
	}

	var blade, rack *Server
	for i := range servers {
		switch servers[i].Identity.SerialNumber {
		case "CZ123":
			blade = &servers[i]
		case "CZ999":
			rack = &servers[i]
		}
	}
	if blade == nil || rack == nil {
		t.Fatalf("missing server: %+v", servers)
	}
	if blade.Processors == nil || blade.Memory == nil || blade.Controllers == nil {
		t.Fatal("slices must be non-nil")
	}
	if len(blade.Processors) != 2 || blade.Processors[0].Manufacturer != "Intel" {
		t.Fatalf("merged processors %+v", blade.Processors)
	}
	if len(blade.Memory) != 2 || blade.Memory[0].SerialNumber != "DIMM1" {
		t.Fatalf("merged memory %+v", blade.Memory)
	}
	if len(blade.Controllers) != 1 || len(blade.Controllers[0].Drives) != 1 {
		t.Fatalf("controller disks %+v", blade.Controllers)
	}
	if blade.Controllers[0].Drives[0].MediaType != "SSD" {
		t.Fatalf("nested drive %+v", blade.Controllers[0].Drives)
	}
	if len(blade.Controllers[0].Volumes) != 1 {
		t.Fatalf("volumes %+v", blade.Controllers[0].Volumes)
	}
	if len(blade.Sources) != 2 {
		t.Fatalf("sources %+v", blade.Sources)
	}
	if rack.Processors == nil || len(rack.Processors) != 0 {
		t.Fatalf("empty processors want [], got %#v", rack.Processors)
	}
	if rack.Memory == nil || len(rack.Memory) != 0 {
		t.Fatalf("empty memory want [], got %#v", rack.Memory)
	}
	if rack.Controllers == nil || len(rack.Controllers) != 0 {
		t.Fatalf("empty controllers want [], got %#v", rack.Controllers)
	}

	raw, err := json.Marshal(rack)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), `"processors":null`) || strings.Contains(string(raw), `"memory":null`) || strings.Contains(string(raw), `"controllers":null`) {
		t.Fatalf("null arrays in json: %s", raw)
	}
}

func TestCollectServersSkipsFailedConfig(t *testing.T) {
	ok := mockAppliance(t, 3800, `{
		"count":1,"total":1,"members":[{"name":"s1","uri":"/rest/server-hardware/1","uuid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","serialNumber":"OK1"}]
	}`, `{
		"name":"s1","uri":"/rest/server-hardware/1","uuid":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","serialNumber":"OK1"
	}`)
	defer ok.Close()

	servers, err := CollectServers(context.Background(), []Config{
		{Endpoint: "https://127.0.0.1:1", Username: "a", Password: "b", InsecureTLS: true, Timeout: time.Millisecond * 200},
		{Endpoint: ok.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: ok.Client()},
	})
	if err == nil {
		t.Fatal("expected error from unreachable config")
	}
	if len(servers) != 1 || servers[0].Identity.SerialNumber != "OK1" {
		t.Fatalf("partial result %+v", servers)
	}
}

func TestCollectServersEmpty(t *testing.T) {
	servers, err := CollectServers(context.Background(), nil)
	if err != nil || servers == nil || len(servers) != 0 {
		t.Fatalf("%v %#v", err, servers)
	}
}

func mockAppliance(t *testing.T, api int, listJSON, getJSON string) *httptest.Server {
	t.Helper()
	ver := `{"minimumVersion":3800,"currentVersion":8800}`
	if api <= 300 {
		ver = `{"minimumVersion":2,"currentVersion":300}`
	}
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, ver)
		case r.URL.Path == "/rest/login-sessions" && r.Method == http.MethodPost:
			if api <= 300 {
				io.WriteString(w, `{"token":"t"}`)
				return
			}
			io.WriteString(w, `{"sessionID":"t"}`)
		case r.URL.Path == "/rest/login-sessions" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/rest/server-hardware" && r.Method == http.MethodGet:
			io.WriteString(w, listJSON)
		case isServerHardwareResource(r.URL.Path):
			io.WriteString(w, getJSON)
		default:
			http.NotFound(w, r)
		}
	}))
}

func isServerHardwareResource(path string) bool {
	if !strings.HasPrefix(path, "/rest/server-hardware/") {
		return false
	}
	suffixes := []string{
		"/firmware", "/localStorage", "/localStorageV2", "/bios",
		"/memory", "/memoryList", "/devices", "/processors",
		"/environmentalConfiguration",
	}
	for _, s := range suffixes {
		if strings.Contains(path, s) {
			return false
		}
	}
	return true
}
