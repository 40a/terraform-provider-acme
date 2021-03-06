package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Server (virtual machine) on Vultr account
type Server struct {
	ID               string      `json:"SUBID"`
	Name             string      `json:"label"`
	OS               string      `json:"os"`
	RAM              string      `json:"ram"`
	Disk             string      `json:"disk"`
	MainIP           string      `json:"main_ip"`
	VCpus            int         `json:"vcpu_count,string"`
	Location         string      `json:"location"`
	RegionID         int         `json:"DCID,string"`
	DefaultPassword  string      `json:"default_password"`
	Created          string      `json:"date_created"`
	PendingCharges   float64     `json:"pending_charges"`
	Status           string      `json:"status"`
	Cost             string      `json:"cost_per_month"`
	CurrentBandwidth float64     `json:"current_bandwidth_gb"`
	AllowedBandwidth float64     `json:"allowed_bandwidth_gb,string"`
	NetmaskV4        string      `json:"netmask_v4"`
	GatewayV4        string      `json:"gateway_v4"`
	PowerStatus      string      `json:"power_status"`
	ServerState      string      `json:"server_state"`
	PlanID           int         `json:"VPSPLANID,string"`
	V6Networks       []V6Network `json:"v6_networks"`
	InternalIP       string      `json:"internal_ip"`
	KVMUrl           string      `json:"kvm_url"`
	AutoBackups      string      `json:"auto_backups"`
	Tag              string      `json:"tag"`
}

// ServerOptions are optional parameters to be used during server creation
type ServerOptions struct {
	IPXEChainURL         string
	ISO                  int
	Script               int
	UserData             string
	Snapshot             string
	SSHKey               string
	IPV6                 bool
	PrivateNetworking    bool
	AutoBackups          bool
	DontNotifyOnActivate bool
}

// V6Network represents a IPv6 network of a Vultr server
type V6Network struct {
	Network     string `json:"v6_network"`
	MainIP      string `json:"v6_main_ip"`
	NetworkSize string `json:"v6_network_size"`
}

// ISOStatus represents an ISO image attached to a Vultr server
type ISOStatus struct {
	State string `json:"state"`
	ISOID string `json:"ISOID"`
}

