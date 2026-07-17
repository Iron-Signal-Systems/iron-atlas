package main

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/Iron-Signal-Systems/iron-atlas/modules/firewall/fortigate"
	"github.com/Iron-Signal-Systems/iron-atlas/modules/firewall/snapshot"
)

const (
	unclassifiedRole     = "unclassified-role"
	unclassifiedKind     = "unclassified-kind"
	unclassifiedSeverity = "unclassified-severity"
	unclassifiedCategory = "unclassified-category"
	unclassifiedTitle    = "Unclassified finding title"
)

type roleKind struct {
	Role string
	Kind string
}

type semanticQualityReport struct {
	Layout             fortigate.FortiGateYAMLLayout
	RecordCounts       []snapshot.RecordCount
	Resolved           map[roleKind]int
	BuiltIns           map[string]int
	Unresolved         map[roleKind]int
	Ambiguous          map[roleKind]int
	FindingSeverities  map[string]int
	FindingCategories  map[string]int
	FindingTitles      map[string]int
	SafeLabelFallbacks int
}

func buildSemanticQualityReport(value *snapshot.FirewallSnapshot, layout fortigate.FortiGateYAMLLayout) semanticQualityReport {
	report := semanticQualityReport{
		Layout:            layout,
		RecordCounts:      snapshot.NormalizedRecordCounts(value),
		Resolved:          make(map[roleKind]int),
		BuiltIns:          make(map[string]int),
		Unresolved:        make(map[roleKind]int),
		Ambiguous:         make(map[roleKind]int),
		FindingSeverities: make(map[string]int),
		FindingCategories: make(map[string]int),
		FindingTitles:     make(map[string]int),
	}
	if value == nil {
		return report
	}

	for _, edge := range value.References.Edges {
		role, roleFallback := safeReferenceRole(edge.Role)
		kind, kindFallback := safeObjectKind(edge.Kind)
		report.Resolved[roleKind{Role: role, Kind: kind}]++
		report.SafeLabelFallbacks += boolInt(roleFallback) + boolInt(kindFallback)
	}
	for _, builtIn := range value.References.BuiltIns {
		kind, fallback := safeObjectKind(builtIn.Kind)
		report.BuiltIns[kind]++
		report.SafeLabelFallbacks += boolInt(fallback)
	}
	for _, unresolved := range value.References.Unresolved {
		role, roleFallback := safeReferenceRole(unresolved.Role)
		kind, kindFallback := safeObjectKind(unresolved.Kind)
		key := roleKind{Role: role, Kind: kind}
		if unresolved.Resolution == snapshot.ReferenceAmbiguous {
			report.Ambiguous[key]++
		} else {
			report.Unresolved[key]++
		}
		report.SafeLabelFallbacks += boolInt(roleFallback) + boolInt(kindFallback)
	}
	for _, finding := range value.Findings {
		severity, severityFallback := safeFindingSeverity(finding.Severity)
		category, categoryFallback := safeFindingCategory(finding.Category)
		title, titleFallback := safeFindingTitle(finding.Title)
		report.FindingSeverities[severity]++
		report.FindingCategories[category]++
		report.FindingTitles[title]++
		report.SafeLabelFallbacks += boolInt(severityFallback) + boolInt(categoryFallback) + boolInt(titleFallback)
	}
	return report
}

