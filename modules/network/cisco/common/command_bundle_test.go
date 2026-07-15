package common

import (
	"strings"
	"testing"
)

func TestParseCommandBundle(t *testing.T) {
	input := "===== COMMAND: show version =====\nCisco IOS Software\n===== END COMMAND =====\n===== COMMAND: show interfaces trunk =====\nGi1/0/48 on 802.1q trunking 999\n===== END COMMAND =====\n"
	bundle, err := ParseCommandBundle(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle) != 2 || bundle["show version"] != "Cisco IOS Software" {
		t.Fatalf("unexpected bundle: %#v", bundle)
	}
}
