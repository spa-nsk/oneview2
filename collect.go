package oneview

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Server is a version-agnostic inventory of one physical machine collected from
// Global Dashboard (API 300) and/or appliance APIs 3800–8800.
//
// Arrays are always non-nil: missing data from an older or thinner API is an
// empty slice, not null. Local disks live on Controllers (not a flat list).
type Server struct {
	Identity     ServerIdentity      `json:"identity"`
	Processors   []ProcessorInfo     `json:"processors"`
	Memory       []MemoryModule      `json:"memory"`
	Controllers  []StorageController `json:"controllers"`
	Devices      []PCIDevice         `json:"devices"`
	NetworkPorts []NetworkPortExport `json:"networkPorts"`
	Firmware     []FirmwareComponent `json:"firmware"`
	APIVersion   int                 `json:"apiVersion,omitempty"`
	Product      string              `json:"product,omitempty"`
	Endpoint     string              `json:"endpoint,omitempty"`
	Sources      []ServerSource      `json:"sources"`
	Warnings     []string            `json:"warnings"`
}

// ServerSource records which OneView endpoint contributed a Server.
type ServerSource struct {
	Endpoint   string `json:"endpoint,omitempty"`
	APIVersion int    `json:"apiVersion,omitempty"`
	Product    string `json:"product,omitempty"`
}

// CollectServers logs into every Config, dumps all servers, and returns a
// deduplicated list. The same physical machine seen on Global Dashboard and on
// an appliance (or on two API versions) is merged: empty fields and empty
// arrays are filled from the richer payload.
//
// A failed appliance does not discard servers already collected; those errors
// are joined and returned alongside the list.
func CollectServers(ctx context.Context, configs []Config) ([]Server, error) {
	out := make([]Server, 0)
	index := map[string]int{}
	var errs []error

	for i, cfg := range configs {
		if err := ctx.Err(); err != nil {
			return out, errors.Join(append(errs, err)...)
		}
		servers, err := collectFromConfig(ctx, cfg)
		if err != nil {
			label := cfg.Endpoint
			if label == "" {
				label = fmt.Sprintf("config[%d]", i)
			}
			errs = append(errs, fmt.Errorf("oneview: %s: %w", label, err))
		}
		for _, s := range servers {
			addOrMergeServer(&out, index, s)
		}
	}

	return out, errors.Join(errs...)
}

func collectFromConfig(ctx context.Context, cfg Config) ([]Server, error) {
	c, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if err := c.Login(ctx); err != nil {
		return nil, err
	}
	defer c.Logout(ctx)

	exports, err := c.ExportServers(ctx, ListOptions{Count: -1}, ExportOptions{})
	if err != nil {
		return nil, err
	}
	src := ServerSource{
		Endpoint:   c.BaseURL(),
		APIVersion: c.APIVersion(),
		Product:    c.Product().String(),
	}
	out := make([]Server, 0, len(exports))
	seen := map[string]int{}
	for _, exp := range exports {
		if exp == nil {
			continue
		}
		s := serverFromExport(exp, src)
		addOrMergeServer(&out, seen, s)
	}
	return out, nil
}

func serverFromExport(exp *ServerExport, src ServerSource) Server {
	s := Server{
		Identity:     exp.Identity,
		Processors:   copyProcessors(exp.Processors.Sockets),
		Memory:       copyMemory(exp.Memory.Modules),
		Controllers:  copyControllers(exp.Storage.Controllers),
		Devices:      copyDevices(exp.Devices),
		NetworkPorts: copyPorts(exp.NetworkPorts),
		Firmware:     copyFirmware(exp.Firmware),
		APIVersion:   exp.APIVersion,
		Product:      exp.Product,
		Endpoint:     src.Endpoint,
		Sources:      []ServerSource{src},
		Warnings:     append([]string{}, exp.Warnings...),
	}
	ensureServerSlices(&s)
	return s
}

func addOrMergeServer(out *[]Server, index map[string]int, s Server) {
	ensureServerSlices(&s)
	keys := serverKeys(s)
	for _, k := range keys {
		if i, ok := index[k]; ok {
			mergeServer(&(*out)[i], s)
			for _, nk := range serverKeys((*out)[i]) {
				index[nk] = i
			}
			return
		}
	}
	i := len(*out)
	*out = append(*out, s)
	for _, k := range keys {
		index[k] = i
	}
}

