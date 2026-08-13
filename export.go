package oneview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FindServerHardware locates a server by URI, UUID, name, serverName or serialNumber.
func (c *Client) FindServerHardware(ctx context.Context, ident string) (*ServerHardware, error) {
	ident = strings.TrimSpace(ident)
	if ident == "" {
		return nil, fmt.Errorf("oneview: empty server identifier")
	}
	if strings.HasPrefix(ident, "/rest/server-hardware/") || looksLikeUUID(ident) {
		hw, err := c.GetServerHardware(ctx, ident)
		if err == nil {
			return hw, nil
		}
		if !IsNotFound(err) {
			return nil, err
		}
	}
	filters := []string{
		fmt.Sprintf("name='%s'", ident),
		fmt.Sprintf("serverName='%s'", ident),
		fmt.Sprintf("serialNumber='%s'", ident),
		fmt.Sprintf("uuid='%s'", ident),
	}
	for _, f := range filters {
		col, err := c.ListServerHardware(ctx, ListOptions{Count: -1, Filter: []string{f}})
		if err != nil {
			continue
		}
		if len(col.Members) == 1 {
			return c.GetServerHardware(ctx, col.Members[0].URI)
		}
		if len(col.Members) > 1 {
			return nil, fmt.Errorf("oneview: identifier %q matched %d servers", ident, len(col.Members))
		}
	}
	return nil, fmt.Errorf("oneview: server %q not found", ident)
}

// ExportServer dumps the detailed configuration of one server (CPU, RAM, disks,
// PCI devices, NICs, firmware, BIOS, profile). ident is URI, UUID, name or serial.
func (c *Client) ExportServer(ctx context.Context, ident string, opts ExportOptions) (*ServerExport, error) {
	hw, err := c.FindServerHardware(ctx, ident)
	if err != nil {
		return nil, err
	}
	return c.exportHardware(ctx, hw, opts)
}

