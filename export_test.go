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

func TestJSONListObjectOrArray(t *testing.T) {
	var a JSONList[PhysicalDrive]
	if err := json.Unmarshal([]byte(`[{"Location":"1:2:3","CapacityMiB":100}]`), &a); err != nil || len(a) != 1 {
		t.Fatalf("array: %v %+v", err, a)
	}
	var b JSONList[PhysicalDrive]
	if err := json.Unmarshal([]byte(`{"Location":"1:2:3","CapacityMiB":100}`), &b); err != nil || len(b) != 1 {
		t.Fatalf("object: %v %+v", err, b)
	}
}

func TestExportServerFromSubresources(t *testing.T) {
	hwJSON := `{
		"name": "Encl1, bay 5",
		"uri": "/rest/server-hardware/abc",
		"uuid": "11111111-1111-1111-1111-111111111111",
		"serialNumber": "CZ123",
		"shortModel": "BL460c Gen10",
		"processorCount": 2,
		"processorCoreCount": 16,
		"processorSpeedMhz": 2600,
		"processorType": "Intel Xeon",
		"memoryMb": 65536,
		"powerState": "On",
		"status": "OK",
		"portMap": {
			"deviceSlots": [{
				"deviceName": "Synergy 3820C",
				"location": "Flb",
				"physicalPorts": [{"portNumber": 1, "type": "Ethernet", "mac": "aa:bb:cc:dd:ee:ff"}]
			}]
		},
		"subResources": {
			"Memory": {
				"collectionState": "Collected",
				"data": [{
					"DeviceLocator": "PROC 1 DIMM 1",
					"CapacityMiB": 32768,
					"MemoryDeviceType": "DDR4",
					"MemoryLocation": {"Socket": 1, "Slot": 1, "Channel": 1, "MemoryController": 1},
					"Oem": {"Hpe": {"DIMMStatus": "GoodInUse"}}
				}]
			},
			"LocalStorage": {
				"collectionState": "Collected",
				"data": [{
					"Name": "Smart Array",
					"Model": "P408i-a",
					"AdapterType": "SmartArray",
					"PhysicalDrives": [{
						"Location": "1:1:1",
						"MediaType": "SSD",
						"InterfaceType": "SATA",
						"CapacityMiB": 228936,
						"Model": "VK000240GWJPH",
						"Status": {"Health": "OK"}
					}],
					"LogicalDrives": [{
						"LogicalDriveName": "RAID1",
						"LogicalDriveNumber": 1,
						"Raid": 1,
						"CapacityMiB": 228936,
						"Status": {"Health": "OK"}
					}]
				}]
			},
			"Devices": {
				"data": [{
					"DeviceType": "LOM/NIC",
					"Name": "Synergy 3820C 10/20/40Gb CNA",
					"Location": "Flb 1",
					"FirmwareVersion": {"Current": {"VersionString": "1.2.3"}}
				}]
			}
		}
	}`
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, `{"minimumVersion":3800,"currentVersion":8800}`)
		case r.URL.Path == "/rest/login-sessions" && r.Method == http.MethodPost:
			io.WriteString(w, `{"sessionID":"t"}`)
		case r.URL.Path == "/rest/server-hardware/abc":
			io.WriteString(w, hwJSON)
		case strings.Contains(r.URL.Path, "/firmware"):
			io.WriteString(w, `{"components":[{"componentName":"System ROM","componentVersion":"I42 01/01/2024","componentLocation":"System Board"}]}`)
		case strings.Contains(r.URL.Path, "/localStorageV2"):
			http.NotFound(w, r)
		case strings.HasSuffix(r.URL.Path, "/bios"), strings.HasSuffix(r.URL.Path, "/environmentalConfiguration"),
			strings.HasSuffix(r.URL.Path, "/memory"), strings.HasSuffix(r.URL.Path, "/memoryList"),
			strings.HasSuffix(r.URL.Path, "/devices"), strings.HasSuffix(r.URL.Path, "/processors"),
			strings.HasSuffix(r.URL.Path, "/localStorage"):
			http.NotFound(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := New(Config{Endpoint: srv.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	exp, err := c.ExportServer(context.Background(), "/rest/server-hardware/abc", ExportOptions{
		SkipProfile: true, SkipHardwareType: true, SkipEnclosure: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if exp.Identity.SerialNumber != "CZ123" {
		t.Fatalf("identity %+v", exp.Identity)
	}
	if exp.Processors.Count != 2 || exp.Processors.TotalCores != 32 {
		t.Fatalf("cpu %+v", exp.Processors)
	}
	if len(exp.Memory.Modules) != 1 || exp.Memory.Modules[0].CapacityMiB != 32768 {
		t.Fatalf("memory %+v", exp.Memory)
	}
	if len(exp.Storage.Drives) != 1 || exp.Storage.Drives[0].MediaType != "SSD" {
		t.Fatalf("drives %+v", exp.Storage)
	}
	if len(exp.Storage.Controllers) != 1 || len(exp.Storage.Controllers[0].Drives) != 1 {
		t.Fatalf("controller disks %+v", exp.Storage.Controllers)
	}
	if len(exp.Storage.Volumes) != 1 || exp.Storage.Volumes[0].RAID != "1" {
		t.Fatalf("volumes %+v", exp.Storage)
	}
	if len(exp.Devices) != 1 || exp.Devices[0].FirmwareVersion != "1.2.3" {
		t.Fatalf("devices %+v", exp.Devices)
	}
	if len(exp.NetworkPorts) != 1 || exp.NetworkPorts[0].MAC != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("nics %+v", exp.NetworkPorts)
	}
	if len(exp.Firmware) != 1 || exp.Firmware[0].Name != "System ROM" {
		t.Fatalf("fw %+v", exp.Firmware)
	}
	sum := exp.Summary()
	if !strings.Contains(sum, "BL460c Gen10") || !strings.Contains(sum, "Disks: 1 physical") {
		t.Fatalf("summary %q", sum)
	}
}

func TestFindServerHardwareByName(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, `{"minimumVersion":3800,"currentVersion":3800}`)
		case r.URL.Path == "/rest/login-sessions":
			io.WriteString(w, `{"sessionID":"t"}`)
		case r.URL.Path == "/rest/server-hardware" && strings.Contains(r.URL.RawQuery, "name"):
			io.WriteString(w, `{"members":[{"name":"bay1","uri":"/rest/server-hardware/abc"}],"count":1,"total":1}`)
		case r.URL.Path == "/rest/server-hardware/abc":
			io.WriteString(w, `{"name":"bay1","uri":"/rest/server-hardware/abc","processorCount":1}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c, _ := New(Config{Endpoint: srv.URL, Username: "a", Password: "b", InsecureTLS: true, HTTPClient: srv.Client()})
	_ = c.Login(context.Background())
	hw, err := c.FindServerHardware(context.Background(), "bay1")
	if err != nil || hw.Name != "bay1" {
		t.Fatalf("%v %+v", err, hw)
	}
}
