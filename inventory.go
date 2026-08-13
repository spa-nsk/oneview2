package oneview

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExportOptions controls which extra OneView resources are pulled into a server dump.
// Zero value (ExportOptions{}) means all sections are included.
type ExportOptions struct {
	SkipFirmware     bool
	SkipLocalStorage bool
	SkipBIOS         bool
	SkipEnvironment  bool
	SkipProfile      bool
	SkipHardwareType bool
	SkipEnclosure    bool
	SkipSubresources bool // do not GET /memory, /devices, /processors, …
}

// ServerExport is a high-level dump of one physical server: identity, CPUs,
// DIMMs, local disks/RAID, PCI devices, NICs, firmware, BIOS and profile.
type ServerExport struct {
	CollectedAt  time.Time           `json:"collectedAt"`
	APIVersion   int                 `json:"apiVersion,omitempty"`
	Product      string              `json:"product,omitempty"`
	Identity     ServerIdentity      `json:"identity"`
	Processors   ProcessorInventory  `json:"processors"`
	Memory       MemoryInventory     `json:"memory"`
	Storage      StorageInventory    `json:"storage"`
	Devices      []PCIDevice         `json:"devices,omitempty"`
	NetworkPorts []NetworkPortExport `json:"networkPorts,omitempty"`
	Firmware     []FirmwareComponent `json:"firmware,omitempty"`
	BIOS         json.RawMessage     `json:"bios,omitempty"`
	Environment  json.RawMessage     `json:"environment,omitempty"`
	Profile      *ServerProfile      `json:"profile,omitempty"`
	HardwareType *ServerHardwareType `json:"hardwareType,omitempty"`
	Enclosure    *Enclosure          `json:"enclosure,omitempty"`
	Hardware     *ServerHardware     `json:"hardware,omitempty"`
	Warnings     []string            `json:"warnings,omitempty"`
}

// ServerIdentity is the inventory header for a server.
type ServerIdentity struct {
	Name                   string `json:"name,omitempty"`
	ServerName             string `json:"serverName,omitempty"`
	URI                    string `json:"uri,omitempty"`
	UUID                   string `json:"uuid,omitempty"`
	SerialNumber           string `json:"serialNumber,omitempty"`
	Model                  string `json:"model,omitempty"`
	ShortModel             string `json:"shortModel,omitempty"`
	PartNumber             string `json:"partNumber,omitempty"`
	AssetTag               string `json:"assetTag,omitempty"`
	FormFactor             string `json:"formFactor,omitempty"`
	Platform               string `json:"platform,omitempty"`
	PowerState             string `json:"powerState,omitempty"`
	Status                 string `json:"status,omitempty"`
	State                  string `json:"state,omitempty"`
	OperatingSystem        string `json:"operatingSystem,omitempty"`
	RomVersion             string `json:"romVersion,omitempty"`
	MpModel                string `json:"mpModel,omitempty"`
	MpFirmwareVersion      string `json:"mpFirmwareVersion,omitempty"`
	MpIPAddress            string `json:"mpIpAddress,omitempty"`
	MpHostName             string `json:"mpHostName,omitempty"`
	Position               int    `json:"position,omitempty"`
	LocationURI            string `json:"locationUri,omitempty"`
	ServerProfileURI       string `json:"serverProfileUri,omitempty"`
	ServerHardwareTypeURI  string `json:"serverHardwareTypeUri,omitempty"`
	ServerHardwareTypeName string `json:"serverHardwareTypeName,omitempty"`
	ApplianceLocation      string `json:"applianceLocation,omitempty"`
	ApplianceName          string `json:"applianceName,omitempty"`
}

// ProcessorInventory is CPU summary plus per-socket details when OneView has them.
type ProcessorInventory struct {
	Count       int             `json:"count"`
	CoresPerCPU int             `json:"coresPerCpu,omitempty"`
	TotalCores  int             `json:"totalCores,omitempty"`
	SpeedMHz    int             `json:"speedMHz,omitempty"`
	Model       string          `json:"model,omitempty"`
	Sockets     []ProcessorInfo `json:"sockets,omitempty"`
}

