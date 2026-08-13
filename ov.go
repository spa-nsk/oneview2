package oneview

import "context"

// Appliance REST API 3800–8800 resource methods that are not in swagger 300.json
// (Global Dashboard is read-mostly and has a smaller surface).

// CreateServerProfile is POST /rest/server-profiles.
func (c *Client) CreateServerProfile(ctx context.Context, p ServerProfile, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URIServerProfiles, p, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// UpdateServerProfile is PUT /rest/server-profiles/{id}.
func (c *Client) UpdateServerProfile(ctx context.Context, p ServerProfile, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, resourcePath(URIServerProfiles, p.URI), p, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteServerProfile is DELETE /rest/server-profiles/{id}.
func (c *Client) DeleteServerProfile(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URIServerProfiles, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// CreateServerProfileTemplate is POST /rest/server-profile-templates.
func (c *Client) CreateServerProfileTemplate(ctx context.Context, t ServerProfileTemplate, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URIServerProfileTemplates, t, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	task, err := c.WaitResponse(ctx, resp)
	return resp, task, err
}

// UpdateServerProfileTemplate is PUT /rest/server-profile-templates/{id}.
func (c *Client) UpdateServerProfileTemplate(ctx context.Context, t ServerProfileTemplate, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, resourcePath(URIServerProfileTemplates, t.URI), t, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	task, err := c.WaitResponse(ctx, resp)
	return resp, task, err
}

// DeleteServerProfileTemplate is DELETE /rest/server-profile-templates/{id}.
func (c *Client) DeleteServerProfileTemplate(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URIServerProfileTemplates, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	task, err := c.WaitResponse(ctx, resp)
	return resp, task, err
}

// ListEthernetNetworks is GET /rest/ethernet-networks.
func (c *Client) ListEthernetNetworks(ctx context.Context, opts ListOptions) (*Collection[EthernetNetwork], error) {
	return list[EthernetNetwork](ctx, c, URIEthernetNetworks, opts)
}

// GetEthernetNetwork is GET /rest/ethernet-networks/{id}.
func (c *Client) GetEthernetNetwork(ctx context.Context, id string) (*EthernetNetwork, error) {
	return get[EthernetNetwork](ctx, c, URIEthernetNetworks, id)
}

// CreateEthernetNetwork is POST /rest/ethernet-networks.
func (c *Client) CreateEthernetNetwork(ctx context.Context, n EthernetNetwork, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URIEthernetNetworks, n, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// UpdateEthernetNetwork is PUT /rest/ethernet-networks/{id}.
func (c *Client) UpdateEthernetNetwork(ctx context.Context, n EthernetNetwork, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, resourcePath(URIEthernetNetworks, n.URI), n, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteEthernetNetwork is DELETE /rest/ethernet-networks/{id}.
func (c *Client) DeleteEthernetNetwork(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URIEthernetNetworks, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// ListFCNetworks is GET /rest/fc-networks.
func (c *Client) ListFCNetworks(ctx context.Context, opts ListOptions) (*Collection[FCNetwork], error) {
	return list[FCNetwork](ctx, c, URIFCNetworks, opts)
}

// GetFCNetwork is GET /rest/fc-networks/{id}.
func (c *Client) GetFCNetwork(ctx context.Context, id string) (*FCNetwork, error) {
	return get[FCNetwork](ctx, c, URIFCNetworks, id)
}

// CreateFCNetwork is POST /rest/fc-networks.
func (c *Client) CreateFCNetwork(ctx context.Context, n FCNetwork, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URIFCNetworks, n, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteFCNetwork is DELETE /rest/fc-networks/{id}.
func (c *Client) DeleteFCNetwork(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URIFCNetworks, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// ListFCoENetworks is GET /rest/fcoe-networks.
func (c *Client) ListFCoENetworks(ctx context.Context, opts ListOptions) (*Collection[EthernetNetwork], error) {
	return list[EthernetNetwork](ctx, c, URIFCoENetworks, opts)
}

// GetFCoENetwork is GET /rest/fcoe-networks/{id}.
func (c *Client) GetFCoENetwork(ctx context.Context, id string) (*EthernetNetwork, error) {
	return get[EthernetNetwork](ctx, c, URIFCoENetworks, id)
}

// ListNetworkSets is GET /rest/network-sets.
func (c *Client) ListNetworkSets(ctx context.Context, opts ListOptions) (*Collection[NetworkSet], error) {
	return list[NetworkSet](ctx, c, URINetworkSets, opts)
}

// GetNetworkSet is GET /rest/network-sets/{id}.
func (c *Client) GetNetworkSet(ctx context.Context, id string) (*NetworkSet, error) {
	return get[NetworkSet](ctx, c, URINetworkSets, id)
}

// CreateNetworkSet is POST /rest/network-sets.
func (c *Client) CreateNetworkSet(ctx context.Context, n NetworkSet, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URINetworkSets, n, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteNetworkSet is DELETE /rest/network-sets/{id}.
func (c *Client) DeleteNetworkSet(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URINetworkSets, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// ListEnclosureGroups is GET /rest/enclosure-groups.
func (c *Client) ListEnclosureGroups(ctx context.Context, opts ListOptions) (*Collection[EnclosureGroup], error) {
	return list[EnclosureGroup](ctx, c, URIEnclosureGroups, opts)
}

// GetEnclosureGroup is GET /rest/enclosure-groups/{id}.
func (c *Client) GetEnclosureGroup(ctx context.Context, id string) (*EnclosureGroup, error) {
	return get[EnclosureGroup](ctx, c, URIEnclosureGroups, id)
}

// CreateEnclosureGroup is POST /rest/enclosure-groups.
func (c *Client) CreateEnclosureGroup(ctx context.Context, g EnclosureGroup) (*EnclosureGroup, error) {
	var out EnclosureGroup
	if _, err := c.PostJSON(ctx, URIEnclosureGroups, g, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateEnclosureGroup is PUT /rest/enclosure-groups/{id}.
func (c *Client) UpdateEnclosureGroup(ctx context.Context, g EnclosureGroup) (*EnclosureGroup, error) {
	var out EnclosureGroup
	if _, err := c.PutJSON(ctx, resourcePath(URIEnclosureGroups, g.URI), g, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteEnclosureGroup is DELETE /rest/enclosure-groups/{id}.
func (c *Client) DeleteEnclosureGroup(ctx context.Context, id string) error {
	_, err := c.DeleteJSON(ctx, resourcePath(URIEnclosureGroups, id), nil)
	return err
}

// ListLogicalEnclosures is GET /rest/logical-enclosures.
func (c *Client) ListLogicalEnclosures(ctx context.Context, opts ListOptions) (*Collection[LogicalEnclosure], error) {
	return list[LogicalEnclosure](ctx, c, URILogicalEnclosures, opts)
}

// GetLogicalEnclosure is GET /rest/logical-enclosures/{id}.
func (c *Client) GetLogicalEnclosure(ctx context.Context, id string) (*LogicalEnclosure, error) {
	return get[LogicalEnclosure](ctx, c, URILogicalEnclosures, id)
}

// UpdateFromGroupLogicalEnclosure is PUT /rest/logical-enclosures/{id}/updateFromGroup.
func (c *Client) UpdateFromGroupLogicalEnclosure(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, joinPath(URILogicalEnclosures, IDFromURI(id), "updateFromGroup"), map[string]any{}, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// ListLogicalInterconnectGroups is GET /rest/logical-interconnect-groups.
func (c *Client) ListLogicalInterconnectGroups(ctx context.Context, opts ListOptions) (*Collection[LogicalInterconnectGroup], error) {
	return list[LogicalInterconnectGroup](ctx, c, URILogicalInterconnectGroups, opts)
}

// GetLogicalInterconnectGroup is GET /rest/logical-interconnect-groups/{id}.
func (c *Client) GetLogicalInterconnectGroup(ctx context.Context, id string) (*LogicalInterconnectGroup, error) {
	return get[LogicalInterconnectGroup](ctx, c, URILogicalInterconnectGroups, id)
}

// CreateLogicalInterconnectGroup is POST /rest/logical-interconnect-groups.
func (c *Client) CreateLogicalInterconnectGroup(ctx context.Context, g LogicalInterconnectGroup) (*LogicalInterconnectGroup, error) {
	var out LogicalInterconnectGroup
	if _, err := c.PostJSON(ctx, URILogicalInterconnectGroups, g, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteLogicalInterconnectGroup is DELETE /rest/logical-interconnect-groups/{id}.
func (c *Client) DeleteLogicalInterconnectGroup(ctx context.Context, id string) error {
	_, err := c.DeleteJSON(ctx, resourcePath(URILogicalInterconnectGroups, id), nil)
	return err
}

// ListUplinkSets is GET /rest/uplink-sets.
func (c *Client) ListUplinkSets(ctx context.Context, opts ListOptions) (*Collection[UplinkSet], error) {
	return list[UplinkSet](ctx, c, URIUplinkSets, opts)
}

// GetUplinkSet is GET /rest/uplink-sets/{id}.
func (c *Client) GetUplinkSet(ctx context.Context, id string) (*UplinkSet, error) {
	return get[UplinkSet](ctx, c, URIUplinkSets, id)
}

// ListServerHardwareTypes is GET /rest/server-hardware-types.
func (c *Client) ListServerHardwareTypes(ctx context.Context, opts ListOptions) (*Collection[ServerHardwareType], error) {
	return list[ServerHardwareType](ctx, c, URIServerHardwareTypes, opts)
}

// GetServerHardwareType is GET /rest/server-hardware-types/{id}.
func (c *Client) GetServerHardwareType(ctx context.Context, id string) (*ServerHardwareType, error) {
	return get[ServerHardwareType](ctx, c, URIServerHardwareTypes, id)
}

// ListFirmwareDrivers is GET /rest/firmware-drivers.
func (c *Client) ListFirmwareDrivers(ctx context.Context, opts ListOptions) (*Collection[FirmwareDriver], error) {
	return list[FirmwareDriver](ctx, c, URIFirmwareDrivers, opts)
}

// GetFirmwareDriver is GET /rest/firmware-drivers/{id}.
func (c *Client) GetFirmwareDriver(ctx context.Context, id string) (*FirmwareDriver, error) {
	return get[FirmwareDriver](ctx, c, URIFirmwareDrivers, id)
}

// CreateStorageVolume is POST /rest/volumes (appliance) or /rest/storage-volumes (GD).
func (c *Client) CreateStorageVolume(ctx context.Context, v StorageVolume, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, c.volumesURI(), v, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteStorageVolume is DELETE /rest/volumes/{id}.
func (c *Client) DeleteStorageVolume(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(c.volumesURI(), id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// AddStorageSystem is POST /rest/storage-systems.
func (c *Client) AddStorageSystem(ctx context.Context, s StorageSystem, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URIStorageSystems, s, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteStorageSystem is DELETE /rest/storage-systems/{id}.
func (c *Client) DeleteStorageSystem(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URIStorageSystems, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// ListInterconnectTypes is GET /rest/interconnect-types.
func (c *Client) ListInterconnectTypes(ctx context.Context, opts ListOptions) (*Collection[Resource], error) {
	return list[Resource](ctx, c, URIInterconnectTypes, opts)
}

// ListRacks is GET /rest/racks.
func (c *Client) ListRacks(ctx context.Context, opts ListOptions) (*Collection[Resource], error) {
	return list[Resource](ctx, c, URIRacks, opts)
}

// GetRack is GET /rest/racks/{id}.
func (c *Client) GetRack(ctx context.Context, id string) (*Resource, error) {
	return get[Resource](ctx, c, URIRacks, id)
}

// ListEvents is GET /rest/events.
func (c *Client) ListEvents(ctx context.Context, opts ListOptions) (*Collection[Resource], error) {
	return list[Resource](ctx, c, URIEvents, opts)
}

// ListIndexResources is GET /rest/index/resources — search across resource types.
func (c *Client) ListIndexResources(ctx context.Context, opts ListOptions) (*Collection[Resource], error) {
	return list[Resource](ctx, c, URIIndexResources, opts)
}

// ComplianceCheckLogicalInterconnect is PUT /rest/logical-interconnects/{id}/compliance.
func (c *Client) ComplianceCheckLogicalInterconnect(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, joinPath(URILogicalInterconnects, IDFromURI(id), "compliance"), map[string]any{}, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// UpdateFromGroupLogicalInterconnect is PUT /rest/logical-interconnects/{id}/updateFromGroup.
func (c *Client) UpdateFromGroupLogicalInterconnect(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, joinPath(URILogicalInterconnects, IDFromURI(id), "updateFromGroup"), map[string]any{}, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}
