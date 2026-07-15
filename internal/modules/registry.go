package modules

import "sort"

type State string

const (
	StateCandidate State = "candidate"
	StatePlanned   State = "planned"
)

type Module struct {
	ID          string
	Name        string
	Category    string
	State       State
	Description string
}

type Registry struct{ modules []Module }

func DefaultRegistry() Registry {
	items := []Module{
		{ID: "firewall.fortigate", Name: "FortiGate", Category: "Firewall", State: StateCandidate, Description: "Native configuration hierarchy parser and semantic-analysis boundary."},
		{ID: "firewall.opnsense", Name: "OPNsense", Category: "Firewall", State: StateCandidate, Description: "XML backup import boundary."},
		{ID: "firewall.pfsense", Name: "pfSense", Category: "Firewall", State: StateCandidate, Description: "XML backup import boundary."},
		{ID: "network.cisco.ios", Name: "Cisco IOS", Category: "Network", State: StateCandidate, Description: "2960-family evidence and topology analysis."},
		{ID: "network.cisco.iosxe", Name: "Cisco IOS XE", Category: "Network", State: StateCandidate, Description: "Catalyst 9200, 9300, and 9500 evidence analysis."},
		{ID: "wireless.cisco.c9800", Name: "Catalyst 9800", Category: "Wireless", State: StatePlanned, Description: "Controller, AP, policy, tag, and client evidence."},
		{ID: "integration.zabbix", Name: "Zabbix", Category: "Integration", State: StateCandidate, Description: "Replaceable sender-protocol adapter for canonical Atlas telemetry."},
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return Registry{modules: items}
}

func (r Registry) List() []Module {
	return append([]Module(nil), r.modules...)
}