// UnmarshalJSON implements json.Unmarshaller on Server.
// This is needed because the Vultr API is inconsistent in it's JSON responses for servers.
// Some fields can change type, from JSON number to JSON string and vice-versa.
func (s *Server) UnmarshalJSON(data []byte) (err error) {
	if s == nil {
		*s = Server{}
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	value := fmt.Sprintf("%v", fields["vcpu_count"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	vcpu, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	s.VCpus = int(vcpu)

	value = fmt.Sprintf("%v", fields["DCID"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	region, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	s.RegionID = int(region)

	value = fmt.Sprintf("%v", fields["VPSPLANID"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	plan, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	s.PlanID = int(plan)

	value = fmt.Sprintf("%v", fields["pending_charges"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	pc, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	s.PendingCharges = pc

	value = fmt.Sprintf("%v", fields["current_bandwidth_gb"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	cb, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	s.CurrentBandwidth = cb

	value = fmt.Sprintf("%v", fields["allowed_bandwidth_gb"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	ab, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	s.AllowedBandwidth = ab

	s.ID = fmt.Sprintf("%v", fields["SUBID"])
	s.Name = fmt.Sprintf("%v", fields["label"])
	s.OS = fmt.Sprintf("%v", fields["os"])
	s.RAM = fmt.Sprintf("%v", fields["ram"])
	s.Disk = fmt.Sprintf("%v", fields["disk"])
	s.MainIP = fmt.Sprintf("%v", fields["main_ip"])
	s.Location = fmt.Sprintf("%v", fields["location"])
	s.DefaultPassword = fmt.Sprintf("%v", fields["default_password"])
	s.Created = fmt.Sprintf("%v", fields["date_created"])
	s.Status = fmt.Sprintf("%v", fields["status"])
	s.Cost = fmt.Sprintf("%v", fields["cost_per_month"])
	s.NetmaskV4 = fmt.Sprintf("%v", fields["netmask_v4"])
	s.GatewayV4 = fmt.Sprintf("%v", fields["gateway_v4"])
	s.PowerStatus = fmt.Sprintf("%v", fields["power_status"])
	s.ServerState = fmt.Sprintf("%v", fields["server_state"])

	v6networks := make([]V6Network, 0)
	if networks, ok := fields["v6_networks"].([]interface{}); ok {
		for _, network := range networks {
			if network, ok := network.(map[string]interface{}); ok {
				v6network := V6Network{
					Network:     fmt.Sprintf("%v", network["v6_network"]),
					MainIP:      fmt.Sprintf("%v", network["v6_main_ip"]),
					NetworkSize: fmt.Sprintf("%v", network["v6_network_size"]),
				}
				v6networks = append(v6networks, v6network)
			}
		}
		s.V6Networks = v6networks
	}

	s.InternalIP = fmt.Sprintf("%v", fields["internal_ip"])
	s.KVMUrl = fmt.Sprintf("%v", fields["kvm_url"])
	s.AutoBackups = fmt.Sprintf("%v", fields["auto_backups"])
	s.Tag = fmt.Sprintf("%v", fields["tag"])

	return
}

func (c *Client) GetServers() (servers []Server, err error) {
	var serverMap map[string]Server
	if err := c.get(`server/list`, &serverMap); err != nil {
		return nil, err
	}

	for _, server := range serverMap {
		servers = append(servers, server)
	}
	return servers, nil
}

func (c *Client) GetServersByTag(tag string) (servers []Server, err error) {
	var serverMap map[string]Server
	if err := c.get(`server/list?tag=`+tag, &serverMap); err != nil {
		return nil, err
	}

	for _, server := range serverMap {
		servers = append(servers, server)
	}
	return servers, nil
}

func (c *Client) GetServer(id string) (server Server, err error) {
	if err := c.get(`server/list?SUBID=`+id, &server); err != nil {
		return Server{}, err
	}
	return server, nil
}

func (c *Client) CreateServer(name string, regionID, planID, osID int, options *ServerOptions) (Server, error) {
	values := url.Values{
		"label":     {name},
		"DCID":      {fmt.Sprintf("%v", regionID)},
		"VPSPLANID": {fmt.Sprintf("%v", planID)},
		"OSID":      {fmt.Sprintf("%v", osID)},
	}

	if options != nil {
		if options.IPXEChainURL != "" {
			values.Add("ipxe_chain_url", options.IPXEChainURL)
		}

		if options.ISO != 0 {
			values.Add("ISOID", fmt.Sprintf("%v", options.ISO))
		}

		if options.Script != 0 {
			values.Add("SCRIPTID", fmt.Sprintf("%v", options.Script))
		}

		if options.UserData != "" {
			values.Add("userdata", base64.StdEncoding.EncodeToString([]byte(options.UserData)))
		}

		if options.Snapshot != "" {
			values.Add("SNAPSHOTID", options.Snapshot)
		}

		if options.SSHKey != "" {
			values.Add("SSHKEYID", options.SSHKey)
		}

		values.Add("enable_ipv6", "no")
		if options.IPV6 {
			values.Set("enable_ipv6", "yes")
		}

		values.Add("enable_private_network", "no")
		if options.PrivateNetworking {
			values.Set("enable_private_network", "yes")
		}

		values.Add("auto_backups", "no")
		if options.AutoBackups {
			values.Set("auto_backups", "yes")
		}

		values.Add("notify_activate", "yes")
		if options.DontNotifyOnActivate {
			values.Set("notify_activate", "no")
		}
	}

	var server Server
	if err := c.post(`server/create`, values, &server); err != nil {
		return Server{}, err
	}
	server.Name = name
	server.RegionID = regionID
	server.PlanID = planID

	return server, nil
}

func (c *Client) RenameServer(id, name string) error {
	values := url.Values{
		"SUBID": {id},
		"label": {name},
	}

	if err := c.post(`server/label_set`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) StartServer(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`server/start`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) HaltServer(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`server/halt`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) RebootServer(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`server/reboot`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReinstallServer(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`server/reinstall`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) ChangeOSofServer(id string, osID int) error {
	values := url.Values{
		"SUBID": {id},
		"OSID":  {fmt.Sprintf("%v", osID)},
	}

	if err := c.post(`server/os_change`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) AttachISOtoServer(id string, isoID int) error {
	values := url.Values{
		"SUBID": {id},
		"ISOID": {fmt.Sprintf("%v", isoID)},
	}

	if err := c.post(`server/iso_attach`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DetachISOfromServer(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`server/iso_detach`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetISOStatusofServer(id string) (isoStatus ISOStatus, err error) {
	if err := c.get(`server/iso_status?SUBID=`+id, &isoStatus); err != nil {
		return ISOStatus{}, err
	}
	return isoStatus, nil
}

func (c *Client) ListOSforServer(id string) (os []OS, err error) {
	var osMap map[string]OS
	if err := c.get(`server/os_change_list?SUBID=`+id, &osMap); err != nil {
		return nil, err
	}

	for _, o := range osMap {
		os = append(os, o)
	}
	return os, nil
}

func (c *Client) DeleteServer(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`server/destroy`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) BandwidthOfServer(id string) (bandwidth []map[string]string, err error) {
	var bandwidthMap map[string][][]string
	if err := c.get(`server/bandwidth?SUBID=`+id, &bandwidthMap); err != nil {
		return nil, err
	}

	// parse incoming bytes
	for _, b := range bandwidthMap["incoming_bytes"] {
		bMap := make(map[string]string)
		bMap["date"] = b[0]
		bMap["incoming"] = b[1]
		bandwidth = append(bandwidth, bMap)
	}

	// parse outgoing bytes (we'll assume that incoming and outgoing dates are always a match)
	for _, b := range bandwidthMap["outgoing_bytes"] {
		for i := range bandwidth {
			if bandwidth[i]["date"] == b[0] {
				bandwidth[i]["outgoing"] = b[1]
				break
			}
		}
	}

	return bandwidth, nil
}