// ExportServers dumps every server matching listOpts (use Count:-1 for all).
func (c *Client) ExportServers(ctx context.Context, listOpts ListOptions, opts ExportOptions) ([]*ServerExport, error) {
	col, err := c.ListServerHardware(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	out := make([]*ServerExport, 0, len(col.Members))
	for i := range col.Members {
		exp, err := c.exportHardware(ctx, &col.Members[i], opts)
		if err != nil {
			return out, fmt.Errorf("oneview: export %s: %w", col.Members[i].Name, err)
		}
		out = append(out, exp)
	}
	return out, nil
}

// WriteServerExportJSON writes one dump as indented JSON.
func WriteServerExportJSON(w io.Writer, exp *ServerExport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(exp)
}

// SaveServerExportJSON writes one dump to path.
func SaveServerExportJSON(path string, exp *ServerExport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && !os.IsExist(err) {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteServerExportJSON(f, exp)
}

func (c *Client) exportHardware(ctx context.Context, listed *ServerHardware, opts ExportOptions) (*ServerExport, error) {
	hw := listed
	if listed.URI != "" {
		full, err := c.GetServerHardware(ctx, listed.URI)
		if err == nil {
			hw = full
		}
	}
	exp := &ServerExport{
		CollectedAt: time.Now().UTC(),
		APIVersion:  c.APIVersion(),
		Product:     c.Product().String(),
		Hardware:    hw,
		Identity:    identityFrom(hw),
	}
	exp.Processors = processorsFrom(hw, nil)
	exp.Memory.TotalMiB = hw.MemoryMb
	exp.NetworkPorts = networkPortsFrom(hw.PortMap)

	if !opts.SkipSubresources {
		c.fillSubresources(ctx, hw, exp)
	} else {
		parseEmbeddedSubresources(hw, exp)
	}

	if !opts.SkipFirmware {
		c.fillFirmware(ctx, hw, exp)
	}
	if !opts.SkipLocalStorage {
		c.fillLocalStorageV2(ctx, hw, exp)
	}
	if !opts.SkipBIOS {
		c.fillOptionalJSON(ctx, joinPath(hw.URI, "bios"), &exp.BIOS, exp, "bios")
	}
	if !opts.SkipEnvironment {
		c.fillOptionalJSON(ctx, joinPath(hw.URI, "environmentalConfiguration"), &exp.Environment, exp, "environmentalConfiguration")
	}
	if !opts.SkipProfile && hw.ServerProfileURI != "" {
		p, err := c.GetServerProfile(ctx, hw.ServerProfileURI)
		if err != nil {
			exp.warn("profile: %v", err)
		} else {
			exp.Profile = p
		}
	}
	if !opts.SkipHardwareType && hw.ServerHardwareTypeURI != "" {
		t, err := c.GetServerHardwareType(ctx, hw.ServerHardwareTypeURI)
		if err != nil {
			exp.warn("hardware type: %v", err)
		} else {
			exp.HardwareType = t
		}
	}
	if !opts.SkipEnclosure && hw.LocationURI != "" {
		enc, err := c.GetEnclosure(ctx, hw.LocationURI)
		if err != nil {
			exp.warn("enclosure: %v", err)
		} else {
			exp.Enclosure = enc
		}
	}
	return exp, nil
}

func (c *Client) fillSubresources(ctx context.Context, hw *ServerHardware, exp *ServerExport) {
	parseEmbeddedSubresources(hw, exp)
	id := IDFromURI(hw.URI)
	fetch := []struct {
		name string
		path string
		fn   func([]byte)
	}{
		{"memory", joinPath(URIServerHardware, id, "memory"), func(b []byte) {
			exp.Memory.Modules = mergeMemory(exp.Memory.Modules, decodeSubData[rawMemory](b, exp, "memory"))
		}},
		{"memoryList", joinPath(URIServerHardware, id, "memoryList"), func(b []byte) {
			exp.Memory.Boards = mergeBoards(exp.Memory.Boards, decodeSubData[MemoryBoard](b, exp, "memoryList"))
		}},
		{"devices", joinPath(URIServerHardware, id, "devices"), func(b []byte) {
			exp.Devices = mergeDevices(exp.Devices, decodeSubData[rawDevice](b, exp, "devices"))
		}},
		{"processors", joinPath(URIServerHardware, id, "processors"), func(b []byte) {
			socks := decodeSubData[rawProcessor](b, exp, "processors")
			if len(socks) > 0 {
				exp.Processors = processorsFrom(hw, socks)
			}
		}},
		{"localStorage", joinPath(URIServerHardware, id, "localStorage"), func(b []byte) {
			if len(exp.Storage.Controllers) > 0 {
				return
			}
			applyLocalStorage(exp, decodeSubData[rawController](b, exp, "localStorage"), "localStorage")
		}},
	}
	for _, f := range fetch {
		var raw json.RawMessage
		if err := c.GetJSON(ctx, f.path, &raw); err != nil {
			if !skipOptional(err) {
				exp.warn("%s: %v", f.name, err)
			}
			continue
		}
		f.fn(raw)
	}
}

func (c *Client) fillFirmware(ctx context.Context, hw *ServerHardware, exp *ServerExport) {
	path := joinPath(hw.URI, "firmware")
	var body json.RawMessage
	if err := c.GetJSON(ctx, path, &body); err != nil {
		if c.IsGlobalDashboard() {
			col, err2 := c.ListServerFirmware(ctx, ListOptions{
				Count:  -1,
				Filter: []string{fmt.Sprintf("serverHardwareUri='%s'", hw.OriginalURI)},
			})
			if err2 != nil {
				col, err2 = c.ListServerFirmware(ctx, ListOptions{
					Count: -1,
					Query: fmt.Sprintf("serverHardwareUri EQ '%s'", hw.URI),
				})
			}
			if err2 != nil {
				exp.warn("firmware: %v", err)
				return
			}
			for _, row := range col.Members {
				exp.Firmware = append(exp.Firmware, FirmwareComponent{
					Name:     row.Name,
					Version:  row.Version,
					Location: row.Location,
					State:    row.State,
					Status:   row.Status,
				})
			}
			return
		}
		if !skipOptional(err) {
			exp.warn("firmware: %v", err)
		}
		return
	}
	exp.Firmware = parseFirmwareBody(body)
}

func (c *Client) fillLocalStorageV2(ctx context.Context, hw *ServerHardware, exp *ServerExport) {
	var raw json.RawMessage
	if err := c.GetJSON(ctx, joinPath(hw.URI, "localStorageV2"), &raw); err != nil {
		if !skipOptional(err) {
			exp.warn("localStorageV2: %v", err)
		}
		return
	}
	ctrls := decodeSubData[rawController](raw, exp, "localStorageV2")
	if len(ctrls) == 0 {
		var wrap struct {
			Controllers JSONList[rawController] `json:"controllers"`
			Members     JSONList[rawController] `json:"members"`
		}
		if json.Unmarshal(raw, &wrap) == nil {
			ctrls = wrap.Controllers
			if len(ctrls) == 0 {
				ctrls = wrap.Members
			}
		}
	}
	if len(ctrls) > 0 {
		applyLocalStorage(exp, ctrls, "localStorageV2")
	}
}

func (c *Client) fillOptionalJSON(ctx context.Context, path string, dest *json.RawMessage, exp *ServerExport, name string) {
	var raw json.RawMessage
	if err := c.GetJSON(ctx, path, &raw); err != nil {
		if !skipOptional(err) {
			exp.warn("%s: %v", name, err)
		}
		return
	}
	*dest = raw
}

func skipOptional(err error) bool {
	var ae *APIError
	if !asAPIError(err, &ae) {
		return false
	}
	switch ae.StatusCode {
	case 400, 404, 405, 412:
		return true
	}
	return false
}

func looksLikeUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}

func identityFrom(hw *ServerHardware) ServerIdentity {
	host := hw.MpHostInfo.MpHostName
	return ServerIdentity{
		Name:                   hw.Name,
		ServerName:             hw.ServerName,
		URI:                    hw.URI,
		UUID:                   first(hw.UUID, hw.ID),
		SerialNumber:           hw.SerialNumber,
		Model:                  hw.Model,
		ShortModel:             hw.ShortModel,
		PartNumber:             hw.PartNumber,
		AssetTag:               hw.AssetTag,
		FormFactor:             hw.FormFactor,
		Platform:               hw.Platform,
		PowerState:             hw.PowerState,
		Status:                 hw.Status,
		State:                  hw.State,
		OperatingSystem:        hw.OperatingSystem,
		RomVersion:             hw.RomVersion,
		MpModel:                hw.MpModel,
		MpFirmwareVersion:      hw.MpFirmwareVersion,
		MpIPAddress:            hw.MpIpAddress,
		MpHostName:             host,
		Position:               hw.Position,
		LocationURI:            hw.LocationURI,
		ServerProfileURI:       hw.ServerProfileURI,
		ServerHardwareTypeURI:  hw.ServerHardwareTypeURI,
		ServerHardwareTypeName: hw.ServerHardwareTypeName,
		ApplianceLocation:      hw.ApplianceLocation,
		ApplianceName:          hw.ApplianceName,
	}
}

func first(vals ...string) string { return nonEmpty(vals...) }

func processorsFrom(hw *ServerHardware, socks []rawProcessor) ProcessorInventory {
	inv := ProcessorInventory{
		Count:       hw.ProcessorCount,
		CoresPerCPU: hw.ProcessorCoreCount,
		SpeedMHz:    hw.ProcessorSpeedMhz,
		Model:       hw.ProcessorType,
	}
	if inv.Count > 0 && inv.CoresPerCPU > 0 {
		inv.TotalCores = inv.Count * inv.CoresPerCPU
	}
	for _, s := range socks {
		p := ProcessorInfo{
			ID:                s.ID,
			Socket:            s.Socket,
			Model:             s.Model,
			Manufacturer:      s.Manufacturer,
			TotalCores:        s.TotalCores,
			TotalThreads:      s.TotalThreads,
			MaxSpeedMHz:       s.MaxSpeedMHz,
			OperatingSpeedMHz: s.OperatingSpeedMHz,
			Health:            s.Status.Health,
			State:             s.Status.State,
		}
		if p.Model == "" {
			p.Model = hw.ProcessorType
		}
		inv.Sockets = append(inv.Sockets, p)
	}
	if len(inv.Sockets) > 0 {
		inv.Count = len(inv.Sockets)
		var cores int
		for _, s := range inv.Sockets {
			cores += s.TotalCores
			if inv.Model == "" {
				inv.Model = s.Model
			}
		}
		if cores > 0 {
			inv.TotalCores = cores
		}
	}
	return inv
}

func parseEmbeddedSubresources(hw *ServerHardware, exp *ServerExport) {
	if hw == nil || hw.SubResources == nil {
		return
	}
	if sr, ok := hw.SubResources["Memory"]; ok {
		exp.Memory.Modules = mergeMemory(exp.Memory.Modules, decodeSubData[rawMemory](sr.Data, exp, "Memory"))
	}
	if sr, ok := hw.SubResources["MemoryList"]; ok {
		exp.Memory.Boards = mergeBoards(exp.Memory.Boards, decodeSubData[MemoryBoard](sr.Data, exp, "MemoryList"))
	}
	if sr, ok := hw.SubResources["Devices"]; ok {
		exp.Devices = mergeDevices(exp.Devices, decodeSubData[rawDevice](sr.Data, exp, "Devices"))
	}
	if sr, ok := hw.SubResources["Processors"]; ok {
		socks := decodeSubData[rawProcessor](sr.Data, exp, "Processors")
		if len(socks) > 0 {
			exp.Processors = processorsFrom(hw, socks)
		}
	}
	if sr, ok := hw.SubResources["LocalStorage"]; ok && len(exp.Storage.Controllers) == 0 {
		applyLocalStorage(exp, decodeSubData[rawController](sr.Data, exp, "LocalStorage"), "subResources")
	}
}

func decodeSubData[T any](raw json.RawMessage, exp *ServerExport, name string) []T {
	if len(raw) == 0 {
		return nil
	}
	var list JSONList[T]
	if err := json.Unmarshal(raw, &list); err == nil && len(list) > 0 {
		return list
	}
	var wrap struct {
		Data    JSONList[T] `json:"data"`
		Members JSONList[T] `json:"members"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if len(wrap.Data) > 0 {
			return wrap.Data
		}
		if len(wrap.Members) > 0 {
			return wrap.Members
		}
	}
	if exp != nil && len(raw) > 2 {
		exp.warn("%s: unexpected payload shape", name)
	}
	return nil
}

func applyLocalStorage(exp *ServerExport, ctrls []rawController, source string) {
	if len(ctrls) == 0 {
		return
	}
	exp.Storage.Source = source
	exp.Storage.Controllers = nil
	exp.Storage.Drives = nil
	exp.Storage.Volumes = nil
	for _, c := range ctrls {
		exp.Storage.Controllers = append(exp.Storage.Controllers, StorageController{
			Name:                 c.Name,
			Model:                c.Model,
			SerialNumber:         c.SerialNumber,
			AdapterType:          c.AdapterType,
			CurrentOperatingMode: c.CurrentOperatingMode,
			Location:             c.Location,
			FirmwareVersion:      c.FirmwareVersion.Current.VersionString,
			CacheMemorySizeMiB:   c.CacheMemorySizeMiB,
			Health:               c.Status.Health,
			State:                c.Status.State,
		})
		for _, d := range c.PhysicalDrives {
			exp.Storage.Drives = append(exp.Storage.Drives, PhysicalDrive{
				Location:        d.Location,
				Model:           d.Model,
				SerialNumber:    d.SerialNumber,
				MediaType:       d.MediaType,
				InterfaceType:   d.InterfaceType,
				CapacityMiB:     d.CapacityMiB,
				CapacityGB:      d.CapacityGB,
				DiskDriveUse:    d.DiskDriveUse,
				FirmwareVersion: d.FirmwareVersion.Current.VersionString,
				Health:          d.Status.Health,
				State:           d.Status.State,
				Encrypted:       d.EncryptedDrive,
			})
		}
		for _, v := range c.LogicalDrives {
			exp.Storage.Volumes = append(exp.Storage.Volumes, LogicalDrive{
				Name:           v.LogicalDriveName,
				Number:         v.LogicalDriveNumber,
				RAID:           fmt.Sprint(v.RAID),
				MediaType:      v.MediaType,
				InterfaceType:  v.InterfaceType,
				CapacityMiB:    v.CapacityMiB,
				Health:         v.Status.Health,
				State:          v.Status.State,
				VolumeUniqueID: v.VolumeUniqueIdentifier,
				Acceleration:   v.AccelerationMethod,
			})
		}
	}
}

func mergeMemory(dst []MemoryModule, raw []rawMemory) []MemoryModule {
	if len(raw) == 0 {
		return dst
	}
	out := make([]MemoryModule, 0, len(raw))
	for _, m := range raw {
		out = append(out, MemoryModule{
			Name:              m.Name,
			DeviceLocator:     m.DeviceLocator,
			Manufacturer:      m.Manufacturer,
			PartNumber:        m.PartNumber,
			SerialNumber:      m.SerialNumber,
			CapacityMiB:       m.CapacityMiB.Int(),
			OperatingSpeedMHz: m.OperatingSpeedMhz.Int(),
			BaseModuleType:    m.BaseModuleType,
			MemoryDeviceType:  m.MemoryDeviceType,
			MemoryType:        m.MemoryType,
			RankCount:         m.RankCount.Int(),
			ErrorCorrection:   m.ErrorCorrection,
			DIMMStatus:        m.Oem.Hpe.DIMMStatus,
			Health:            m.Status.Health,
			State:             m.Status.State,
			Socket:            m.MemoryLocation.Socket.Int(),
			Controller:        m.MemoryLocation.MemoryController.Int(),
			Channel:           m.MemoryLocation.Channel.Int(),
			Slot:              m.MemoryLocation.Slot.Int(),
			Attributes:        m.Oem.Hpe.Attributes,
		})
	}
	return out
}

func mergeBoards(dst []MemoryBoard, boards []MemoryBoard) []MemoryBoard {
	if len(boards) == 0 {
		return dst
	}
	return boards
}

func mergeDevices(dst []PCIDevice, raw []rawDevice) []PCIDevice {
	if len(raw) == 0 {
		return dst
	}
	out := make([]PCIDevice, 0, len(raw))
	for _, d := range raw {
		out = append(out, PCIDevice{
			ID:                d.ID,
			Name:              d.Name,
			DeviceType:        d.DeviceType,
			Location:          d.Location,
			Manufacturer:      d.Manufacturer,
			PartNumber:        d.PartNumber,
			ProductPartNumber: d.ProductPartNumber,
			SerialNumber:      d.SerialNumber,
			FirmwareVersion:   d.FirmwareVersion.Current.VersionString,
			Health:            d.Status.Health,
			State:             d.Status.State,
		})
	}
	return out
}

func networkPortsFrom(pm PortMap) []NetworkPortExport {
	var out []NetworkPortExport
	for _, slot := range pm.DeviceSlots {
		for _, p := range slot.PhysicalPorts {
			row := NetworkPortExport{
				DeviceName:   slot.DeviceName,
				DeviceSlot:   slot.Location,
				PortNumber:   p.PortNumber,
				Type:         p.Type,
				MAC:          p.MAC,
				WWN:          p.WWN,
				Interconnect: p.InterconnectURI,
			}
			for _, v := range p.VirtualPorts {
				if v.MAC != "" {
					row.VirtualMACs = append(row.VirtualMACs, v.MAC)
				}
			}
			out = append(out, row)
		}
	}
	return out
}

func parseFirmwareBody(raw json.RawMessage) []FirmwareComponent {
	var wrap struct {
		Components []struct {
			ComponentName     string `json:"componentName"`
			ComponentVersion  string `json:"componentVersion"`
			ComponentLocation string `json:"componentLocation"`
			Name              string `json:"name"`
			Version           string `json:"version"`
			Location          string `json:"location"`
			State             string `json:"state"`
			Status            string `json:"status"`
		} `json:"components"`
		Members JSONList[ServerFirmware] `json:"members"`
	}
	if json.Unmarshal(raw, &wrap) != nil {
		return nil
	}
	var out []FirmwareComponent
	for _, c := range wrap.Components {
		out = append(out, FirmwareComponent{
			Name:     nonEmpty(c.ComponentName, c.Name),
			Version:  nonEmpty(c.ComponentVersion, c.Version),
			Location: nonEmpty(c.ComponentLocation, c.Location),
			State:    c.State,
			Status:   c.Status,
		})
	}
	for _, m := range wrap.Members {
		out = append(out, FirmwareComponent{
			Name:     m.Name,
			Version:  m.Version,
			Location: m.Location,
			State:    m.State,
			Status:   m.Status,
		})
	}
	if len(out) == 0 {
		var one ServerFirmware
		if json.Unmarshal(raw, &one) == nil && one.Name != "" {
			out = append(out, FirmwareComponent{Name: one.Name, Version: one.Version, Location: one.Location, State: one.State, Status: one.Status})
		}
	}
	return out
}

type rawProcessor struct {
	ID                string       `json:"Id"`
	Socket            string       `json:"Socket"`
	Model             string       `json:"Model"`
	Manufacturer      string       `json:"Manufacturer"`
	TotalCores        int          `json:"TotalCores"`
	TotalThreads      int          `json:"TotalThreads"`
	MaxSpeedMHz       int          `json:"MaxSpeedMHz"`
	OperatingSpeedMHz int          `json:"OperatingSpeedMHz"`
	Status            HealthStatus `json:"Status"`
}

type rawMemory struct {
	Name              string  `json:"Name"`
	DeviceLocator     string  `json:"DeviceLocator"`
	Manufacturer      string  `json:"Manufacturer"`
	PartNumber        string  `json:"PartNumber"`
	SerialNumber      string  `json:"SerialNumber"`
	CapacityMiB       FlexInt `json:"CapacityMiB"`
	OperatingSpeedMhz FlexInt `json:"OperatingSpeedMhz"`
	BaseModuleType    string  `json:"BaseModuleType"`
	MemoryDeviceType  string  `json:"MemoryDeviceType"`
	MemoryType        string  `json:"MemoryType"`
	RankCount         FlexInt `json:"RankCount"`
	ErrorCorrection   string  `json:"ErrorCorrection"`
	MemoryLocation    struct {
		Channel          FlexInt `json:"Channel"`
		MemoryController FlexInt `json:"MemoryController"`
		Slot             FlexInt `json:"Slot"`
		Socket           FlexInt `json:"Socket"`
	} `json:"MemoryLocation"`
	Oem struct {
		Hpe struct {
			DIMMStatus string   `json:"DIMMStatus"`
			Attributes []string `json:"Attributes"`
		} `json:"Hpe"`
	} `json:"Oem"`
	Status HealthStatus `json:"Status"`
}

type rawDevice struct {
	ID                string                `json:"Id"`
	Name              string                `json:"Name"`
	DeviceType        string                `json:"DeviceType"`
	Location          string                `json:"Location"`
	Manufacturer      string                `json:"Manufacturer"`
	PartNumber        string                `json:"PartNumber"`
	ProductPartNumber string                `json:"ProductPartNumber"`
	SerialNumber      string                `json:"SerialNumber"`
	FirmwareVersion   NestedFirmwareVersion `json:"FirmwareVersion"`
	Status            HealthStatus          `json:"Status"`
}

type rawController struct {
	Name                 string                `json:"Name"`
	Model                string                `json:"Model"`
	SerialNumber         string                `json:"SerialNumber"`
	AdapterType          string                `json:"AdapterType"`
	CurrentOperatingMode string                `json:"CurrentOperatingMode"`
	Location             string                `json:"Location"`
	FirmwareVersion      NestedFirmwareVersion `json:"FirmwareVersion"`
	CacheMemorySizeMiB   int                   `json:"CacheMemorySizeMiB"`
	Status               HealthStatus          `json:"Status"`
	PhysicalDrives       JSONList[rawDrive]    `json:"PhysicalDrives"`
	LogicalDrives        JSONList[rawVolume]   `json:"LogicalDrives"`
}

type rawDrive struct {
	Location        string                `json:"Location"`
	Model           string                `json:"Model"`
	SerialNumber    string                `json:"SerialNumber"`
	MediaType       string                `json:"MediaType"`
	InterfaceType   string                `json:"InterfaceType"`
	CapacityMiB     int                   `json:"CapacityMiB"`
	CapacityGB      int                   `json:"CapacityGB"`
	DiskDriveUse    string                `json:"DiskDriveUse"`
	FirmwareVersion NestedFirmwareVersion `json:"FirmwareVersion"`
	EncryptedDrive  bool                  `json:"EncryptedDrive"`
	Status          HealthStatus          `json:"Status"`
}

type rawVolume struct {
	LogicalDriveName       string       `json:"LogicalDriveName"`
	LogicalDriveNumber     int          `json:"LogicalDriveNumber"`
	RAID                   any          `json:"Raid"`
	MediaType              string       `json:"MediaType"`
	InterfaceType          string       `json:"InterfaceType"`
	CapacityMiB            int          `json:"CapacityMiB"`
	VolumeUniqueIdentifier string       `json:"VolumeUniqueIdentifier"`
	AccelerationMethod     string       `json:"AccelerationMethod"`
	Status                 HealthStatus `json:"Status"`
}