func serverKeys(s Server) []string {
	var keys []string
	if u := normalizeKey(s.Identity.UUID); u != "" {
		keys = append(keys, "uuid:"+u)
	}
	if n := normalizeKey(s.Identity.SerialNumber); n != "" {
		keys = append(keys, "serial:"+n)
	}
	if u := hardwareID(s.Identity.URI); u != "" {
		keys = append(keys, "id:"+u)
	}
	return keys
}

func hardwareID(uri string) string {
	id := strings.ToLower(strings.TrimSpace(IDFromURI(uri)))
	if id == "" || id == "server-hardware" {
		return ""
	}
	return id
}

func normalizeKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func mergeServer(dst *Server, src Server) {
	mergeIdentity(&dst.Identity, src.Identity)
	dst.Processors = preferSlice(dst.Processors, src.Processors, processorsRicher)
	dst.Memory = preferSlice(dst.Memory, src.Memory, memoryRicher)
	dst.Controllers = preferSlice(dst.Controllers, src.Controllers, controllersRicher)
	dst.Devices = preferSlice(dst.Devices, src.Devices, devicesRicher)
	dst.NetworkPorts = preferSlice(dst.NetworkPorts, src.NetworkPorts, portsRicher)
	dst.Firmware = preferSlice(dst.Firmware, src.Firmware, firmwareRicher)
	if dst.APIVersion == 0 {
		dst.APIVersion = src.APIVersion
	}
	if dst.Product == "" {
		dst.Product = src.Product
	}
	if dst.Endpoint == "" {
		dst.Endpoint = src.Endpoint
	}
	dst.Sources = mergeSources(dst.Sources, src.Sources)
	dst.Warnings = append(dst.Warnings, src.Warnings...)
	ensureServerSlices(dst)
}

func mergeIdentity(dst *ServerIdentity, src ServerIdentity) {
	dst.Name = first(dst.Name, src.Name)
	dst.ServerName = first(dst.ServerName, src.ServerName)
	dst.URI = first(dst.URI, src.URI)
	dst.UUID = first(dst.UUID, src.UUID)
	dst.SerialNumber = first(dst.SerialNumber, src.SerialNumber)
	dst.Model = first(dst.Model, src.Model)
	dst.ShortModel = first(dst.ShortModel, src.ShortModel)
	dst.PartNumber = first(dst.PartNumber, src.PartNumber)
	dst.AssetTag = first(dst.AssetTag, src.AssetTag)
	dst.FormFactor = first(dst.FormFactor, src.FormFactor)
	dst.Platform = first(dst.Platform, src.Platform)
	dst.PowerState = first(dst.PowerState, src.PowerState)
	dst.Status = first(dst.Status, src.Status)
	dst.State = first(dst.State, src.State)
	dst.OperatingSystem = first(dst.OperatingSystem, src.OperatingSystem)
	dst.RomVersion = first(dst.RomVersion, src.RomVersion)
	dst.MpModel = first(dst.MpModel, src.MpModel)
	dst.MpFirmwareVersion = first(dst.MpFirmwareVersion, src.MpFirmwareVersion)
	dst.MpIPAddress = first(dst.MpIPAddress, src.MpIPAddress)
	dst.MpHostName = first(dst.MpHostName, src.MpHostName)
	if dst.Position == 0 {
		dst.Position = src.Position
	}
	dst.LocationURI = first(dst.LocationURI, src.LocationURI)
	dst.ServerProfileURI = first(dst.ServerProfileURI, src.ServerProfileURI)
	dst.ServerHardwareTypeURI = first(dst.ServerHardwareTypeURI, src.ServerHardwareTypeURI)
	dst.ServerHardwareTypeName = first(dst.ServerHardwareTypeName, src.ServerHardwareTypeName)
	dst.ApplianceLocation = first(dst.ApplianceLocation, src.ApplianceLocation)
	dst.ApplianceName = first(dst.ApplianceName, src.ApplianceName)
}

func preferSlice[T any](dst, src []T, richer func([]T) bool) []T {
	if len(src) == 0 {
		if dst == nil {
			return []T{}
		}
		return dst
	}
	if len(dst) == 0 {
		return src
	}
	if richer(src) && !richer(dst) {
		return src
	}
	if len(src) > len(dst) {
		return src
	}
	return dst
}

