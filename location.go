package oneview

import (
	"context"
	"strings"
)

// DatacenterInfo is the datacenter that contains the server's rack.
type DatacenterInfo struct {
	Name            string `json:"name,omitempty"`
	URI             string `json:"uri,omitempty"`
	UUID            string `json:"uuid,omitempty"`
	Status          string `json:"status,omitempty"`
	State           string `json:"state,omitempty"`
	CoolingCapacity int    `json:"coolingCapacity,omitempty"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	Currency        string `json:"currency,omitempty"`
}

// RackInfo is the rack that holds the enclosure or the rack-mount server.
type RackInfo struct {
	Name          string `json:"name,omitempty"`
	URI           string `json:"uri,omitempty"`
	UUID          string `json:"uuid,omitempty"`
	Model         string `json:"model,omitempty"`
	SerialNumber  string `json:"serialNumber,omitempty"`
	PartNumber    string `json:"partNumber,omitempty"`
	UHeight       int    `json:"uHeight,omitempty"`
	TopUSlot      int    `json:"topUSlot,omitempty"`
	MountUHeight  int    `json:"mountUHeight,omitempty"`
	MountLocation string `json:"mountLocation,omitempty"`
	ThermalLimit  int    `json:"thermalLimit,omitempty"`
	Status        string `json:"status,omitempty"`
	State         string `json:"state,omitempty"`
}

// EnclosureInfo is the blade chassis that holds the server, if any.
// Rack-mount servers leave this empty.
type EnclosureInfo struct {
	Name                string `json:"name,omitempty"`
	URI                 string `json:"uri,omitempty"`
	UUID                string `json:"uuid,omitempty"`
	SerialNumber        string `json:"serialNumber,omitempty"`
	Model               string `json:"model,omitempty"`
	EnclosureType       string `json:"enclosureType,omitempty"`
	PartNumber          string `json:"partNumber,omitempty"`
	RackName            string `json:"rackName,omitempty"`
	GroupURI            string `json:"enclosureGroupUri,omitempty"`
	LogicalEnclosureURI string `json:"logicalEnclosureUri,omitempty"`
	DeviceBayCount      int    `json:"deviceBayCount,omitempty"`
	Status              string `json:"status,omitempty"`
	State               string `json:"state,omitempty"`
	ActiveOAPreferredIP string `json:"activeOaPreferredIP,omitempty"`
}

// BayInfo is the device bay occupied by the server inside the enclosure.
type BayInfo struct {
	BayNumber      int    `json:"bayNumber,omitempty"`
	DevicePresence string `json:"devicePresence,omitempty"`
	DeviceURI      string `json:"deviceUri,omitempty"`
	PowerState     string `json:"powerState,omitempty"`
	DeviceBayType  string `json:"deviceBayType,omitempty"`
	Model          string `json:"model,omitempty"`
	ProfileURI     string `json:"profileUri,omitempty"`
	Covered        bool   `json:"covered,omitempty"`
}

type locationCatalog struct {
	racks []Rack
	dcs   []Datacenter
}

func (c *Client) loadLocationCatalog(ctx context.Context) locationCatalog {
	cat := locationCatalog{racks: []Rack{}, dcs: []Datacenter{}}
	if col, err := c.ListRacks(ctx, ListOptions{Count: -1}); err == nil && col != nil {
		cat.racks = col.Members
	}
	if col, err := c.ListDatacenters(ctx, ListOptions{Count: -1}); err == nil && col != nil {
		cat.dcs = col.Members
	}
	return cat
}

func fillServerLocation(s *Server, exp *ServerExport, cat locationCatalog) {
	if s == nil {
		return
	}
	var hw *ServerHardware
	var enc *Enclosure
	if exp != nil {
		hw = exp.Hardware
		enc = exp.Enclosure
	}
	s.EnclosureInfo = enclosureInfoFrom(enc, hw)
	s.BayInfo = bayInfoFrom(hw, enc)
	s.Rack, s.Datacenter = cat.rackAndDatacenter(hw, enc, s.EnclosureInfo)
}

func enclosureInfoFrom(enc *Enclosure, hw *ServerHardware) EnclosureInfo {
	info := EnclosureInfo{}
	if hw != nil {
		info.URI = hw.LocationURI
	}
	if enc == nil {
		return info
	}
	info.Name = enc.Name
	info.URI = first(enc.URI, enc.OriginalURI, info.URI)
	info.UUID = first(enc.UUID, enc.ID)
	info.SerialNumber = enc.SerialNumber
	info.Model = first(enc.EnclosureModel, enc.Name)
	info.EnclosureType = enc.EnclosureType
	info.PartNumber = enc.PartNumber
	info.RackName = enc.RackName
	info.GroupURI = enc.EnclosureGroupURI
	info.LogicalEnclosureURI = enc.LogicalEnclosureURI
	info.DeviceBayCount = enc.DeviceBayCount.Int()
	info.Status = enc.Status
	info.State = enc.State
	info.ActiveOAPreferredIP = enc.ActiveOAPreferredIP
	return info
}

func bayInfoFrom(hw *ServerHardware, enc *Enclosure) BayInfo {
	info := BayInfo{}
	if hw != nil {
		info.BayNumber = hw.Position.Int()
		info.DeviceURI = first(hw.URI, hw.OriginalURI)
		info.PowerState = hw.PowerState
		info.ProfileURI = hw.ServerProfileURI
	}
	if enc == nil {
		return info
	}
	var match *EnclosureDeviceBay
	for i := range enc.DeviceBays {
		b := &enc.DeviceBays[i]
		if hw != nil && (sameResource(b.DeviceURI, hw.URI) || sameResource(b.DeviceURI, hw.OriginalURI)) {
			match = b
			break
		}
	}
	if match == nil && hw != nil && hw.Position.Int() > 0 {
		for i := range enc.DeviceBays {
			b := &enc.DeviceBays[i]
			if b.BayNumber.Int() == hw.Position.Int() {
				match = b
				break
			}
		}
	}
	if match == nil {
		return info
	}
	info.BayNumber = match.BayNumber.Int()
	info.DevicePresence = match.DevicePresence
	info.DeviceURI = first(match.DeviceURI, info.DeviceURI)
	info.PowerState = first(match.BayPowerState, info.PowerState)
	info.DeviceBayType = match.DeviceBayType
	info.Model = match.Model
	info.ProfileURI = first(match.ProfileURI, info.ProfileURI)
	info.Covered = match.Covered.Bool()
	return info
}

func (cat locationCatalog) rackAndDatacenter(hw *ServerHardware, enc *Enclosure, encInfo EnclosureInfo) (RackInfo, DatacenterInfo) {
	var rack RackInfo
	var mount RackMount
	found := false
	candidates := make([]string, 0, 6)
	if hw != nil {
		candidates = append(candidates, hw.URI, hw.OriginalURI)
	}
	if enc != nil {
		candidates = append(candidates, enc.URI, enc.OriginalURI)
	}
	candidates = append(candidates, encInfo.URI)

	for _, r := range cat.racks {
		for _, m := range r.RackMounts {
			if anySameResource(m.MountURI, candidates...) {
				rack = rackInfoFrom(r)
				mount = m
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		name := encInfo.RackName
		if enc != nil {
			name = first(name, enc.RackName)
		}
		if name != "" {
			for _, r := range cat.racks {
				if strings.EqualFold(r.Name, name) {
					rack = rackInfoFrom(r)
					found = true
					break
				}
			}
		}
	}
	if found {
		rack.TopUSlot = mount.TopUSlot.Int()
		rack.MountUHeight = mount.UHeight.Int()
		rack.MountLocation = mount.Location
		if rack.Name == "" {
			rack.Name = encInfo.RackName
		}
	} else if encInfo.RackName != "" {
		rack.Name = encInfo.RackName
	}

	dc := cat.datacenterForRack(rack)
	return rack, dc
}

func (cat locationCatalog) datacenterForRack(rack RackInfo) DatacenterInfo {
	if rack.URI == "" && rack.Name == "" {
		return DatacenterInfo{}
	}
	for _, dc := range cat.dcs {
		for _, item := range dc.Contents {
			if sameResource(item.ResourceURI, rack.URI) {
				return datacenterInfoFrom(dc)
			}
		}
		for _, item := range dc.RackInventory {
			if sameResource(item.OriginalURI, rack.URI) || (rack.Name != "" && strings.EqualFold(item.Name, rack.Name)) {
				return datacenterInfoFrom(dc)
			}
		}
	}
	return DatacenterInfo{}
}

func rackInfoFrom(r Rack) RackInfo {
	return RackInfo{
		Name:         r.Name,
		URI:          first(r.URI, r.OriginalURI),
		UUID:         first(r.UUID, r.ID),
		Model:        r.Model,
		SerialNumber: r.SerialNumber,
		PartNumber:   r.PartNumber,
		UHeight:      r.UHeight.Int(),
		ThermalLimit: r.ThermalLimit.Int(),
		Status:       r.Status,
		State:        r.State,
	}
}

func datacenterInfoFrom(dc Datacenter) DatacenterInfo {
	return DatacenterInfo{
		Name:            dc.Name,
		URI:             first(dc.URI, dc.OriginalURI),
		UUID:            first(dc.UUID, dc.ID),
		Status:          dc.Status,
		State:           dc.State,
		CoolingCapacity: dc.CoolingCapacity.Int(),
		Width:           dc.Width.Int(),
		Height:          dc.Height.Int(),
		Currency:        dc.Currency,
	}
}

func anySameResource(uri string, candidates ...string) bool {
	for _, c := range candidates {
		if sameResource(uri, c) {
			return true
		}
	}
	return false
}

func sameResource(a, b string) bool {
	ka, kb := resourceKey(a), resourceKey(b)
	if ka == "" || kb == "" {
		return false
	}
	return ka == kb
}

func resourceKey(uri string) string {
	u := strings.ToLower(strings.TrimRight(strings.TrimSpace(uri), "/"))
	if u == "" {
		return ""
	}
	if i := strings.Index(u, "/rest/"); i >= 0 {
		return u[i:]
	}
	return u
}

func mergeDatacenterInfo(dst *DatacenterInfo, src DatacenterInfo) {
	dst.Name = first(dst.Name, src.Name)
	dst.URI = first(dst.URI, src.URI)
	dst.UUID = first(dst.UUID, src.UUID)
	dst.Status = first(dst.Status, src.Status)
	dst.State = first(dst.State, src.State)
	if dst.CoolingCapacity == 0 {
		dst.CoolingCapacity = src.CoolingCapacity
	}
	if dst.Width == 0 {
		dst.Width = src.Width
	}
	if dst.Height == 0 {
		dst.Height = src.Height
	}
	dst.Currency = first(dst.Currency, src.Currency)
}

func mergeRackInfo(dst *RackInfo, src RackInfo) {
	dst.Name = first(dst.Name, src.Name)
	dst.URI = first(dst.URI, src.URI)
	dst.UUID = first(dst.UUID, src.UUID)
	dst.Model = first(dst.Model, src.Model)
	dst.SerialNumber = first(dst.SerialNumber, src.SerialNumber)
	dst.PartNumber = first(dst.PartNumber, src.PartNumber)
	if dst.UHeight == 0 {
		dst.UHeight = src.UHeight
	}
	if dst.TopUSlot == 0 {
		dst.TopUSlot = src.TopUSlot
	}
	if dst.MountUHeight == 0 {
		dst.MountUHeight = src.MountUHeight
	}
	dst.MountLocation = first(dst.MountLocation, src.MountLocation)
	if dst.ThermalLimit == 0 {
		dst.ThermalLimit = src.ThermalLimit
	}
	dst.Status = first(dst.Status, src.Status)
	dst.State = first(dst.State, src.State)
}

func mergeEnclosureInfo(dst *EnclosureInfo, src EnclosureInfo) {
	dst.Name = first(dst.Name, src.Name)
	dst.URI = first(dst.URI, src.URI)
	dst.UUID = first(dst.UUID, src.UUID)
	dst.SerialNumber = first(dst.SerialNumber, src.SerialNumber)
	dst.Model = first(dst.Model, src.Model)
	dst.EnclosureType = first(dst.EnclosureType, src.EnclosureType)
	dst.PartNumber = first(dst.PartNumber, src.PartNumber)
	dst.RackName = first(dst.RackName, src.RackName)
	dst.GroupURI = first(dst.GroupURI, src.GroupURI)
	dst.LogicalEnclosureURI = first(dst.LogicalEnclosureURI, src.LogicalEnclosureURI)
	if dst.DeviceBayCount == 0 {
		dst.DeviceBayCount = src.DeviceBayCount
	}
	dst.Status = first(dst.Status, src.Status)
	dst.State = first(dst.State, src.State)
	dst.ActiveOAPreferredIP = first(dst.ActiveOAPreferredIP, src.ActiveOAPreferredIP)
}

func mergeBayInfo(dst *BayInfo, src BayInfo) {
	if dst.BayNumber == 0 {
		dst.BayNumber = src.BayNumber
	}
	dst.DevicePresence = first(dst.DevicePresence, src.DevicePresence)
	dst.DeviceURI = first(dst.DeviceURI, src.DeviceURI)
	dst.PowerState = first(dst.PowerState, src.PowerState)
	dst.DeviceBayType = first(dst.DeviceBayType, src.DeviceBayType)
	dst.Model = first(dst.Model, src.Model)
	dst.ProfileURI = first(dst.ProfileURI, src.ProfileURI)
	if !dst.Covered {
		dst.Covered = src.Covered
	}
}
