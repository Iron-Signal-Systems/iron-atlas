package fortigate

import (
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Iron-Signal-Systems/iron-atlas/modules/firewall/snapshot"
)

type normalizer struct {
	doc      *YAMLDocument
	out      *snapshot.FirewallSnapshot
	objects  map[string][]registeredObject
	findingN int
}

type registeredObject struct {
	ID   snapshot.ObjectID
	Kind snapshot.ObjectKind
}

type namedNode struct {
	Name     string
	Node     *YAMLNode
	Sequence int
	Path     string
}

func NormalizeFortiGateYAML(doc *YAMLDocument) (*snapshot.FirewallSnapshot, error) {
	if doc == nil || doc.Root == nil {
		return nil, fmt.Errorf("YAML document is required")
	}
	n := &normalizer{
		doc:     doc,
		objects: make(map[string][]registeredObject),
		out: &snapshot.FirewallSnapshot{
			Source: snapshot.SnapshotSource{
				Format:         "fortios-yaml",
				Vendor:         "fortinet",
				Platform:       "fortigate",
				FortiOSVersion: detectFortiOSVersion(doc),
			},
			Captured: time.Now().UTC(),
		},
	}

	n.parseGlobal()
	vdoms := n.vdomEntries()
	if len(vdoms) == 0 {
		vdoms = []namedNode{{Name: "root", Node: doc.Root, Sequence: 1, Path: "$.vdom.root"}}
	}
	for _, vdom := range vdoms {
		n.parseVDOM(vdom)
	}
	n.finalize()
	return n.out, nil
}

func (n *normalizer) parseGlobal() {
	global := n.doc.Root.Child("global")
	if global == nil {
		return
	}
	systemGlobal := childAlias(global, "system_global")
	if systemGlobal == nil {
		return
	}
	n.out.Device.Hostname = scalarField(systemGlobal, "hostname")
	n.out.Device.Alias = scalarField(systemGlobal, "alias")
	n.out.Device.SerialNumber = firstNonEmpty(scalarField(systemGlobal, "serial-number"), scalarField(systemGlobal, "serial_number"))
	n.out.Device.Model = scalarField(systemGlobal, "model")
	n.out.Device.Version = firstNonEmpty(scalarField(systemGlobal, "version"), n.out.Source.FortiOSVersion)
	n.out.Device.VDOMMode = scalarField(systemGlobal, "vdom-mode")
}

func (n *normalizer) vdomEntries() []namedNode {
	vdom := n.doc.Root.Child("vdom")
	if vdom == nil {
		return nil
	}
	return tableEntries(vdom, "$.vdom")
}

func (n *normalizer) parseVDOM(entry namedNode) {
	scope := snapshot.ObjectScope{VDOM: entry.Name}
	domainID := objectID(scope, "routing-domain", entry.Name)
	n.out.Domains = append(n.out.Domains, snapshot.RoutingDomain{
		ID:     domainID,
		Name:   entry.Name,
		Scope:  scope,
		Source: location(entry.Node, entry.Path),
	})

	n.parseInterfaces(entry.Node, scope, entry.Path)
	n.parseObjects(entry.Node, scope, entry.Path)
	n.parseNAT(entry.Node, scope, entry.Path)
	n.parseSDWAN(entry.Node, scope, entry.Path)
	n.parseQoS(entry.Node, scope, entry.Path)
	n.parseRoutes(entry.Node, scope, entry.Path)
	n.parseDHCP(entry.Node, scope, entry.Path)
	n.parseVPN(entry.Node, scope, entry.Path)
	n.parsePolicies(entry.Node, scope, entry.Path)
}

func (n *normalizer) register(scope snapshot.ObjectScope, kind snapshot.ObjectKind, name string, id snapshot.ObjectID) {
	if strings.TrimSpace(name) == "" {
		return
	}
	key := referenceKey(scope, kind, name)
	n.objects[key] = append(n.objects[key], registeredObject{ID: id, Kind: kind})
}