// ProcessorInfo is one CPU socket (Redfish/iLO inventory when present).
type ProcessorInfo struct {
	ID                string `json:"id,omitempty"`
	Socket            string `json:"socket,omitempty"`
	Model             string `json:"model,omitempty"`
	Manufacturer      string `json:"manufacturer,omitempty"`
	TotalCores        int    `json:"totalCores,omitempty"`
	TotalThreads      int    `json:"totalThreads,omitempty"`
	MaxSpeedMHz       int    `json:"maxSpeedMHz,omitempty"`
	OperatingSpeedMHz int    `json:"operatingSpeedMHz,omitempty"`
	Health            string `json:"health,omitempty"`
	State             string `json:"state,omitempty"`
}

// MemoryInventory is DIMM-level and per-board memory.
type MemoryInventory struct {
	TotalMiB int            `json:"totalMiB,omitempty"`
	Modules  []MemoryModule `json:"modules,omitempty"`
	Boards   []MemoryBoard  `json:"boards,omitempty"`
}

// MemoryModule is one DIMM from subResources.Memory.
type MemoryModule struct {
	Name              string   `json:"name,omitempty"`
	DeviceLocator     string   `json:"deviceLocator,omitempty"`
	Manufacturer      string   `json:"manufacturer,omitempty"`
	PartNumber        string   `json:"partNumber,omitempty"`
	SerialNumber      string   `json:"serialNumber,omitempty"`
	CapacityMiB       int      `json:"capacityMiB,omitempty"`
	OperatingSpeedMHz int      `json:"operatingSpeedMHz,omitempty"`
	BaseModuleType    string   `json:"baseModuleType,omitempty"`
	MemoryDeviceType  string   `json:"memoryDeviceType,omitempty"`
	MemoryType        string   `json:"memoryType,omitempty"`
	RankCount         int      `json:"rankCount,omitempty"`
	ErrorCorrection   string   `json:"errorCorrection,omitempty"`
	DIMMStatus        string   `json:"dimmStatus,omitempty"`
	Health            string   `json:"health,omitempty"`
	State             string   `json:"state,omitempty"`
	Socket            int      `json:"socket,omitempty"`
	Controller        int      `json:"memoryController,omitempty"`
	Channel           int      `json:"channel,omitempty"`
	Slot              int      `json:"slot,omitempty"`
	Attributes        []string `json:"attributes,omitempty"`
}

// MemoryBoard is subResources.MemoryList (per-CPU board totals).
type MemoryBoard struct {
	CPUNumber            int `json:"boardCpuNumber,omitempty"`
	NumberOfSockets      int `json:"boardNumberOfSockets,omitempty"`
	OperationalFrequency int `json:"boardOperationalFrequency,omitempty"`
	OperationalVoltage   int `json:"boardOperationalVoltage,omitempty"`
	TotalMemorySizeMiB   int `json:"boardTotalMemorySize,omitempty"`
}

// StorageInventory is RAID controllers, logical volumes and physical disks.
type StorageInventory struct {
	Controllers []StorageController `json:"controllers,omitempty"`
	Drives      []PhysicalDrive     `json:"drives,omitempty"`
	Volumes     []LogicalDrive      `json:"volumes,omitempty"`
	Source      string              `json:"source,omitempty"` // localStorage | localStorageV2 | subResources
}

// StorageController is a Smart Array / HBA from LocalStorage inventory.
type StorageController struct {
	Name                 string `json:"name,omitempty"`
	Model                string `json:"model,omitempty"`
	SerialNumber         string `json:"serialNumber,omitempty"`
	AdapterType          string `json:"adapterType,omitempty"`
	CurrentOperatingMode string `json:"currentOperatingMode,omitempty"`
	Location             string `json:"location,omitempty"`
	FirmwareVersion      string `json:"firmwareVersion,omitempty"`
	CacheMemorySizeMiB   int    `json:"cacheMemorySizeMiB,omitempty"`
	Health               string `json:"health,omitempty"`
	State                string `json:"state,omitempty"`
}