func renderSemanticQuality(value *snapshot.FirewallSnapshot, layout fortigate.FortiGateYAMLLayout) string {
	report := buildSemanticQualityReport(value, layout)
	var b bytes.Buffer
	fmt.Fprintln(&b, "Upload-safe FortiGate semantic-quality report")
	fmt.Fprintln(&b, "Privacy mode: aggregate-only")

	fmt.Fprintln(&b, "\nLayout coverage:")
	fmt.Fprintf(&b, "- Root entries: %d\n", report.Layout.RootEntries)
	fmt.Fprintf(&b, "- Recognized direct section kinds: %d\n", len(report.Layout.RecognizedDirectSections))
	fmt.Fprintf(&b, "- Recognized nested section kinds: %d\n", len(report.Layout.RecognizedNestedSections))
	fmt.Fprintf(&b, "- Detected VDOM containers: %d\n", report.Layout.DetectedVDOMContainerCount)
	fmt.Fprintf(&b, "- Unrecognized root entries: %d\n", report.Layout.UnrecognizedRootEntryCount)

	fmt.Fprintln(&b, "\nNormalized records by kind:")
	fmt.Fprintf(&b, "- total: %d\n", snapshot.TotalNormalizedRecords(value))
	for _, record := range report.RecordCounts {
		fmt.Fprintf(&b, "- %s: %d\n", record.Kind, record.Count)
	}

	fmt.Fprintln(&b, "\nReference totals:")
	fmt.Fprintf(&b, "- resolved: %d\n", sumRoleKindCounts(report.Resolved))
	fmt.Fprintf(&b, "- built-in: %d\n", sumStringCounts(report.BuiltIns))
	fmt.Fprintf(&b, "- unresolved: %d\n", sumRoleKindCounts(report.Unresolved))
	fmt.Fprintf(&b, "- ambiguous: %d\n", sumRoleKindCounts(report.Ambiguous))

	renderRoleKindCounts(&b, "Resolved references by role and kind", report.Resolved)
	renderStringCounts(&b, "Built-in references by kind", report.BuiltIns)
	renderRoleKindCounts(&b, "Unresolved references by role and kind", report.Unresolved)
	renderRoleKindCounts(&b, "Ambiguous references by role and kind", report.Ambiguous)

	fmt.Fprintln(&b, "\nFinding totals:")
	if value == nil {
		fmt.Fprintln(&b, "- total: 0")
	} else {
		fmt.Fprintf(&b, "- total: %d\n", len(value.Findings))
	}
	renderStringCounts(&b, "Findings by severity", report.FindingSeverities)
	renderStringCounts(&b, "Findings by category", report.FindingCategories)
	renderStringCounts(&b, "Findings by title", report.FindingTitles)

	fmt.Fprintln(&b, "\nCoverage warnings:")
	warnings := 0
	if report.Layout.UnrecognizedRootEntryCount > 0 {
		fmt.Fprintf(&b, "- unrecognized-root-entry-coverage: %d\n", report.Layout.UnrecognizedRootEntryCount)
		warnings++
	}
	if count := sumRoleKindCounts(report.Unresolved); count > 0 {
		fmt.Fprintf(&b, "- unresolved-reference-coverage: %d\n", count)
		warnings++
	}
	if count := sumRoleKindCounts(report.Ambiguous); count > 0 {
		fmt.Fprintf(&b, "- ambiguous-reference-coverage: %d\n", count)
		warnings++
	}
	if value != nil && len(value.Findings) > 0 {
		fmt.Fprintf(&b, "- finding-review-required: %d\n", len(value.Findings))
		warnings++
	}
	if report.SafeLabelFallbacks > 0 {
		fmt.Fprintf(&b, "- unclassified-safe-label-buckets: %d\n", report.SafeLabelFallbacks)
		warnings++
	}
	if warnings == 0 {
		fmt.Fprintln(&b, "- none")
	}

	fmt.Fprintln(&b, "\nPrivate vendor names, object identifiers, source paths, scalar values, finding details, VDOM names, input paths, and normalized JSON: omitted")
	return b.String()
}

func renderRoleKindCounts(b *bytes.Buffer, title string, counts map[roleKind]int) {
	fmt.Fprintf(b, "\n%s:\n", title)
	if len(counts) == 0 {
		fmt.Fprintln(b, "- none")
		return
	}
	keys := make([]roleKind, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Role != keys[j].Role {
			return keys[i].Role < keys[j].Role
		}
		return keys[i].Kind < keys[j].Kind
	})
	for _, key := range keys {
		fmt.Fprintf(b, "- %s | %s: %d\n", key.Role, key.Kind, counts[key])
	}
}

func renderStringCounts(b *bytes.Buffer, title string, counts map[string]int) {
	fmt.Fprintf(b, "\n%s:\n", title)
	if len(counts) == 0 {
		fmt.Fprintln(b, "- none")
		return
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(b, "- %s: %d\n", key, counts[key])
	}
}

func sumRoleKindCounts(counts map[roleKind]int) int {
	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}

func sumStringCounts(counts map[string]int) int {
	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}

func safeReferenceRole(value string) (string, bool) {
	if _, ok := safeReferenceRoles[value]; ok {
		return value, false
	}
	return unclassifiedRole, true
}

func safeObjectKind(value snapshot.ObjectKind) (string, bool) {
	if _, ok := safeObjectKinds[value]; ok {
		return string(value), false
	}
	return unclassifiedKind, true
}

