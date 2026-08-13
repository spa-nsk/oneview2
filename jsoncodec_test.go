package oneview

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// GD swagger 300 example shape: eTag as timestamp-number string, hostOsType null,
// powerLock bool (schema says string), description null, extra GD fields.
const swagger300HardwareJSON = `{
  "applianceLocation": "172.1.1.11",
  "applianceName": "ci-005056b3524e",
  "appluri": "/rest/appliances/7df6c7fb-8d5d-4fef-94ea-cc29799c44cc",
  "assetTag": "[Unknown]",
  "category": "server-hardware",
  "created": "2017-07-16T22:08:19.816Z",
  "description": null,
  "eTag": "1505729209522",
  "formFactor": "FullHeight",
  "generation": "Gen10",
  "groups": [{"name": "group-lab", "uri": "/rest/groups/96c3af19-b928-4719-96dc-4e51e3dc51fb"}],
  "hostOsType": null,
  "id": "660b26ce-8201-4484-9fcb-2fd16c977e40",
  "licensingIntent": "OneView",
  "locationUri": "/rest/enclosures/09SGH100X6J1",
  "memoryMb": 16384,
  "model": "ProLiant BL660c Gen10",
  "modified": "2017-09-18T10:06:49.522Z",
  "mpDnsName": "172.18.6.1",
  "mpFirmwareVersion": "2.20 Nov 01 2014",
  "mpHostInfo": {"mpHostName": "172.18.6.1", "mpIpAddresses": [{"address": "172.18.6.1", "type": "Undefined"}]},
  "mpIpAddress": "172.18.6.1",
  "mpModel": "iLO4",
  "mpState": "OK",
  "name": "Encl1, bay 1",
  "originalUri": "/rest/server-hardware/31393736-3831-4753-4831-30305837524E",
  "partNumber": "679118-B21",
  "physicalServerHardwareUri": null,
  "platform": "BladeServer",
  "position": 1,
  "powerLock": false,
  "powerState": "Off",
  "processorCoreCount": 8,
  "processorCount": 4,
  "processorSpeedMhz": 2000,
  "processorType": "Quad-Core Intel Xeon @ 2.0GHz",
  "refreshState": "NotRefreshing",
  "serialNumber": "SGH100X7RN",
  "serverGroupUri": "/rest/enclosure-groups/b43512bd-4ba9-4438-bdf5-812ad394bbb9",
  "serverHardwareTypeUri": "/rest/server-hardware-types/1987B1BF-245A-434A-9AAC-C4B76931BE9D",
  "serverName": "",
  "serverProfileUri": "/rest/server-profiles/a00c15ea-f98d-4e8e-b253-6fd4288ab4c8",
  "shortModel": "BL660c Gen9",
  "state": "ProfileApplied",
  "stateReason": "NotApplicable",
  "status": "Warning",
  "type": "server-hardware",
  "uidState": "Off",
  "uri": "/rest/server-hardware/660b26ce-8201-4484-9fcb-2fd16c977e40",
  "uuid": "660b26ce-8201-4484-9fcb-2fd16c977e40",
  "virtualSerialNumber": null,
  "virtualUuid": null
}`

// Appliance 3800–8800 shape: numeric eTag, string powerLock (GD schema),
// hostOsType as OS name, null position on rack servers, extra scope fields.
const appliance8800HardwareJSON = `{
  "type": "server-hardware-14",
  "uri": "/rest/server-hardware/30373237-3132-4D32-3235-303930524D57",
  "category": "server-hardware",
  "eTag": 1441147370086,
  "name": "rack-01",
  "description": null,
  "status": "OK",
  "state": "NoProfileApplied",
  "configurationState": "Managed",
  "generation": "Gen11",
  "hostname": "ilo.example.com",
  "hostOsType": "Windows",
  "intelligentProvisioningVersion": "4.12.0.0",
  "licensingIntent": "OneView",
  "locationUri": null,
  "maintenanceMode": "false",
  "memoryMb": "262144",
  "model": "ProLiant DL380 Gen11",
  "mpHostsAndRanges": ["172.1.1.1-172.1.1.10"],
  "oneTimeBoot": "USB",
  "position": null,
  "powerLock": "true",
  "powerState": "On",
  "processorCoreCount": null,
  "processorCount": 2,
  "processorSpeedMhz": 2300,
  "scopesUri": "/rest/scopes/resources/rest/server-hardware/1",
  "scopeUris": ["/rest/scopes/Default"],
  "initialScopeUris": [],
  "uidState": "Off",
  "uuid": "30373237-3132-4D32-3235-303930524D57"
}`

func TestDecodeJSONSwagger300Hardware(t *testing.T) {
	var hw ServerHardware
	if err := DecodeJSON([]byte(swagger300HardwareJSON), &hw); err != nil {
		t.Fatal(err)
	}
	if hw.Name != "Encl1, bay 1" {
		t.Fatalf("name %q", hw.Name)
	}
	if hw.ETag.String() != "1505729209522" {
		t.Fatalf("etag %q", hw.ETag)
	}
	if hw.HostOsType.String() != "" {
		t.Fatalf("null hostOsType %q", hw.HostOsType)
	}
	if hw.PowerLock.Bool() {
		t.Fatal("powerLock false")
	}
	if hw.MemoryMb.Int() != 16384 || hw.Position.Int() != 1 || hw.ProcessorCount.Int() != 4 {
		t.Fatalf("ints %+v", hw)
	}
	if hw.ApplianceLocation != "172.1.1.11" || hw.Generation != "Gen10" {
		t.Fatalf("gd fields %+v", hw)
	}
	if hw.Description != "" {
		t.Fatalf("null description %q", hw.Description)
	}
}