func (n *normalizer) ref(owner snapshot.ObjectID, role string, kind snapshot.ObjectKind, name string, scope snapshot.ObjectScope, source snapshot.SourceLocation) snapshot.ObjectRef {
	ref := snapshot.ObjectRef{Kind: kind, VendorName: name, Scope: scope, Source: source}
	if isBuiltinReference(kind, name) {
		ref.Resolution = snapshot.ReferenceBuiltIn
		return ref
	}

	candidates := n.lookup(scope, kind, name)
	if len(candidates) == 1 {
		ref.ID = candidates[0].ID
		ref.Kind = candidates[0].Kind
		ref.Resolution = snapshot.ReferenceResolved
		n.out.References.Edges = append(n.out.References.Edges, snapshot.ReferenceEdge{From: owner, To: ref.ID, Role: role, Source: source})
		return ref
	}
	if len(candidates) > 1 {
		ref.Resolution = snapshot.ReferenceAmbiguous
		n.addUnresolved(owner, role, ref)
		n.addFinding("warning", "reference", "Ambiguous object reference", fmt.Sprintf("%s references %q as %s, but multiple objects match", owner, name, kind), owner, source)
		return ref
	}
	ref.Resolution = snapshot.ReferenceUnresolved
	n.addUnresolved(owner, role, ref)
	n.addFinding("warning", "reference", "Unresolved object reference", fmt.Sprintf("%s references missing %s %q", owner, kind, name), owner, source)
	return ref
}

func (n *normalizer) refs(owner snapshot.ObjectID, role string, kind snapshot.ObjectKind, names []string, scope snapshot.ObjectScope, source snapshot.SourceLocation) []snapshot.ObjectRef {
	refs := make([]snapshot.ObjectRef, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		refs = append(refs, n.ref(owner, role, kind, name, scope, source))
	}
	return refs
}

func (n *normalizer) lookup(scope snapshot.ObjectScope, kind snapshot.ObjectKind, name string) []registeredObject {
	var kinds []snapshot.ObjectKind
	switch kind {
	case snapshot.ObjectKindInterfaceOrZone:
		kinds = []snapshot.ObjectKind{snapshot.ObjectKindInterface, snapshot.ObjectKindZone, snapshot.ObjectKindSDWANZone}
	case snapshot.ObjectKindAddress:
		kinds = []snapshot.ObjectKind{snapshot.ObjectKindAddress, snapshot.ObjectKindAddressGroup, snapshot.ObjectKindVIP}
	case snapshot.ObjectKindService:
		kinds = []snapshot.ObjectKind{snapshot.ObjectKindService, snapshot.ObjectKindServiceGroup}
	default:
		kinds = []snapshot.ObjectKind{kind}
	}
	var result []registeredObject
	for _, candidateKind := range kinds {
		result = append(result, n.objects[referenceKey(scope, candidateKind, name)]...)
		if scope.VDOM != "" {
			result = append(result, n.objects[referenceKey(snapshot.ObjectScope{}, candidateKind, name)]...)
		}
	}
	return deduplicateRegistered(result)
}

func deduplicateRegistered(in []registeredObject) []registeredObject {
	seen := make(map[snapshot.ObjectID]bool)
	out := make([]registeredObject, 0, len(in))
	for _, candidate := range in {
		if seen[candidate.ID] {
			continue
		}
		seen[candidate.ID] = true
		out = append(out, candidate)
	}
	return out
}

func (n *normalizer) addUnresolved(owner snapshot.ObjectID, role string, ref snapshot.ObjectRef) {
	n.out.References.Unresolved = append(n.out.References.Unresolved, snapshot.UnresolvedReference{
		From: owner, Kind: ref.Kind, VendorName: ref.VendorName, Role: role, Scope: ref.Scope, Source: ref.Source,
	})
}

func (n *normalizer) addFinding(severity, category, title, detail string, object snapshot.ObjectID, source snapshot.SourceLocation) {
	n.findingN++
	n.out.Findings = append(n.out.Findings, snapshot.Finding{
		ID: fmt.Sprintf("FGY-%04d", n.findingN), Severity: severity, Category: category,
		Title: title, Detail: detail, ObjectID: object, State: snapshot.EvidenceConfigured, Source: source,
	})
}

func (n *normalizer) finalize() {
	sort.SliceStable(n.out.Policies, func(i, j int) bool {
		if n.out.Policies[i].Scope.VDOM != n.out.Policies[j].Scope.VDOM {
			return n.out.Policies[i].Scope.VDOM < n.out.Policies[j].Scope.VDOM
		}
		return n.out.Policies[i].Sequence < n.out.Policies[j].Sequence
	})
	sort.SliceStable(n.out.Routing.Routes, func(i, j int) bool {
		if n.out.Routing.Routes[i].Scope.VDOM != n.out.Routing.Routes[j].Scope.VDOM {
			return n.out.Routing.Routes[i].Scope.VDOM < n.out.Routing.Routes[j].Scope.VDOM
		}
		return n.out.Routing.Routes[i].Sequence < n.out.Routing.Routes[j].Sequence
	})
	if n.out.Device.Hostname == "" {
		n.addFinding("info", "metadata", "Hostname not present", "The normalized snapshot does not contain a FortiGate hostname", "", snapshot.SourceLocation{YAMLPath: "$.global.system_global.hostname"})
	}
}

