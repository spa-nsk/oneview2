package oneview

import "encoding/json"

// URI constants from swagger 300.json and HPE OneView REST API 3800–8800.
const (
	URIResourceAlerts         = "/rest/resource-alerts"
	URIAlerts                 = "/rest/alerts"
	URIAlertSettings          = "/rest/admin-settings/alert-settings"
	URIAppliances             = "/rest/appliances"
	URIAuditLogSettings       = "/rest/audit-logs/settings"
	URIAuditLogTest           = "/rest/audit-logs/test-forwarding"
	URIAuditLogs              = "/rest/audit-logs"
	URIRemoteCertificate      = "/rest/certificates/https/remote"
	URIServerCertificates     = "/rest/certificates/servers"
	URIConvergedSystems       = "/rest/converged-systems"
	URIDatacenters            = "/rest/datacenters"
	URIDriveEnclosures        = "/rest/drive-enclosures"
	URIEnclosures             = "/rest/enclosures"
	URIGroups                 = "/rest/groups"
	URIInterconnects          = "/rest/interconnects"
	URILogicalInterconnects   = "/rest/logical-interconnects"
	URIManagedSANs            = "/rest/managed-sans"
	URINetworkInterfaces      = "/rest/appliance/network-interfaces"
	URITimeLocale             = "/rest/appliance/configuration/time-locale"
	URISANManagers            = "/rest/san-managers"
	URIServerFirmware         = "/rest/server-firmware"
	URIServerHardware         = "/rest/server-hardware"
	URIServerProfiles         = "/rest/server-profiles"
	URIServerProfileTemplates = "/rest/server-profile-templates"
	URIStoragePools           = "/rest/storage-pools"
	URIStorageSystems         = "/rest/storage-systems"
	URIStorageVolumes         = "/rest/storage-volumes"
	URIVolumes                = "/rest/volumes"

	// OneView appliance REST API 3800–8800 (not in Global Dashboard swagger).
	URIEthernetNetworks          = "/rest/ethernet-networks"
	URIFCNetworks                = "/rest/fc-networks"
	URIFCoENetworks              = "/rest/fcoe-networks"
	URINetworkSets               = "/rest/network-sets"
	URIConnectionTemplates       = "/rest/connection-templates"
	URIEnclosureGroups           = "/rest/enclosure-groups"
	URILogicalEnclosures         = "/rest/logical-enclosures"
	URILogicalInterconnectGroups = "/rest/logical-interconnect-groups"
	URIUplinkSets                = "/rest/uplink-sets"
	URIInterconnectTypes         = "/rest/interconnect-types"
	URIServerHardwareTypes       = "/rest/server-hardware-types"
	URIStorageVolumeTemplates    = "/rest/storage-volume-templates"
	URIStorageVolumeAttachments  = "/rest/storage-volume-attachments"
	URIFirmwareDrivers           = "/rest/firmware-drivers"
	URIFirmwareBundles           = "/rest/firmware-bundles"
	URIIDPools                   = "/rest/id-pools"
	URIIDPoolsIPv4Ranges         = "/rest/id-pools/ipv4/ranges"
	URIIDPoolsIPv4Subnets        = "/rest/id-pools/ipv4/subnets"
	URIRacks                     = "/rest/racks"
	URIPowerDevices              = "/rest/power-devices"
	URILabels                    = "/rest/labels"
	URIEvents                    = "/rest/events"
	URISwitches                  = "/rest/switches"
	URIFabrics                   = "/rest/fabrics"
	URIHypervisorManagers        = "/rest/hypervisor-managers"
	URIHypervisorClusterProfiles = "/rest/hypervisor-cluster-profiles"
	URIIndexResources            = "/rest/index/resources"
)

// ResourceAlert is GET /rest/resource-alerts/{id} (Global Dashboard) and
// is compatible with GET /rest/alerts/{id} on the appliance.
type ResourceAlert struct {
	Resource
	ActivityURI          string             `json:"activityUri,omitempty"`
	AlertState           string             `json:"alertState,omitempty"`
	AlertTypeID          string             `json:"alertTypeID,omitempty"`
	ApplianceModel       string             `json:"applianceModel,omitempty"`
	ApplianceVersion     string             `json:"applianceVersion,omitempty"`
	AssignedToUser       string             `json:"assignedToUser,omitempty"`
	AssociatedEventURIs  []string           `json:"associatedEventUris,omitempty"`
	AssociatedResource   AssociatedResource `json:"associatedResource,omitempty"`
	ChangeLog            []ChangeLogEntry   `json:"changeLog,omitempty"`
	ClearedByUser        string             `json:"clearedByUser,omitempty"`
	ClearedTime          string             `json:"clearedTime,omitempty"`
	CorrectiveAction     string             `json:"correctiveAction,omitempty"`
	Description          string             `json:"description,omitempty"`
	LifeCycle            bool               `json:"lifeCycle,omitempty"`
	PhysicalResourceType string             `json:"physicalResourceType,omitempty"`
	ResourceID           string             `json:"resourceID,omitempty"`
	ResourceURI          string             `json:"resourceUri,omitempty"`
	Severity             string             `json:"severity,omitempty"`
	ServiceEvent         bool               `json:"serviceEvent,omitempty"`
	ServiceEventSource   bool               `json:"serviceEventSource,omitempty"`
	Urgency              string             `json:"urgency,omitempty"`
}