func TestDecodeJSONAppliance8800Hardware(t *testing.T) {
	var hw ServerHardware
	if err := DecodeJSON([]byte(appliance8800HardwareJSON), &hw); err != nil {
		t.Fatal(err)
	}
	if hw.ETag.String() != "1441147370086" {
		t.Fatalf("numeric etag %q", hw.ETag)
	}
	if hw.HostOsType.String() != "Windows" {
		t.Fatalf("hostOsType %q", hw.HostOsType)
	}
	if !hw.PowerLock.Bool() {
		t.Fatal("string powerLock true")
	}
	if hw.Position.Int() != 0 {
		t.Fatalf("null position %d", hw.Position)
	}
	if hw.MemoryMb.Int() != 262144 {
		t.Fatalf("string memoryMb %d", hw.MemoryMb)
	}
	if hw.MaintenanceMode.Bool() {
		t.Fatal("string maintenanceMode false")
	}
	if hw.ConfigurationState != "Managed" || hw.Hostname != "ilo.example.com" {
		t.Fatalf("3800 fields %+v", hw)
	}
	if hw.ScopesURI == "" || hw.OneTimeBoot != "USB" || hw.Generation != "Gen11" {
		t.Fatalf("extra 3800 fields %+v", hw)
	}
	if hw.ProcessorCoreCount.Int() != 0 || hw.ProcessorCount.Int() != 2 {
		t.Fatalf("cpu %+v", hw)
	}
}

func TestDecodeJSONCoercesNestedAndCollections(t *testing.T) {
	raw := `{
		"type": "server-hardware-list-14",
		"start": "0",
		"count": 1,
		"total": "1",
		"members": {
			"name": "solo",
			"eTag": 12,
			"powerLock": "false",
			"hideUnusedFlexNics": "true",
			"enclosureBay": null,
			"memoryMb": 1024
		}
	}`
	var col Collection[ServerHardware]
	if err := DecodeJSON([]byte(raw), &col); err != nil {
		t.Fatal(err)
	}
	if col.Start != 0 || col.Total != 1 || len(col.Members) != 1 {
		t.Fatalf("collection %+v", col)
	}
	if col.Members[0].Name != "solo" || col.Members[0].ETag.String() != "12" {
		t.Fatalf("member %+v", col.Members[0])
	}
	if col.Members[0].PowerLock.Bool() {
		t.Fatal("powerLock")
	}

	var p ServerProfile
	if err := DecodeJSON([]byte(`{"enclosureBay":null,"hideUnusedFlexNics":"true","inProgress":0}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.EnclosureBay.Int() != 0 || !p.HideUnusedFlexNics.Bool() || p.InProgress.Bool() {
		t.Fatalf("profile %+v", p)
	}

	var task Task
	if err := DecodeJSON([]byte(`{"percentComplete":100.0,"hidden":"false","userInitiated":1}`), &task); err != nil {
		t.Fatal(err)
	}
	if task.PercentComplete.Int() != 100 || task.Hidden.Bool() || !task.UserInitiated.Bool() {
		t.Fatalf("task %+v", task)
	}
}

func TestGetServerHardwareMixedVersions(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/version":
			io.WriteString(w, `{"minimumVersion":3800,"currentVersion":8800}`)
		case r.URL.Path == "/rest/login-sessions":
			w.Write([]byte(`{"sessionID":"s"}`))
		case r.URL.Path == "/rest/server-hardware/gd":
			io.WriteString(w, swagger300HardwareJSON)
		case r.URL.Path == "/rest/server-hardware/ov":
			io.WriteString(w, appliance8800HardwareJSON)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := New(Config{Endpoint: srv.URL, Username: "a", Password: "p", InsecureTLS: true, HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	gd, err := c.GetServerHardware(context.Background(), "gd")
	if err != nil {
		t.Fatal(err)
	}
	if gd.Name != "Encl1, bay 1" || gd.PowerLock.Bool() || gd.MemoryMb.Int() != 16384 {
		t.Fatalf("gd %+v", gd)
	}
	ov, err := c.GetServerHardware(context.Background(), "ov")
	if err != nil {
		t.Fatal(err)
	}
	if ov.HostOsType.String() != "Windows" || !ov.PowerLock.Bool() || ov.Position.Int() != 0 {
		t.Fatalf("ov %+v", ov)
	}
}

func TestDecodeJSONRejectsInvalid(t *testing.T) {
	var hw ServerHardware
	if err := DecodeJSON([]byte(`{`), &hw); err == nil {
		t.Fatal("expected error")
	}
}

func TestStrictUnmarshalStillWorksForAlignedJSON(t *testing.T) {
	var hw ServerHardware
	if err := json.Unmarshal([]byte(`{"name":"x","powerLock":true,"memoryMb":8}`), &hw); err != nil {
		t.Fatal(err)
	}
	if hw.Name != "x" || !hw.PowerLock.Bool() || hw.MemoryMb.Int() != 8 {
		t.Fatalf("%+v", hw)
	}
}