func processorsRicher(v []ProcessorInfo) bool {
	for _, p := range v {
		if p.Manufacturer != "" || p.ID != "" || p.TotalThreads > 0 {
			return true
		}
	}
	return false
}

func memoryRicher(v []MemoryModule) bool {
	for _, m := range v {
		if m.SerialNumber != "" || m.PartNumber != "" || m.DeviceLocator != "" {
			return true
		}
	}
	return false
}

func controllersRicher(v []StorageController) bool {
	for _, c := range v {
		if len(c.Drives) > 0 || c.SerialNumber != "" || c.Model != "" {
			return true
		}
	}
	return false
}

func devicesRicher(v []PCIDevice) bool {
	return len(v) > 0 && v[0].Name != ""
}

func portsRicher(v []NetworkPortExport) bool {
	for _, p := range v {
		if p.MAC != "" {
			return true
		}
	}
	return false
}

func firmwareRicher(v []FirmwareComponent) bool {
	return len(v) > 0 && v[0].Version != ""
}

func mergeSources(dst, src []ServerSource) []ServerSource {
	seen := map[string]struct{}{}
	out := make([]ServerSource, 0, len(dst)+len(src))
	for _, s := range append(dst, src...) {
		key := strings.ToLower(s.Endpoint) + "|" + fmt.Sprintf("%d", s.APIVersion)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, s)
	}
	return out
}

func ensureServerSlices(s *Server) {
	if s.Processors == nil {
		s.Processors = []ProcessorInfo{}
	}
	if s.Memory == nil {
		s.Memory = []MemoryModule{}
	}
	if s.Controllers == nil {
		s.Controllers = []StorageController{}
	}
	for i := range s.Controllers {
		if s.Controllers[i].Drives == nil {
			s.Controllers[i].Drives = []PhysicalDrive{}
		}
		if s.Controllers[i].Volumes == nil {
			s.Controllers[i].Volumes = []LogicalDrive{}
		}
	}
	if s.Devices == nil {
		s.Devices = []PCIDevice{}
	}
	if s.NetworkPorts == nil {
		s.NetworkPorts = []NetworkPortExport{}
	}
	if s.Firmware == nil {
		s.Firmware = []FirmwareComponent{}
	}
	if s.Sources == nil {
		s.Sources = []ServerSource{}
	}
	if s.Warnings == nil {
		s.Warnings = []string{}
	}
}

func copyProcessors(in []ProcessorInfo) []ProcessorInfo {
	if len(in) == 0 {
		return []ProcessorInfo{}
	}
	out := make([]ProcessorInfo, len(in))
	copy(out, in)
	return out
}

func copyMemory(in []MemoryModule) []MemoryModule {
	if len(in) == 0 {
		return []MemoryModule{}
	}
	out := make([]MemoryModule, len(in))
	copy(out, in)
	return out
}

func copyControllers(in []StorageController) []StorageController {
	if len(in) == 0 {
		return []StorageController{}
	}
	out := make([]StorageController, len(in))
	for i, c := range in {
		c.Drives = append([]PhysicalDrive{}, c.Drives...)
		if c.Drives == nil {
			c.Drives = []PhysicalDrive{}
		}
		c.Volumes = append([]LogicalDrive{}, c.Volumes...)
		if c.Volumes == nil {
			c.Volumes = []LogicalDrive{}
		}
		out[i] = c
	}
	return out
}

func copyDevices(in []PCIDevice) []PCIDevice {
	if len(in) == 0 {
		return []PCIDevice{}
	}
	out := make([]PCIDevice, len(in))
	copy(out, in)
	return out
}

func copyPorts(in []NetworkPortExport) []NetworkPortExport {
	if len(in) == 0 {
		return []NetworkPortExport{}
	}
	out := make([]NetworkPortExport, len(in))
	copy(out, in)
	return out
}

func copyFirmware(in []FirmwareComponent) []FirmwareComponent {
	if len(in) == 0 {
		return []FirmwareComponent{}
	}
	out := make([]FirmwareComponent, len(in))
	copy(out, in)
	return out
}
