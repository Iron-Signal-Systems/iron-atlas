package snapshot

import "testing"

func TestNormalizedRecordCountsAreStableAndComplete(t *testing.T) {
	value := &FirewallSnapshot{
		Domains:    []RoutingDomain{{}},
		Interfaces: []Interface{{}, {}},
		Routing: RoutingConfiguration{
			BGP:    &RoutingProtocol{},
			Routes: []Route{{}},
		},
		Objects: ObjectCatalog{Addresses: []AddressObject{{}, {}, {}}},
	}

	counts := NormalizedRecordCounts(value)
	if len(counts) == 0 {
		t.Fatal("expected fixed record counts")
	}
	for index := 1; index < len(counts); index++ {
		if counts[index-1].Kind >= counts[index].Kind {
			t.Fatalf("record kinds are not in stable order: %q before %q", counts[index-1].Kind, counts[index].Kind)
		}
	}

	got := make(map[string]int, len(counts))
	for _, count := range counts {
		got[count.Kind] = count.Count
	}
	for kind, want := range map[string]int{
		"address":           3,
		"bgp-configuration": 1,
		"interface":         2,
		"route":             1,
		"routing-domain":    1,
		"policy-route":      0,
	} {
		if got[kind] != want {
			t.Fatalf("unexpected %s count: got %d want %d", kind, got[kind], want)
		}
	}
	if got, want := TotalNormalizedRecords(value), 8; got != want {
		t.Fatalf("unexpected total normalized records: got %d want %d", got, want)
	}
}

func TestNormalizedRecordCountsHandleNil(t *testing.T) {
	if got := TotalNormalizedRecords(nil); got != 0 {
		t.Fatalf("unexpected nil total: %d", got)
	}
	if counts := NormalizedRecordCounts(nil); len(counts) == 0 {
		t.Fatal("nil snapshot should still return fixed zero-count kinds")
	}
}