func safeFindingSeverity(value string) (string, bool) {
	if _, ok := safeFindingSeverities[value]; ok {
		return value, false
	}
	return unclassifiedSeverity, true
}

func safeFindingCategory(value string) (string, bool) {
	if _, ok := safeFindingCategories[value]; ok {
		return value, false
	}
	return unclassifiedCategory, true
}

func safeFindingTitle(value string) (string, bool) {
	if _, ok := safeFindingTitles[value]; ok {
		return value, false
	}
	return unclassifiedTitle, true
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

var safeReferenceRoles = map[string]struct{}{
	"address-associated-interface":       {},
	"address-group-exclude-member":       {},
	"address-group-member":               {},
	"central-snat-destination-address":   {},
	"central-snat-destination-interface": {},
	"central-snat-ip-pool":               {},
	"central-snat-source-address":        {},
	"central-snat-source-interface":      {},
	"dhcp-interface":                     {},
	"parent-interface":                   {},
	"policy-destination-address":         {},
	"policy-destination-interface":       {},
	"policy-forward-shaper":              {},
	"policy-ip-pool":                     {},
	"policy-per-ip-shaper":               {},
	"policy-reverse-shaper":              {},
	"policy-route-egress":                {},
	"policy-route-ingress":               {},
	"policy-schedule":                    {},
	"policy-service":                     {},
	"policy-source-address":              {},
	"policy-source-interface":            {},
	"policy-user":                        {},
	"policy-user-group":                  {},
	"route-egress":                       {},
	"route-map-interface":                {},
	"route-map-prefix-list":              {},
	"sdwan-destination-address":          {},
	"sdwan-health-member":                {},
	"sdwan-member-interface":             {},
	"sdwan-member-zone":                  {},
	"sdwan-preferred-member":             {},
	"sdwan-required-sla":                 {},
	"sdwan-service":                      {},
	"sdwan-source-address":               {},
	"service-group-member":               {},
	"shaping-destination-address":        {},
	"shaping-destination-interface":      {},
	"shaping-forward-shaper":             {},
	"shaping-per-ip-shaper":              {},
	"shaping-reverse-shaper":             {},
	"shaping-service":                    {},
	"shaping-source-address":             {},
	"shaping-source-interface":           {},
	"vip-external-interface":             {},
	"vip-group-member":                   {},
	"vpn-interface":                      {},
	"zone-member":                        {},
}

var safeObjectKinds = map[snapshot.ObjectKind]struct{}{
	snapshot.ObjectKindInterface:       {},
	snapshot.ObjectKindZone:            {},
	snapshot.ObjectKindInterfaceOrZone: {},
	snapshot.ObjectKindVLAN:            {},
	snapshot.ObjectKindSubnet:          {},
	snapshot.ObjectKindAddress:         {},
	snapshot.ObjectKindAddressGroup:    {},
	snapshot.ObjectKindService:         {},
	snapshot.ObjectKindServiceGroup:    {},
	snapshot.ObjectKindSchedule:        {},
	snapshot.ObjectKindUser:            {},
	snapshot.ObjectKindUserGroup:       {},
	snapshot.ObjectKindVIP:             {},
	snapshot.ObjectKindIPPool:          {},
	snapshot.ObjectKindTrafficShaper:   {},
	snapshot.ObjectKindPerIPShaper:     {},
	snapshot.ObjectKindSDWANMember:     {},
	snapshot.ObjectKindSDWANZone:       {},
	snapshot.ObjectKindSDWANHealth:     {},
	snapshot.ObjectKindRouteMap:        {},
	snapshot.ObjectKindPrefixList:      {},
	snapshot.ObjectKindVPN:             {},
	snapshot.ObjectKindTag:             {},
}

var safeFindingSeverities = map[string]struct{}{
	"info":    {},
	"warning": {},
	"error":   {},
}

var safeFindingCategories = map[string]struct{}{
	"interface": {},
	"metadata":  {},
	"object":    {},
	"policy":    {},
	"reference": {},
	"routing":   {},
}

var safeFindingTitles = map[string]struct{}{
	"Ambiguous object reference":  {},
	"Disabled firewall policy":    {},
	"Disabled static route":       {},
	"Empty address group":         {},
	"Hostname not present":        {},
	"Invalid interface address":   {},
	"Unresolved object reference": {},
}
