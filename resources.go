package oneview

import "context"

func list[T any](ctx context.Context, c *Client, path string, opts ListOptions) (*Collection[T], error) {
	return GetAll[T](ctx, c, path, opts)
}

func get[T any](ctx context.Context, c *Client, path, id string) (*T, error) {
	var v T
	if err := c.GetJSON(ctx, resourcePath(path, id), &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// ListResourceAlerts is GET /rest/resource-alerts (Global Dashboard swagger).
func (c *Client) ListResourceAlerts(ctx context.Context, opts ListOptions) (*Collection[ResourceAlert], error) {
	return list[ResourceAlert](ctx, c, URIResourceAlerts, opts)
}

// GetResourceAlert is GET /rest/resource-alerts/{id}.
func (c *Client) GetResourceAlert(ctx context.Context, id string) (*ResourceAlert, error) {
	return get[ResourceAlert](ctx, c, URIResourceAlerts, id)
}

// ListAlerts is GET /rest/alerts on an OneView appliance (API 3800–8800).
func (c *Client) ListAlerts(ctx context.Context, opts ListOptions) (*Collection[ResourceAlert], error) {
	if c.IsGlobalDashboard() {
		return c.ListResourceAlerts(ctx, opts)
	}
	return list[ResourceAlert](ctx, c, URIAlerts, opts)
}

// GetAlert is GET /rest/alerts/{id} (appliance) or /rest/resource-alerts/{id} (GD).
func (c *Client) GetAlert(ctx context.Context, id string) (*ResourceAlert, error) {
	if c.IsGlobalDashboard() {
		return c.GetResourceAlert(ctx, id)
	}
	return get[ResourceAlert](ctx, c, URIAlerts, id)
}

// GetAlertSettings is GET /rest/admin-settings/alert-settings.
func (c *Client) GetAlertSettings(ctx context.Context) (*AlertSettings, error) {
	var v AlertSettings
	if err := c.GetJSON(ctx, URIAlertSettings, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateAlertSettings is PUT /rest/admin-settings/alert-settings.
func (c *Client) UpdateAlertSettings(ctx context.Context, s AlertSettings) (*AlertSettings, error) {
	var out AlertSettings
	if _, err := c.PutJSON(ctx, URIAlertSettings, s, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAppliances is GET /rest/appliances (Global Dashboard).
func (c *Client) ListAppliances(ctx context.Context, opts ListOptions) (*Collection[Appliance], error) {
	return list[Appliance](ctx, c, URIAppliances, opts)
}

// GetAppliance is GET /rest/appliances/{id}.
func (c *Client) GetAppliance(ctx context.Context, id string) (*Appliance, error) {
	return get[Appliance](ctx, c, URIAppliances, id)
}

// AddAppliance is POST /rest/appliances. Returns 202 and a task on success.
func (c *Client) AddAppliance(ctx context.Context, a ApplianceAdd, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URIAppliances, a, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// DeleteAppliance is DELETE /rest/appliances/{id}.
func (c *Client) DeleteAppliance(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	resp, err := c.DeleteJSON(ctx, resourcePath(URIAppliances, id), nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// PatchAppliance is PATCH /rest/appliances/{id}.
func (c *Client) PatchAppliance(ctx context.Context, id string, ops []PatchOp, wait bool) (*Response, *Task, error) {
	resp, err := c.PatchJSON(ctx, resourcePath(URIAppliances, id), ops, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// GetApplianceSSO is GET /rest/appliances/{id}/sso.
func (c *Client) GetApplianceSSO(ctx context.Context, id string) (*SSOURL, error) {
	var v SSOURL
	if err := c.GetJSON(ctx, joinPath(URIAppliances, id, "sso"), &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// GetAuditLogSettings is GET /rest/audit-logs/settings.
func (c *Client) GetAuditLogSettings(ctx context.Context) (*AuditLogSettings, error) {
	var v AuditLogSettings
	if err := c.GetJSON(ctx, URIAuditLogSettings, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateAuditLogSettings is PUT /rest/audit-logs/settings.
func (c *Client) UpdateAuditLogSettings(ctx context.Context, s AuditLogSettings) (*AuditLogSettings, error) {
	var out AuditLogSettings
	if _, err := c.PutJSON(ctx, URIAuditLogSettings, s, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TestAuditLogForwarding is POST /rest/audit-logs/test-forwarding.
func (c *Client) TestAuditLogForwarding(ctx context.Context, s AuditLogSettings) (*Response, error) {
	return c.PostJSON(ctx, URIAuditLogTest, s, nil)
}

// GetRemoteCertificate is GET /rest/certificates/https/remote/{address}.
func (c *Client) GetRemoteCertificate(ctx context.Context, address string) (*RemoteCertificate, error) {
	return get[RemoteCertificate](ctx, c, URIRemoteCertificate, address)
}

// AddServerCertificate is POST /rest/certificates/servers.
func (c *Client) AddServerCertificate(ctx context.Context, cert ServerCertificate) (*ServerCertificate, error) {
	var out ServerCertificate
	if _, err := c.PostJSON(ctx, URIServerCertificates, cert, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteServerCertificate is DELETE /rest/certificates/servers/{aliasName}.
func (c *Client) DeleteServerCertificate(ctx context.Context, alias string) error {
	_, err := c.DeleteJSON(ctx, resourcePath(URIServerCertificates, alias), nil)
	return err
}

// ListConvergedSystems is GET /rest/converged-systems.
func (c *Client) ListConvergedSystems(ctx context.Context, opts ListOptions) (*Collection[ConvergedSystem], error) {
	return list[ConvergedSystem](ctx, c, URIConvergedSystems, opts)
}

// GetConvergedSystem is GET /rest/converged-systems/{id}.
func (c *Client) GetConvergedSystem(ctx context.Context, id string) (*ConvergedSystem, error) {
	return get[ConvergedSystem](ctx, c, URIConvergedSystems, id)
}

// ListDatacenters is GET /rest/datacenters.
func (c *Client) ListDatacenters(ctx context.Context, opts ListOptions) (*Collection[Datacenter], error) {
	return list[Datacenter](ctx, c, URIDatacenters, opts)
}

// GetDatacenter is GET /rest/datacenters/{id}.
func (c *Client) GetDatacenter(ctx context.Context, id string) (*Datacenter, error) {
	return get[Datacenter](ctx, c, URIDatacenters, id)
}

// ListDriveEnclosures is GET /rest/drive-enclosures.
func (c *Client) ListDriveEnclosures(ctx context.Context, opts ListOptions) (*Collection[DriveEnclosure], error) {
	return list[DriveEnclosure](ctx, c, URIDriveEnclosures, opts)
}

// GetDriveEnclosure is GET /rest/drive-enclosures/{id}.
func (c *Client) GetDriveEnclosure(ctx context.Context, id string) (*DriveEnclosure, error) {
	return get[DriveEnclosure](ctx, c, URIDriveEnclosures, id)
}

// ListEnclosures is GET /rest/enclosures.
func (c *Client) ListEnclosures(ctx context.Context, opts ListOptions) (*Collection[Enclosure], error) {
	return list[Enclosure](ctx, c, URIEnclosures, opts)
}

// GetEnclosure is GET /rest/enclosures/{id}.
func (c *Client) GetEnclosure(ctx context.Context, id string) (*Enclosure, error) {
	return get[Enclosure](ctx, c, URIEnclosures, id)
}

// GetEnclosureOASSO is GET /rest/enclosures/{id}/oaSsoUrl.
func (c *Client) GetEnclosureOASSO(ctx context.Context, id string) (*SSOURL, error) {
	var v SSOURL
	if err := c.GetJSON(ctx, joinPath(URIEnclosures, IDFromURI(id), "oaSsoUrl"), &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// ListGroups is GET /rest/groups (Global Dashboard).
func (c *Client) ListGroups(ctx context.Context, opts ListOptions) (*Collection[Group], error) {
	col, err := list[Group](ctx, c, URIGroups, opts)
	if err != nil {
		return list[Group](ctx, c, URIGroups+"/", opts)
	}
	return col, nil
}

// GetGroup is GET /rest/groups/{id}.
func (c *Client) GetGroup(ctx context.Context, id string) (*Group, error) {
	return get[Group](ctx, c, URIGroups, id)
}

// CreateGroup is POST /rest/groups.
func (c *Client) CreateGroup(ctx context.Context, g GroupCreate) (*Group, error) {
	var out Group
	if _, err := c.PostJSON(ctx, URIGroups+"/", g, &out); err != nil {
		if _, err2 := c.PostJSON(ctx, URIGroups, g, &out); err2 != nil {
			return nil, err
		}
	}
	return &out, nil
}

// DeleteGroup is DELETE /rest/groups/{id}.
func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	_, err := c.DeleteJSON(ctx, resourcePath(URIGroups, id), nil)
	return err
}

// PatchGroup is PATCH /rest/groups/{id}.
func (c *Client) PatchGroup(ctx context.Context, id string, ops []PatchOp) (*Group, error) {
	var out Group
	if _, err := c.PatchJSON(ctx, resourcePath(URIGroups, id), ops, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGroupMembers is GET /rest/groups/{id}/members.
func (c *Client) GetGroupMembers(ctx context.Context, id string, opts ListOptions) (*Collection[Resource], error) {
	return list[Resource](ctx, c, joinPath(URIGroups, IDFromURI(id), "members"), opts)
}

// ListInterconnects is GET /rest/interconnects.
func (c *Client) ListInterconnects(ctx context.Context, opts ListOptions) (*Collection[Interconnect], error) {
	return list[Interconnect](ctx, c, URIInterconnects, opts)
}

// GetInterconnect is GET /rest/interconnects/{id}.
func (c *Client) GetInterconnect(ctx context.Context, id string) (*Interconnect, error) {
	return get[Interconnect](ctx, c, URIInterconnects, id)
}

// ListLogicalInterconnects is GET /rest/logical-interconnects.
func (c *Client) ListLogicalInterconnects(ctx context.Context, opts ListOptions) (*Collection[LogicalInterconnect], error) {
	return list[LogicalInterconnect](ctx, c, URILogicalInterconnects, opts)
}

// GetLogicalInterconnect is GET /rest/logical-interconnects/{id}.
func (c *Client) GetLogicalInterconnect(ctx context.Context, id string) (*LogicalInterconnect, error) {
	return get[LogicalInterconnect](ctx, c, URILogicalInterconnects, id)
}

// ListManagedSANs is GET /rest/managed-sans.
func (c *Client) ListManagedSANs(ctx context.Context, opts ListOptions) (*Collection[ManagedSAN], error) {
	return list[ManagedSAN](ctx, c, URIManagedSANs, opts)
}

// GetManagedSAN is GET /rest/managed-sans/{id}.
func (c *Client) GetManagedSAN(ctx context.Context, id string) (*ManagedSAN, error) {
	return get[ManagedSAN](ctx, c, URIManagedSANs, id)
}

// GetNetworkInterfaces is GET /rest/appliance/network-interfaces.
func (c *Client) GetNetworkInterfaces(ctx context.Context) (*NetworkInterfaces, error) {
	var v NetworkInterfaces
	if err := c.GetJSON(ctx, URINetworkInterfaces, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateNetworkInterfaces is POST /rest/appliance/network-interfaces.
func (c *Client) UpdateNetworkInterfaces(ctx context.Context, n NetworkInterfaces, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URINetworkInterfaces, n, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// GetTimeLocale is GET /rest/appliance/configuration/time-locale.
func (c *Client) GetTimeLocale(ctx context.Context) (*TimeLocale, error) {
	var v TimeLocale
	if err := c.GetJSON(ctx, URITimeLocale, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateTimeLocale is POST /rest/appliance/configuration/time-locale.
func (c *Client) UpdateTimeLocale(ctx context.Context, t TimeLocale, wait bool) (*Response, *Task, error) {
	resp, err := c.PostJSON(ctx, URITimeLocale, t, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	task, err := c.WaitResponse(ctx, resp)
	return resp, task, err
}

// ListSANManagers is GET /rest/san-managers.
func (c *Client) ListSANManagers(ctx context.Context, opts ListOptions) (*Collection[SANManager], error) {
	return list[SANManager](ctx, c, URISANManagers, opts)
}

// GetSANManager is GET /rest/san-managers/{id}.
func (c *Client) GetSANManager(ctx context.Context, id string) (*SANManager, error) {
	return get[SANManager](ctx, c, URISANManagers, id)
}

// ListServerFirmware is GET /rest/server-firmware (Global Dashboard).
func (c *Client) ListServerFirmware(ctx context.Context, opts ListOptions) (*Collection[ServerFirmware], error) {
	return list[ServerFirmware](ctx, c, URIServerFirmware, opts)
}

// GetServerFirmware is GET /rest/server-firmware/{id}.
func (c *Client) GetServerFirmware(ctx context.Context, id string) (*ServerFirmware, error) {
	return get[ServerFirmware](ctx, c, URIServerFirmware, id)
}

// ListServerHardware is GET /rest/server-hardware.
func (c *Client) ListServerHardware(ctx context.Context, opts ListOptions) (*Collection[ServerHardware], error) {
	return list[ServerHardware](ctx, c, URIServerHardware, opts)
}

// GetServerHardware is GET /rest/server-hardware/{id}.
func (c *Client) GetServerHardware(ctx context.Context, id string) (*ServerHardware, error) {
	return get[ServerHardware](ctx, c, URIServerHardware, id)
}

// GetServerHardwareILOSSO is GET /rest/server-hardware/{id}/iloSsoUrl.
func (c *Client) GetServerHardwareILOSSO(ctx context.Context, id string) (*SSOURL, error) {
	var v SSOURL
	if err := c.GetJSON(ctx, joinPath(URIServerHardware, IDFromURI(id), "iloSsoUrl"), &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// SetServerPower is PUT /rest/server-hardware/{id}/powerState (appliance 3800–8800).
func (c *Client) SetServerPower(ctx context.Context, id string, req PowerStateRequest, wait bool) (*Response, *Task, error) {
	resp, err := c.PutJSON(ctx, joinPath(URIServerHardware, IDFromURI(id), "powerState"), req, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// RefreshServerHardware is PUT /rest/server-hardware/{id}/refreshState.
func (c *Client) RefreshServerHardware(ctx context.Context, id string, wait bool) (*Response, *Task, error) {
	body := map[string]string{"refreshState": "RefreshPending"}
	resp, err := c.PutJSON(ctx, joinPath(URIServerHardware, IDFromURI(id), "refreshState"), body, nil)
	if err != nil {
		return resp, nil, err
	}
	if !wait {
		return resp, nil, nil
	}
	t, err := c.WaitResponse(ctx, resp)
	return resp, t, err
}

// ListServerProfiles is GET /rest/server-profiles.
func (c *Client) ListServerProfiles(ctx context.Context, opts ListOptions) (*Collection[ServerProfile], error) {
	return list[ServerProfile](ctx, c, URIServerProfiles, opts)
}

// GetServerProfile is GET /rest/server-profiles/{id}.
func (c *Client) GetServerProfile(ctx context.Context, id string) (*ServerProfile, error) {
	return get[ServerProfile](ctx, c, URIServerProfiles, id)
}

// ListServerProfileTemplates is GET /rest/server-profile-templates.
func (c *Client) ListServerProfileTemplates(ctx context.Context, opts ListOptions) (*Collection[ServerProfileTemplate], error) {
	return list[ServerProfileTemplate](ctx, c, URIServerProfileTemplates, opts)
}

// GetServerProfileTemplate is GET /rest/server-profile-templates/{id}.
func (c *Client) GetServerProfileTemplate(ctx context.Context, id string) (*ServerProfileTemplate, error) {
	return get[ServerProfileTemplate](ctx, c, URIServerProfileTemplates, id)
}

func (c *Client) volumesURI() string {
	if c.IsGlobalDashboard() {
		return URIStorageVolumes
	}
	return URIVolumes
}

// ListStoragePools is GET /rest/storage-pools.
func (c *Client) ListStoragePools(ctx context.Context, opts ListOptions) (*Collection[StoragePool], error) {
	return list[StoragePool](ctx, c, URIStoragePools, opts)
}

// GetStoragePool is GET /rest/storage-pools/{id}.
func (c *Client) GetStoragePool(ctx context.Context, id string) (*StoragePool, error) {
	return get[StoragePool](ctx, c, URIStoragePools, id)
}

// ListStorageSystems is GET /rest/storage-systems.
func (c *Client) ListStorageSystems(ctx context.Context, opts ListOptions) (*Collection[StorageSystem], error) {
	return list[StorageSystem](ctx, c, URIStorageSystems, opts)
}

// GetStorageSystem is GET /rest/storage-systems/{id}.
func (c *Client) GetStorageSystem(ctx context.Context, id string) (*StorageSystem, error) {
	return get[StorageSystem](ctx, c, URIStorageSystems, id)
}

// ListStorageVolumes is GET /rest/storage-volumes (GD) or /rest/volumes (appliance).
func (c *Client) ListStorageVolumes(ctx context.Context, opts ListOptions) (*Collection[StorageVolume], error) {
	return list[StorageVolume](ctx, c, c.volumesURI(), opts)
}

// GetStorageVolume is GET /rest/storage-volumes/{id} or /rest/volumes/{id}.
func (c *Client) GetStorageVolume(ctx context.Context, id string) (*StorageVolume, error) {
	return get[StorageVolume](ctx, c, c.volumesURI(), id)
}