func tableEntries(node *YAMLNode, path string) []namedNode {
	if node == nil {
		return nil
	}
	var entries []namedNode
	switch node.Kind {
	case YAMLSequence:
		for i, item := range node.Seq {
			if item.Kind != YAMLMapping {
				continue
			}
			for _, name := range item.Order {
				entries = append(entries, namedNode{Name: name, Node: item.Map[name], Sequence: i + 1, Path: fmt.Sprintf("%s[%d].%s", path, i, name)})
			}
		}
	case YAMLMapping:
		for i, name := range node.Order {
			entries = append(entries, namedNode{Name: name, Node: node.Map[name], Sequence: i + 1, Path: path + "." + name})
		}
	}
	return entries
}

func childAlias(node *YAMLNode, aliases ...string) *YAMLNode {
	if node == nil || node.Kind != YAMLMapping {
		return nil
	}
	for _, alias := range aliases {
		if child := node.Map[alias]; child != nil {
			return child
		}
	}
	for key, child := range node.Map {
		keyNorm := normalizeKey(key)
		for _, alias := range aliases {
			if keyNorm == normalizeKey(alias) {
				return child
			}
		}
	}
	return nil
}

func scalarField(node *YAMLNode, aliases ...string) string {
	child := childAlias(node, aliases...)
	if child == nil {
		return ""
	}
	if child.Kind == YAMLScalar {
		return child.Value
	}
	values := child.Scalars()
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func valuesField(node *YAMLNode, aliases ...string) []string {
	child := childAlias(node, aliases...)
	if child == nil {
		return nil
	}
	if child.Kind == YAMLSequence {
		var values []string
		for _, item := range child.Seq {
			if item.Kind == YAMLScalar {
				values = append(values, item.Value)
			} else if item.Kind == YAMLMapping {
				values = append(values, item.Order...)
			}
		}
		return cleanValues(values)
	}
	if child.Kind == YAMLScalar {
		return cleanValues(splitFortiScalar(child.Value))
	}
	return cleanValues(child.Order)
}

func splitFortiScalar(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	var values []string
	var current strings.Builder
	quoted := rune(0)
	escaped := false
	flush := func() {
		if current.Len() > 0 {
			values = append(values, current.String())
			current.Reset()
		}
	}
	for _, r := range value {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if quoted == '"' && r == '\\' {
			escaped = true
			continue
		}
		if r == '\'' || r == '"' {
			if quoted == 0 {
				quoted = r
			} else if quoted == r {
				quoted = 0
			} else {
				current.WriteRune(r)
			}
			continue
		}
		if unicode.IsSpace(r) && quoted == 0 {
			flush()
			continue
		}
		current.WriteRune(r)
	}
	flush()
	return values
}

func cleanValues(in []string) []string {
	out := make([]string, 0, len(in))
	for _, value := range in {
		value = strings.TrimSpace(strings.Trim(value, "\"'"))
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func enabledField(node *YAMLNode, defaultValue bool, aliases ...string) bool {
	value := strings.ToLower(scalarField(node, aliases...))
	if value == "" {
		return defaultValue
	}
	switch value {
	case "enable", "enabled", "yes", "true", "1", "up":
		return true
	case "disable", "disabled", "no", "false", "0", "down":
		return false
	default:
		return defaultValue
	}
}

func intField(node *YAMLNode, aliases ...string) int {
	value, _ := strconv.Atoi(scalarField(node, aliases...))
	return value
}

func uint32Field(node *YAMLNode, aliases ...string) uint32 {
	value, _ := strconv.ParseUint(scalarField(node, aliases...), 10, 32)
	return uint32(value)
}

func uint64Field(node *YAMLNode, aliases ...string) uint64 {
	value, _ := strconv.ParseUint(scalarField(node, aliases...), 10, 64)
	return value
}

func uint16Value(value string) uint16 {
	parsed, _ := strconv.ParseUint(value, 10, 16)
	return uint16(parsed)
}

func objectID(scope snapshot.ObjectScope, kind, name string) snapshot.ObjectID {
	vdom := scope.VDOM
	if vdom == "" {
		vdom = "global"
	}
	return snapshot.ObjectID(fmt.Sprintf("%s/%s/%s", safeID(vdom), safeID(kind), safeID(name)))
}

func referenceKey(scope snapshot.ObjectScope, kind snapshot.ObjectKind, name string) string {
	return strings.ToLower(scope.VDOM + "\x00" + string(kind) + "\x00" + strings.TrimSpace(name))
}

func safeID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func normalizeKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer("-", "_", " ", "_", ".", "_").Replace(value)
	return value
}

func location(node *YAMLNode, path string) snapshot.SourceLocation {
	if node == nil {
		return snapshot.SourceLocation{YAMLPath: path}
	}
	return snapshot.SourceLocation{YAMLPath: path, Line: node.Line, Column: node.Column}
}

func isBuiltinReference(kind snapshot.ObjectKind, name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	switch kind {
	case snapshot.ObjectKindAddress, snapshot.ObjectKindService, snapshot.ObjectKindInterfaceOrZone:
		return name == "all" || name == "any" || name == "any4" || name == "any6"
	case snapshot.ObjectKindSchedule:
		return name == "always"
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func canonicalIPv4(address, mask string) (string, bool) {
	ip := net.ParseIP(address)
	maskIP := net.ParseIP(mask)
	if ip == nil || maskIP == nil {
		return "", false
	}
	ip4 := ip.To4()
	mask4 := maskIP.To4()
	if ip4 == nil || mask4 == nil {
		return "", false
	}
	ones, bits := net.IPMask(mask4).Size()
	if ones < 0 || bits != 32 {
		return "", false
	}
	prefix, err := netip.ParsePrefix(fmt.Sprintf("%s/%d", ip4.String(), ones))
	if err != nil {
		return "", false
	}
	return prefix.Masked().String(), true
}

func canonicalPrefix(values []string) (prefix, address, mask string, ok bool) {
	if len(values) == 0 {
		return "", "", "", false
	}
	if strings.Contains(values[0], "/") {
		parsed, err := netip.ParsePrefix(values[0])
		if err != nil {
			return "", values[0], "", false
		}
		return parsed.Masked().String(), parsed.Addr().String(), "", true
	}
	if len(values) < 2 {
		return "", values[0], "", false
	}
	prefix, ok = canonicalIPv4(values[0], values[1])
	return prefix, values[0], values[1], ok
}

func parsePortRanges(values []string) []snapshot.PortRange {
	var ranges []snapshot.PortRange
	for _, value := range values {
		for _, token := range strings.FieldsFunc(value, func(r rune) bool { return r == ',' || unicode.IsSpace(r) }) {
			parts := strings.SplitN(token, "-", 2)
			start, err := strconv.ParseUint(parts[0], 10, 16)
			if err != nil {
				continue
			}
			end := start
			if len(parts) == 2 {
				parsedEnd, err := strconv.ParseUint(parts[1], 10, 16)
				if err != nil {
					continue
				}
				end = parsedEnd
			}
			ranges = append(ranges, snapshot.PortRange{Start: uint16(start), End: uint16(end)})
		}
	}
	return ranges
}

func parseDSCP(value string) *snapshot.DSCPValue {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	upper := strings.ToUpper(value)
	if numeric, ok := dscpNames[upper]; ok {
		return &snapshot.DSCPValue{Value: numeric, CanonicalName: upper, VendorName: value, Original: value}
	}
	parsed, err := strconv.ParseUint(value, 0, 8)
	if err != nil || parsed > 63 {
		return nil
	}
	name := canonicalDSCPName(uint8(parsed))
	return &snapshot.DSCPValue{Value: uint8(parsed), CanonicalName: name, VendorName: value, Original: value}
}

var dscpNames = map[string]uint8{
	"CS0": 0, "CS1": 8, "CS2": 16, "CS3": 24, "CS4": 32, "CS5": 40, "CS6": 48, "CS7": 56,
	"AF11": 10, "AF12": 12, "AF13": 14, "AF21": 18, "AF22": 20, "AF23": 22,
	"AF31": 26, "AF32": 28, "AF33": 30, "AF41": 34, "AF42": 36, "AF43": 38,
	"EF": 46,
}

func canonicalDSCPName(value uint8) string {
	for name, candidate := range dscpNames {
		if candidate == value {
			return name
		}
	}
	return ""
}
