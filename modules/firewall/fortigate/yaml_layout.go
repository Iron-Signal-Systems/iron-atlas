package fortigate

import "sort"

// FortiGateYAMLLayout describes only vendor schema structure. It never
// contains source scalar values, unrecognized mapping keys, or VDOM names, so
// callers can include it in upload-safe diagnostics.
type FortiGateYAMLLayout struct {
	RootKind                   string
	RootEntries                int
	CanonicalGlobal            bool
	CanonicalVDOM              bool
	RecognizedDirectSections   []string
	RecognizedNestedSections   []string
	NestedMappingCount         int
	DetectedVDOMContainerCount int
	UnrecognizedRootEntryCount int
}

var fortiGateGlobalSectionNames = map[string]struct{}{
	"system_global": {},
}

// These are the FortiOS CMDB sections the current normalizer consumes. They
// are public schema labels, not configuration values. Keeping the detector
// aligned with the normalizer prevents a successful admission followed by an
// empty snapshot.
var fortiGateVDOMSectionNames = map[string]struct{}{
	"firewall_address":               {},
	"firewall_addrgrp":               {},
	"firewall_central_snat_map":      {},
	"firewall_ippool":                {},
	"firewall_policy":                {},
	"firewall_schedule_onetime":      {},
	"firewall_schedule_recurring":    {},
	"firewall_service_custom":        {},
	"firewall_service_group":         {},
	"firewall_shaper_per_ip_shaper":  {},
	"firewall_shaper_traffic_shaper": {},
	"firewall_shaping_policy":        {},
	"firewall_vip":                   {},
	"firewall_vipgrp":                {},
	"router_bgp":                     {},
	"router_ospf":                    {},
	"router_policy":                  {},
	"router_prefix_list":             {},
	"router_prefix_list6":            {},
	"router_rip":                     {},
	"router_route_map":               {},
	"router_static":                  {},
	"system_dhcp_server":             {},
	"system_interface":               {},
	"system_sdwan":                   {},
	"system_virtual_wan_link":        {},
	"system_zone":                    {},
	"vpn_ipsec_phase1_interface":     {},
}

// DiagnoseFortiGateYAMLLayout returns a content-free description suitable for
// support logs. Only recognized public FortiOS section names are retained.
func DiagnoseFortiGateYAMLLayout(doc *YAMLDocument) FortiGateYAMLLayout {
	var diagnostic FortiGateYAMLLayout
	if doc == nil || doc.Root == nil {
		diagnostic.RootKind = "missing"
		return diagnostic
	}

	root := doc.Root
	diagnostic.RootKind = yamlKindName(root.Kind)
	if root.Kind != YAMLMapping {
		return diagnostic
	}

	diagnostic.RootEntries = len(root.Order)
	diagnostic.CanonicalGlobal = root.Child("global") != nil
	diagnostic.CanonicalVDOM = root.Child("vdom") != nil
	direct := make(map[string]struct{})
	nested := make(map[string]struct{})
	for _, key := range root.Order {
		child := root.Map[key]
		normalized := normalizeKey(key)
		if isRecognizedFortiGateSection(normalized) {
			direct[normalized] = struct{}{}
			continue
		}
		if normalized == "global" || normalized == "vdom" {
			continue
		}
		if child != nil && child.Kind == YAMLMapping {
			diagnostic.NestedMappingCount++
			sections := recognizedFortiGateSections(child)
			if len(sections) > 0 {
				diagnostic.DetectedVDOMContainerCount++
				for _, section := range sections {
					nested[section] = struct{}{}
				}
				continue
			}
		}
		diagnostic.UnrecognizedRootEntryCount++
	}

	diagnostic.RecognizedDirectSections = sortedKeys(direct)
	diagnostic.RecognizedNestedSections = sortedKeys(nested)
	return diagnostic
}

func fortiGateGlobalScope(root *YAMLNode) *YAMLNode {
	if root == nil || root.Kind != YAMLMapping {
		return nil
	}
	if global := root.Child("global"); global != nil && global.Kind == YAMLMapping {
		return global
	}
	if childAlias(root, "system_global") != nil {
		return root
	}
	return nil
}

func fortiGateVDOMEntries(root *YAMLNode) []namedNode {
	if root == nil || root.Kind != YAMLMapping {
		return nil
	}
	if vdom := root.Child("vdom"); vdom != nil {
		return tableEntries(vdom, "$.vdom")
	}

	var entries []namedNode
	for _, key := range root.Order {
		normalized := normalizeKey(key)
		if normalized == "global" || normalized == "vdom" || isRecognizedFortiGateSection(normalized) {
			continue
		}
		child := root.Map[key]
		if child == nil || child.Kind != YAMLMapping || len(recognizedFortiGateVDOMSections(child)) == 0 {
			continue
		}
		entries = append(entries, namedNode{
			Name:     key,
			Node:     child,
			Sequence: len(entries) + 1,
			Path:     "$.vdom[detected]",
		})
	}
	if len(entries) > 0 {
		return entries
	}
	if len(recognizedFortiGateVDOMSections(root)) > 0 {
		return []namedNode{{Name: "root", Node: root, Sequence: 1, Path: "$.vdom[direct]"}}
	}
	return nil
}

func recognizedFortiGateSections(node *YAMLNode) []string {
	if node == nil || node.Kind != YAMLMapping {
		return nil
	}
	seen := make(map[string]struct{})
	for _, key := range node.Order {
		normalized := normalizeKey(key)
		if isRecognizedFortiGateSection(normalized) {
			seen[normalized] = struct{}{}
		}
	}
	return sortedKeys(seen)
}

func recognizedFortiGateVDOMSections(node *YAMLNode) []string {
	if node == nil || node.Kind != YAMLMapping {
		return nil
	}
	seen := make(map[string]struct{})
	for _, key := range node.Order {
		normalized := normalizeKey(key)
		if _, ok := fortiGateVDOMSectionNames[normalized]; ok {
			seen[normalized] = struct{}{}
		}
	}
	return sortedKeys(seen)
}

func isRecognizedFortiGateSection(section string) bool {
	if _, ok := fortiGateGlobalSectionNames[section]; ok {
		return true
	}
	_, ok := fortiGateVDOMSectionNames[section]
	return ok
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func yamlKindName(kind YAMLKind) string {
	switch kind {
	case YAMLScalar:
		return "scalar"
	case YAMLMapping:
		return "mapping"
	case YAMLSequence:
		return "sequence"
	default:
		return "unknown"
	}
}
