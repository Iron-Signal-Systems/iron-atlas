package fortigate

import (
	"strings"
	"testing"
)

func TestParseNativePreservesHierarchy(t *testing.T) {
	input := `
#config-version=FGT60F-7.4.0
config system interface
    edit "wan1"
        set ip 192.0.2.2 255.255.255.0
        set role wan
    next
end
config router static
    edit 1
        set dst 0.0.0.0 0.0.0.0
        set gateway 192.0.2.1
        set device "wan1"
    next
end
`
	doc, err := ParseNative(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Root.Children) != 2 {
		t.Fatalf("expected 2 config blocks, got %d", len(doc.Root.Children))
	}
	iface := doc.Root.Children[0].Children[0]
	if iface.Name != "wan1" || iface.Values["set role"][0] != "wan" {
		t.Fatalf("unexpected interface parse: %#v", iface)
	}
}

func TestParseNativeRejectsUnclosedEdit(t *testing.T) {
	_, err := ParseNative(strings.NewReader("config system interface\nedit wan1\nend\n"))
	if err == nil {
		t.Fatal("expected malformed hierarchy error")
	}
}