// PhysicalDrive is a HDD/SSD/NVMe behind a controller.
type PhysicalDrive struct {
	Location        string `json:"location,omitempty"`
	Model           string `json:"model,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	MediaType       string `json:"mediaType,omitempty"`
	InterfaceType   string `json:"interfaceType,omitempty"`
	CapacityMiB     int    `json:"capacityMiB,omitempty"`
	CapacityGB      int    `json:"capacityGB,omitempty"`
	DiskDriveUse    string `json:"diskDriveUse,omitempty"`
	FirmwareVersion string `json:"firmwareVersion,omitempty"`
	Health          string `json:"health,omitempty"`
	State           string `json:"state,omitempty"`
	Encrypted       bool   `json:"encrypted,omitempty"`
}

// LogicalDrive is a RAID volume on a controller.
type LogicalDrive struct {
	Name           string `json:"name,omitempty"`
	Number         int    `json:"number,omitempty"`
	RAID           string `json:"raid,omitempty"`
	MediaType      string `json:"mediaType,omitempty"`
	InterfaceType  string `json:"interfaceType,omitempty"`
	CapacityMiB    int    `json:"capacityMiB,omitempty"`
	Health         string `json:"health,omitempty"`
	State          string `json:"state,omitempty"`
	VolumeUniqueID string `json:"volumeUniqueIdentifier,omitempty"`
	Acceleration   string `json:"accelerationMethod,omitempty"`
}

// PCIDevice is a Devices subresource entry (NIC, GPU, HBA, backplane, …).
type PCIDevice struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	DeviceType        string `json:"deviceType,omitempty"`
	Location          string `json:"location,omitempty"`
	Manufacturer      string `json:"manufacturer,omitempty"`
	PartNumber        string `json:"partNumber,omitempty"`
	ProductPartNumber string `json:"productPartNumber,omitempty"`
	SerialNumber      string `json:"serialNumber,omitempty"`
	FirmwareVersion   string `json:"firmwareVersion,omitempty"`
	Health            string `json:"health,omitempty"`
	State             string `json:"state,omitempty"`
}

// NetworkPortExport is a physical adapter port from portMap.
type NetworkPortExport struct {
	DeviceName   string   `json:"deviceName,omitempty"`
	DeviceSlot   string   `json:"deviceSlot,omitempty"`
	PortNumber   int      `json:"portNumber,omitempty"`
	Type         string   `json:"type,omitempty"`
	MAC          string   `json:"mac,omitempty"`
	WWN          string   `json:"wwn,omitempty"`
	VirtualMACs  []string `json:"virtualMacs,omitempty"`
	Interconnect string   `json:"interconnectUri,omitempty"`
}

// FirmwareComponent is one firmware/driver inventory row.
type FirmwareComponent struct {
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
	Location string `json:"location,omitempty"`
	State    string `json:"state,omitempty"`
	Status   string `json:"status,omitempty"`
}

// HealthStatus is the Redfish {Health, State} pair.
type HealthStatus struct {
	Health string `json:"Health,omitempty"`
	State  string `json:"State,omitempty"`
}

// NestedFirmwareVersion is {Current:{VersionString}}.
type NestedFirmwareVersion struct {
	Current struct {
		VersionString string `json:"VersionString,omitempty"`
	} `json:"Current,omitempty"`
}

func (e *ServerExport) warn(format string, args ...any) {
	e.Warnings = append(e.Warnings, fmt.Sprintf(format, args...))
}

// Summary returns a short human-readable inventory of the server.
func (e *ServerExport) Summary() string {
	if e == nil {
		return ""
	}
	var b strings.Builder
	id := e.Identity
	fmt.Fprintf(&b, "%s  model=%s  serial=%s  power=%s  status=%s\n",
		nonEmpty(id.Name, id.ServerName), id.ShortModel, id.SerialNumber, id.PowerState, id.Status)
	fmt.Fprintf(&b, "  CPU: %d x %s  cores/cpu=%d  speed=%d MHz\n",
		e.Processors.Count, e.Processors.Model, e.Processors.CoresPerCPU, e.Processors.SpeedMHz)
	fmt.Fprintf(&b, "  RAM: %d MiB  (%d DIMMs)\n", e.Memory.TotalMiB, len(e.Memory.Modules))
	fmt.Fprintf(&b, "  Disks: %d physical  %d logical  %d controllers\n",
		len(e.Storage.Drives), len(e.Storage.Volumes), len(e.Storage.Controllers))
	fmt.Fprintf(&b, "  Devices: %d   NIC ports: %d   Firmware components: %d\n",
		len(e.Devices), len(e.NetworkPorts), len(e.Firmware))
	if len(e.Warnings) > 0 {
		fmt.Fprintf(&b, "  warnings: %d\n", len(e.Warnings))
	}
	return b.String()
}

func nonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