// ChangeLogEntry is an alert history line.
type ChangeLogEntry struct {
	UserName       string `json:"username,omitempty"`
	OnlyForUser    bool   `json:"onlyForUser,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Timestamp      string `json:"timestamp,omitempty"`
	AlertState     string `json:"alertState,omitempty"`
	AssignedToUser string `json:"assignedToUser,omitempty"`
}

// AlertSettings is GET/PUT /rest/admin-settings/alert-settings.
type AlertSettings struct {
	Resource
	EmailFilterCategories []string `json:"emailFilterCategories,omitempty"`
	EmailFilterEnabled    bool     `json:"emailFilterEnabled,omitempty"`
	Enabled               bool     `json:"enabled,omitempty"`
	Password              string   `json:"password,omitempty"`
	Port                  int      `json:"port,omitempty"`
	SenderEmailAddress    string   `json:"senderEmailAddress,omitempty"`
	SMTPServer            string   `json:"smtpServer,omitempty"`
	Username              string   `json:"username,omitempty"`
}

// Appliance is a Global Dashboard registered OneView appliance
// (POST/GET /rest/appliances).
type Appliance struct {
	Resource
	Address           string         `json:"address,omitempty"`
	ApplianceURI      string         `json:"applianceUri,omitempty"`
	CurrentAPIVersion int            `json:"currentApiVersion,omitempty"`
	Hostname          string         `json:"hostname,omitempty"`
	LoginDomain       string         `json:"loginDomain,omitempty"`
	Model             string         `json:"model,omitempty"`
	ModuleName        string         `json:"moduleName,omitempty"`
	Password          string         `json:"password,omitempty"`
	PlatformType      string         `json:"platformType,omitempty"`
	SerialNumber      string         `json:"serialNumber,omitempty"`
	ServerCert        map[string]any `json:"serverCert,omitempty"`
	Username          string         `json:"username,omitempty"`
	Version           string         `json:"version,omitempty"`
	ForceTrusted      bool           `json:"forceTrusted,omitempty"`
}

// ApplianceAdd is POST /rest/appliances.
type ApplianceAdd struct {
	Address       string `json:"address"`
	ApplianceName string `json:"applianceName,omitempty"`
	Category      string `json:"category,omitempty"`
	LoginDomain   string `json:"loginDomain"`
	Password      string `json:"password"`
	Username      string `json:"username"`
	ForceTrusted  bool   `json:"forceTrusted,omitempty"`
}

// AuditLogSettings is GET/PUT /rest/audit-logs/settings.
type AuditLogSettings struct {
	Resource
	Enabled  bool   `json:"enabled,omitempty"`
	Facility string `json:"facility,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// ServerCertificate is POST /rest/certificates/servers.
type ServerCertificate struct {
	AliasName          string   `json:"aliasName,omitempty"`
	CertificateDetails []any    `json:"certificateDetails,omitempty"`
	Type               string   `json:"type,omitempty"`
	URI                string   `json:"uri,omitempty"`
	Certificate        string   `json:"certificate,omitempty"`
	Base64SSLCertData  string   `json:"base64SSLCertData,omitempty"`
	ServerCertBase64   []string `json:"serverCertBase64,omitempty"`
}

// RemoteCertificate is GET /rest/certificates/https/remote/{address}.
type RemoteCertificate struct {
	Type               string   `json:"type,omitempty"`
	CertificateDetails []any    `json:"certificateDetails,omitempty"`
	ServerCertBase64   []string `json:"serverCertBase64,omitempty"`
}

// Group is Global Dashboard /rest/groups.
type Group struct {
	Resource
	ParentName     string `json:"parentName,omitempty"`
	ParentURI      string `json:"parentUri,omitempty"`
	ResourceNumber int    `json:"resourceNumber,omitempty"`
}

// GroupCreate is POST /rest/groups.
type GroupCreate struct {
	Name      string `json:"name"`
	ParentURI string `json:"parentUri,omitempty"`
}

// Enclosure is GET /rest/enclosures/{id}.
type Enclosure struct {
	Resource
	ActiveOAPreferredIP  string           `json:"activeOaPreferredIP,omitempty"`
	AssetTag             string           `json:"assetTag,omitempty"`
	DeviceBayCount       int              `json:"deviceBayCount,omitempty"`
	DeviceBays           []map[string]any `json:"deviceBays,omitempty"`
	EnclosureGroupURI    string           `json:"enclosureGroupUri,omitempty"`
	EnclosureType        string           `json:"enclosureType,omitempty"`
	FanBayCount          int              `json:"fanBayCount,omitempty"`
	FanBays              []map[string]any `json:"fanBays,omitempty"`
	FwBaselineName       string           `json:"fwBaselineName,omitempty"`
	FwBaselineURI        string           `json:"fwBaselineUri,omitempty"`
	InterconnectBayCount int              `json:"interconnectBayCount,omitempty"`
	InterconnectBays     []map[string]any `json:"interconnectBays,omitempty"`
	IsFwManaged          bool             `json:"isFwManaged,omitempty"`
	LicensingIntent      string           `json:"licensingIntent,omitempty"`
	LogicalEnclosureURI  string           `json:"logicalEnclosureUri,omitempty"`
	ManagerBays          []map[string]any `json:"managerBays,omitempty"`
	OABays               []map[string]any `json:"oaBays,omitempty"`
	PartNumber           string           `json:"partNumber,omitempty"`
	PowerSupplyBayCount  int              `json:"powerSupplyBayCount,omitempty"`
	PowerSupplyBays      []map[string]any `json:"powerSupplyBays,omitempty"`
	RackName             string           `json:"rackName,omitempty"`
	RefreshState         string           `json:"refreshState,omitempty"`
	SerialNumber         string           `json:"serialNumber,omitempty"`
	StandbyOAPreferredIP string           `json:"standbyOaPreferredIP,omitempty"`
	UIDState             string           `json:"uidState,omitempty"`
	UUID                 string           `json:"uuid,omitempty"`
	Version              string           `json:"version,omitempty"`
}

// ServerHardware is GET /rest/server-hardware/{id}.
type ServerHardware struct {
	Resource
	AssetTag                   string                         `json:"assetTag,omitempty"`
	FormFactor                 string                         `json:"formFactor,omitempty"`
	HostOsType                 int                            `json:"hostOsType,omitempty"`
	LicensingIntent            string                         `json:"licensingIntent,omitempty"`
	LocationURI                string                         `json:"locationUri,omitempty"`
	MaintenanceMode            bool                           `json:"maintenanceMode,omitempty"`
	MemoryMb                   int                            `json:"memoryMb,omitempty"`
	Model                      string                         `json:"model,omitempty"`
	MpDnsName                  string                         `json:"mpDnsName,omitempty"`
	MpFirmwareVersion          string                         `json:"mpFirmwareVersion,omitempty"`
	MpHostInfo                 MpHostInfo                     `json:"mpHostInfo,omitempty"`
	MpIpAddress                string                         `json:"mpIpAddress,omitempty"`
	MpLicenseType              string                         `json:"mpLicenseType,omitempty"`
	MpModel                    string                         `json:"mpModel,omitempty"`
	MpState                    string                         `json:"mpState,omitempty"`
	OperatingSystem            string                         `json:"operatingSystem,omitempty"`
	PartNumber                 string                         `json:"partNumber,omitempty"`
	PhysicalServerHardwareURI  string                         `json:"physicalServerHardwareUri,omitempty"`
	Platform                   string                         `json:"platform,omitempty"`
	PortMap                    PortMap                        `json:"portMap,omitempty"`
	SubResources               map[string]HardwareSubResource `json:"subResources,omitempty"`
	Position                   int                            `json:"position,omitempty"`
	PowerLock                  string                         `json:"powerLock,omitempty"`
	PowerState                 string                         `json:"powerState,omitempty"`
	ProcessorCoreCount         int                            `json:"processorCoreCount,omitempty"`
	ProcessorCount             int                            `json:"processorCount,omitempty"`
	ProcessorSpeedMhz          int                            `json:"processorSpeedMhz,omitempty"`
	ProcessorType              string                         `json:"processorType,omitempty"`
	RefreshState               string                         `json:"refreshState,omitempty"`
	RemoteSupportSettings      map[string]any                 `json:"remoteSupportSettings,omitempty"`
	RomVersion                 string                         `json:"romVersion,omitempty"`
	SerialNumber               string                         `json:"serialNumber,omitempty"`
	ServerGroupURI             string                         `json:"serverGroupUri,omitempty"`
	ServerHardwareTypeName     string                         `json:"serverHardwareTypeName,omitempty"`
	ServerHardwareTypeURI      string                         `json:"serverHardwareTypeUri,omitempty"`
	ServerName                 string                         `json:"serverName,omitempty"`
	ServerProfileURI           string                         `json:"serverProfileUri,omitempty"`
	ShortModel                 string                         `json:"shortModel,omitempty"`
	SupportDataCollectionState string                         `json:"supportDataCollectionState,omitempty"`
	SupportDataCollectionType  string                         `json:"supportDataCollectionType,omitempty"`
	SupportState               string                         `json:"supportState,omitempty"`
	UIDState                   string                         `json:"uidState,omitempty"`
	UUID                       string                         `json:"uuid,omitempty"`
	VirtualSerialNumber        string                         `json:"virtualSerialNumber,omitempty"`
	VirtualUUID                string                         `json:"virtualUuid,omitempty"`
}

// HardwareSubResource is one iLO inventory collection on GET /rest/server-hardware/{id}
// (Memory, LocalStorage, Devices, AdvancedMemoryProtection, MemoryList, …).
type HardwareSubResource struct {
	URI             string          `json:"uri,omitempty"`
	CollectionState string          `json:"collectionState,omitempty"`
	Modified        string          `json:"modified,omitempty"`
	Name            string          `json:"name,omitempty"`
	Count           int             `json:"count,omitempty"`
	ETag            FlexString      `json:"etag,omitempty"`
	Data            json.RawMessage `json:"data,omitempty"`
}

// PortMap is the adapter/slot map of a server.
type PortMap struct {
	DeviceSlots []DeviceSlot `json:"deviceSlots,omitempty"`
}

// DeviceSlot is one adapter location (LOM, Mezz, PCI, …).
type DeviceSlot struct {
	DeviceName    string         `json:"deviceName,omitempty"`
	DeviceNumber  int            `json:"deviceNumber,omitempty"`
	Location      string         `json:"location,omitempty"`
	SlotNumber    int            `json:"slotNumber,omitempty"`
	PhysicalPorts []PhysicalPort `json:"physicalPorts,omitempty"`
}

// PhysicalPort is a NIC/HBA port on an adapter.
type PhysicalPort struct {
	InterconnectPort         int           `json:"interconnectPort,omitempty"`
	InterconnectURI          string        `json:"interconnectUri,omitempty"`
	MAC                      string        `json:"mac,omitempty"`
	PhysicalInterconnectPort int           `json:"physicalInterconnectPort,omitempty"`
	PhysicalInterconnectURI  string        `json:"physicalInterconnectUri,omitempty"`
	PortNumber               int           `json:"portNumber,omitempty"`
	Type                     string        `json:"type,omitempty"`
	WWN                      string        `json:"wwn,omitempty"`
	VirtualPorts             []VirtualPort `json:"virtualPorts,omitempty"`
}

// VirtualPort is a FlexNIC on a physical adapter port.
type VirtualPort struct {
	CurrentAllocatedVirtualFunctionCount int    `json:"currentAllocatedVirtualFunctionCount,omitempty"`
	MAC                                  string `json:"mac,omitempty"`
	PortFunction                         string `json:"portFunction,omitempty"`
	PortNumber                           int    `json:"portNumber,omitempty"`
	WWNN                                 string `json:"wwnn,omitempty"`
	WWPN                                 string `json:"wwpn,omitempty"`
}

// MpHostInfo is iLO hostname / IP information.
type MpHostInfo struct {
	MpHostName    string        `json:"mpHostName,omitempty"`
	MpIpAddresses []MpIPAddress `json:"mpIpAddresses,omitempty"`
}

// MpIPAddress is a management processor address.
type MpIPAddress struct {
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

// PowerStateRequest is PUT /rest/server-hardware/{id}/powerState.
type PowerStateRequest struct {
	PowerState   string `json:"powerState"`
	PowerControl string `json:"powerControl,omitempty"`
}

// SSOURL is returned by iLO / OA SSO endpoints.
type SSOURL struct {
	SSOUrl    string `json:"ssoUrl,omitempty"`
	IloSsoUrl string `json:"iloSsoUrl,omitempty"`
	Type      string `json:"type,omitempty"`
}

// ServerProfile is GET /rest/server-profiles/{id}.
type ServerProfile struct {
	Resource
	Affinity                 string           `json:"affinity,omitempty"`
	AssociatedServer         string           `json:"associatedServer,omitempty"`
	Bios                     map[string]any   `json:"bios,omitempty"`
	Boot                     map[string]any   `json:"boot,omitempty"`
	BootMode                 map[string]any   `json:"bootMode,omitempty"`
	ConnectionSettings       map[string]any   `json:"connectionSettings,omitempty"`
	Connections              []map[string]any `json:"connections,omitempty"`
	EnclosureBay             int              `json:"enclosureBay,omitempty"`
	EnclosureGroupURI        string           `json:"enclosureGroupUri,omitempty"`
	EnclosureURI             string           `json:"enclosureUri,omitempty"`
	Firmware                 map[string]any   `json:"firmware,omitempty"`
	HideUnusedFlexNics       bool             `json:"hideUnusedFlexNics,omitempty"`
	IscsiInitiatorName       string           `json:"iscsiInitiatorName,omitempty"`
	IscsiInitiatorNameType   string           `json:"iscsiInitiatorNameType,omitempty"`
	LocalStorage             map[string]any   `json:"localStorage,omitempty"`
	MacType                  string           `json:"macType,omitempty"`
	ProfileName              string           `json:"profileName,omitempty"`
	RefreshState             string           `json:"refreshState,omitempty"`
	SanStorage               map[string]any   `json:"sanStorage,omitempty"`
	SerialNumber             string           `json:"serialNumber,omitempty"`
	SerialNumberType         string           `json:"serialNumberType,omitempty"`
	ServerHardwareTypeURI    string           `json:"serverHardwareTypeUri,omitempty"`
	ServerHardwareURI        string           `json:"serverHardwareUri,omitempty"`
	ServerProfileTemplateURI string           `json:"serverProfileTemplateUri,omitempty"`
	TemplateCompliance       string           `json:"templateCompliance,omitempty"`
	UUID                     string           `json:"uuid,omitempty"`
	WwnType                  string           `json:"wwnType,omitempty"`
}

// ServerProfileTemplate is GET /rest/server-profile-templates/{id}.
type ServerProfileTemplate struct {
	Resource
	Affinity              string           `json:"affinity,omitempty"`
	Bios                  map[string]any   `json:"bios,omitempty"`
	Boot                  map[string]any   `json:"boot,omitempty"`
	BootMode              map[string]any   `json:"bootMode,omitempty"`
	ConnectionSettings    map[string]any   `json:"connectionSettings,omitempty"`
	Connections           []map[string]any `json:"connections,omitempty"`
	EnclosureGroupURI     string           `json:"enclosureGroupUri,omitempty"`
	Firmware              map[string]any   `json:"firmware,omitempty"`
	HideUnusedFlexNics    bool             `json:"hideUnusedFlexNics,omitempty"`
	LocalStorage          map[string]any   `json:"localStorage,omitempty"`
	MacType               string           `json:"macType,omitempty"`
	SanStorage            map[string]any   `json:"sanStorage,omitempty"`
	SerialNumberType      string           `json:"serialNumberType,omitempty"`
	ServerHardwareTypeURI string           `json:"serverHardwareTypeUri,omitempty"`
	ServerProfileReapply  map[string]any   `json:"serverProfileReapplyEnabled,omitempty"`
	WwnType               string           `json:"wwnType,omitempty"`
}

// Interconnect is GET /rest/interconnects/{id}.
type Interconnect struct {
	Resource
	EnclosureName          string           `json:"enclosureName,omitempty"`
	EnclosureURI           string           `json:"enclosureUri,omitempty"`
	EnclosureType          string           `json:"enclosureType,omitempty"`
	HostName               string           `json:"hostName,omitempty"`
	InterconnectIP         string           `json:"interconnectIP,omitempty"`
	InterconnectLocation   map[string]any   `json:"interconnectLocation,omitempty"`
	InterconnectTypeURI    string           `json:"interconnectTypeUri,omitempty"`
	LogicalInterconnectURI string           `json:"logicalInterconnectUri,omitempty"`
	Model                  string           `json:"model,omitempty"`
	PartNumber             string           `json:"partNumber,omitempty"`
	Ports                  []map[string]any `json:"ports,omitempty"`
	PowerState             string           `json:"powerState,omitempty"`
	ProductName            string           `json:"productName,omitempty"`
	SerialNumber           string           `json:"serialNumber,omitempty"`
}

// LogicalInterconnect is GET /rest/logical-interconnects/{id}.
type LogicalInterconnect struct {
	Resource
	ConsistencyStatus           string           `json:"consistencyStatus,omitempty"`
	DomainURI                   string           `json:"domainUri,omitempty"`
	EnclosureURIs               []string         `json:"enclosureUris,omitempty"`
	EthernetSettings            map[string]any   `json:"ethernetSettings,omitempty"`
	FabricURI                   string           `json:"fabricUri,omitempty"`
	InterconnectMap             map[string]any   `json:"interconnectMap,omitempty"`
	Interconnects               []string         `json:"interconnects,omitempty"`
	LogicalInterconnectGroupURI string           `json:"logicalInterconnectGroupUri,omitempty"`
	StackingHealth              string           `json:"stackingHealth,omitempty"`
	UplinkSets                  []map[string]any `json:"uplinkSets,omitempty"`
}

// DriveEnclosure is GET /rest/drive-enclosures/{id}.
type DriveEnclosure struct {
	Resource
	BayNumber      int              `json:"bayNumber,omitempty"`
	DeviceBays     []map[string]any `json:"deviceBays,omitempty"`
	DriveBays      []map[string]any `json:"driveBays,omitempty"`
	EnclosureURI   string           `json:"enclosureUri,omitempty"`
	IOAdapterCount int              `json:"ioAdapterCount,omitempty"`
	Model          string           `json:"model,omitempty"`
	ProductName    string           `json:"productName,omitempty"`
	SerialNumber   string           `json:"serialNumber,omitempty"`
	WWN            string           `json:"wwn,omitempty"`
}

// Datacenter is GET /rest/datacenters/{id}.
type Datacenter struct {
	Resource
	Contents                []map[string]any `json:"contents,omitempty"`
	CoolingCapacity         int              `json:"coolingCapacity,omitempty"`
	CostPerWatt             float64          `json:"costPerWatt,omitempty"`
	Currency                string           `json:"currency,omitempty"`
	DefaultPowerLineVoltage int              `json:"defaultPowerLineVoltage,omitempty"`
	DeratingType            string           `json:"deratingType,omitempty"`
	Height                  int              `json:"height,omitempty"`
	Width                   int              `json:"width,omitempty"`
}

// ConvergedSystem is GET /rest/converged-systems/{id}.
type ConvergedSystem struct {
	Resource
	LicensingIntent string           `json:"licensingIntent,omitempty"`
	Model           string           `json:"model,omitempty"`
	PartNumber      string           `json:"partNumber,omitempty"`
	PrimaryOAIPv4   string           `json:"primaryOAIPv4Address,omitempty"`
	SerialNumber    string           `json:"serialNumber,omitempty"`
	Subsystems      []map[string]any `json:"subsystems,omitempty"`
}

// ManagedSAN is GET /rest/managed-sans/{id}.
type ManagedSAN struct {
	Resource
	DeviceAliases map[string]any `json:"deviceAliases,omitempty"`
	FabricURI     string         `json:"fabricUri,omitempty"`
	IsExpected    bool           `json:"isExpected,omitempty"`
	PrincipalWWN  string         `json:"principalSwitchWwn,omitempty"`
	SanPolicy     map[string]any `json:"sanPolicy,omitempty"`
	WWN           string         `json:"wwn,omitempty"`
}

// SANManager is GET /rest/san-managers/{id}.
type SANManager struct {
	Resource
	DeviceID                 string         `json:"deviceId,omitempty"`
	DeviceSpecificAttributes map[string]any `json:"deviceSpecificAttributes,omitempty"`
	IsInternal               bool           `json:"isInternal,omitempty"`
	ProviderDisplayName      string         `json:"providerDisplayName,omitempty"`
	ProviderURI              string         `json:"providerUri,omitempty"`
}

// StorageSystem is GET /rest/storage-systems/{id}.
type StorageSystem struct {
	Resource
	Credentials              map[string]any   `json:"credentials,omitempty"`
	DeviceSpecificAttributes map[string]any   `json:"deviceSpecificAttributes,omitempty"`
	Family                   string           `json:"family,omitempty"`
	Firmware                 string           `json:"firmware,omitempty"`
	Hostname                 string           `json:"hostname,omitempty"`
	ManagedDomain            string           `json:"managedDomain,omitempty"`
	Model                    string           `json:"model,omitempty"`
	Ports                    []map[string]any `json:"ports,omitempty"`
	SerialNumber             string           `json:"serialNumber,omitempty"`
	TotalCapacity            string           `json:"totalCapacity,omitempty"`
	WWN                      string           `json:"wwn,omitempty"`
}

// StoragePool is GET /rest/storage-pools/{id}.
type StoragePool struct {
	Resource
	AllocatedCapacity        string         `json:"allocatedCapacity,omitempty"`
	DeviceSpecificAttributes map[string]any `json:"deviceSpecificAttributes,omitempty"`
	FreeCapacity             string         `json:"freeCapacity,omitempty"`
	IsManaged                bool           `json:"isManaged,omitempty"`
	StorageSystemURI         string         `json:"storageSystemUri,omitempty"`
	TotalCapacity            string         `json:"totalCapacity,omitempty"`
}

// StorageVolume is GET /rest/storage-volumes/{id} or /rest/volumes/{id}.
type StorageVolume struct {
	Resource
	AllocatedCapacity        string         `json:"allocatedCapacity,omitempty"`
	DeviceSpecificAttributes map[string]any `json:"deviceSpecificAttributes,omitempty"`
	DeviceVolumeName         string         `json:"deviceVolumeName,omitempty"`
	IsPermanent              bool           `json:"isPermanent,omitempty"`
	IsShareable              bool           `json:"isShareable,omitempty"`
	ProvisionType            string         `json:"provisionType,omitempty"`
	ProvisionedCapacity      string         `json:"provisionedCapacity,omitempty"`
	RaidLevel                string         `json:"raidLevel,omitempty"`
	Shareable                bool           `json:"shareable,omitempty"`
	StoragePoolURI           string         `json:"storagePoolUri,omitempty"`
	StorageSystemURI         string         `json:"storageSystemUri,omitempty"`
	VolumeTemplateURI        string         `json:"volumeTemplateUri,omitempty"`
	WWN                      string         `json:"wwn,omitempty"`
}

// ServerFirmware is GET /rest/server-firmware/{id} (Global Dashboard).
type ServerFirmware struct {
	Resource
	BaselineName          string           `json:"baselineName,omitempty"`
	Components            []map[string]any `json:"components,omitempty"`
	Location              string           `json:"location,omitempty"`
	Version               string           `json:"version,omitempty"`
	ServerHardwareURI     string           `json:"serverHardwareUri,omitempty"`
	ServerHardwareTypeURI string           `json:"serverHardwareTypeUri,omitempty"`
	ServerName            string           `json:"serverName,omitempty"`
	ServerModel           string           `json:"serverModel,omitempty"`
}

// TimeLocale is GET /rest/appliance/configuration/time-locale.
type TimeLocale struct {
	Type              string   `json:"type,omitempty"`
	DateTime          string   `json:"dateTime,omitempty"`
	Locale            string   `json:"locale,omitempty"`
	LocaleDisplayName string   `json:"localeDisplayName,omitempty"`
	NtpServers        []string `json:"ntpServers,omitempty"`
	Timezone          string   `json:"timezone,omitempty"`
	PollingInterval   string   `json:"pollingInterval,omitempty"`
}

// NetworkInterfaces is GET /rest/appliance/network-interfaces.
type NetworkInterfaces struct {
	Type              string           `json:"type,omitempty"`
	ApplianceNetworks []map[string]any `json:"applianceNetworks,omitempty"`
	Nics              []map[string]any `json:"nics,omitempty"`
}

// EthernetNetwork is GET /rest/ethernet-networks/{id} (appliance 3800–8800).
type EthernetNetwork struct {
	Resource
	VlanId                int    `json:"vlanId,omitempty"`
	EthernetNetworkType   string `json:"ethernetNetworkType,omitempty"`
	Purpose               string `json:"purpose,omitempty"`
	SmartLink             bool   `json:"smartLink,omitempty"`
	PrivateNetwork        bool   `json:"privateNetwork,omitempty"`
	ConnectionTemplateURI string `json:"connectionTemplateUri,omitempty"`
	SubnetURI             string `json:"subnetUri,omitempty"`
}

// FCNetwork is GET /rest/fc-networks/{id}.
type FCNetwork struct {
	Resource
	FabricType              string `json:"fabricType,omitempty"`
	LinkStabilityTime       int    `json:"linkStabilityTime,omitempty"`
	AutoLoginRedistribution bool   `json:"autoLoginRedistribution,omitempty"`
	ManagedSanURI           string `json:"managedSanUri,omitempty"`
	ConnectionTemplateURI   string `json:"connectionTemplateUri,omitempty"`
}

// NetworkSet is GET /rest/network-sets/{id}.
type NetworkSet struct {
	Resource
	NativeNetworkURI      string   `json:"nativeNetworkUri,omitempty"`
	NetworkURIs           []string `json:"networkUris,omitempty"`
	ConnectionTemplateURI string   `json:"connectionTemplateUri,omitempty"`
}

// EnclosureGroup is GET /rest/enclosure-groups/{id}.
type EnclosureGroup struct {
	Resource
	EnclosureCount          int              `json:"enclosureCount,omitempty"`
	InterconnectBayMappings []map[string]any `json:"interconnectBayMappings,omitempty"`
	IpAddressingMode        string           `json:"ipAddressingMode,omitempty"`
	Ipv6AddressingMode      string           `json:"ipv6AddressingMode,omitempty"`
	PowerMode               string           `json:"powerMode,omitempty"`
	PortMappingCount        int              `json:"portMappingCount,omitempty"`
	StackingMode            string           `json:"stackingMode,omitempty"`
	OsDeploymentSettings    map[string]any   `json:"osDeploymentSettings,omitempty"`
}

// LogicalEnclosure is GET /rest/logical-enclosures/{id}.
type LogicalEnclosure struct {
	Resource
	DeleteFailed              bool           `json:"deleteFailed,omitempty"`
	DeploymentManagerSettings map[string]any `json:"deploymentManagerSettings,omitempty"`
	EnclosureGroupURI         string         `json:"enclosureGroupUri,omitempty"`
	EnclosureURIs             []string       `json:"enclosureUris,omitempty"`
	Firmware                  map[string]any `json:"firmware,omitempty"`
	LogicalInterconnectURIs   []string       `json:"logicalInterconnectUris,omitempty"`
	PowerMode                 string         `json:"powerMode,omitempty"`
	ScalingState              string         `json:"scalingState,omitempty"`
}

// LogicalInterconnectGroup is GET /rest/logical-interconnect-groups/{id}.
type LogicalInterconnectGroup struct {
	Resource
	EnclosureIndexes        []int            `json:"enclosureIndexes,omitempty"`
	EnclosureType           string           `json:"enclosureType,omitempty"`
	EthernetSettings        map[string]any   `json:"ethernetSettings,omitempty"`
	FabricURI               string           `json:"fabricUri,omitempty"`
	InterconnectMapTemplate map[string]any   `json:"interconnectMapTemplate,omitempty"`
	InternalNetworkURIs     []string         `json:"internalNetworkUris,omitempty"`
	QualityOfService        map[string]any   `json:"qualityOfService,omitempty"`
	RedundancyType          string           `json:"redundancyType,omitempty"`
	StackingHealth          string           `json:"stackingHealth,omitempty"`
	UplinkSets              []map[string]any `json:"uplinkSets,omitempty"`
}

// UplinkSet is GET /rest/uplink-sets/{id}.
type UplinkSet struct {
	Resource
	ConnectionMode                string           `json:"connectionMode,omitempty"`
	EthernetNetworkType           string           `json:"ethernetNetworkType,omitempty"`
	LacpTimer                     string           `json:"lacpTimer,omitempty"`
	LogicalInterconnectURI        string           `json:"logicalInterconnectUri,omitempty"`
	ManualLoginRedistributionType string           `json:"manualLoginRedistributionType,omitempty"`
	NativeNetworkURI              string           `json:"nativeNetworkUri,omitempty"`
	NetworkType                   string           `json:"networkType,omitempty"`
	NetworkURIs                   []string         `json:"networkUris,omitempty"`
	PortConfigInfos               []map[string]any `json:"portConfigInfos,omitempty"`
	PrimaryPort                   map[string]any   `json:"primaryPort,omitempty"`
}

// FirmwareDriver is GET /rest/firmware-drivers/{id}.
type FirmwareDriver struct {
	Resource
	BaselineShortName  string           `json:"baselineShortName,omitempty"`
	BundleSize         int              `json:"bundleSize,omitempty"`
	BundleType         string           `json:"bundleType,omitempty"`
	FwComponents       []map[string]any `json:"fwComponents,omitempty"`
	Hotfixes           []map[string]any `json:"hotfixes,omitempty"`
	ReleaseDate        string           `json:"releaseDate,omitempty"`
	ResourceId         string           `json:"resourceId,omitempty"`
	Version            string           `json:"version,omitempty"`
	SupportedLanguages string           `json:"supportedLanguages,omitempty"`
}

// ServerHardwareType is GET /rest/server-hardware-types/{id}.
type ServerHardwareType struct {
	Resource
	Adapters            []map[string]any `json:"adapters,omitempty"`
	BiosSettings        []map[string]any `json:"biosSettings,omitempty"`
	BootTargetOptions   []string         `json:"bootTargetOptions,omitempty"`
	Capabilities        []string         `json:"capabilities,omitempty"`
	FormFactor          string           `json:"formFactor,omitempty"`
	Model               string           `json:"model,omitempty"`
	PxeBootPolicy       string           `json:"pxeBootPolicy,omitempty"`
	StorageCapabilities map[string]any   `json:"storageCapabilities,omitempty"`
}
