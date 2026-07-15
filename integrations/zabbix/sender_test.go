package zabbix

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"
)

func TestEncodeUsesZabbixHeaderAndJSONPayload(t *testing.T) {
	packet, err := Encode([]Metric{{Host: "iron-atlas", Key: "atlas.health", Value: "1", Clock: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if string(packet[:5]) != string(protocolHeader) {
		t.Fatalf("invalid header: %q", packet[:5])
	}
	if got, want := binary.LittleEndian.Uint64(packet[5:13]), uint64(len(packet)-13); got != want {
		t.Fatalf("length %d != %d", got, want)
	}
	var body map[string]any
	if err := json.Unmarshal(packet[13:], &body); err != nil {
		t.Fatal(err)
	}
	if body["request"] != "sender data" {
		t.Fatalf("unexpected request: %#v", body)
	}
}

func TestDecodeResponse(t *testing.T) {
	payload := []byte(`{"response":"success","info":"processed: 1; failed: 0"}`)
	packet := append([]byte{}, protocolHeader...)
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(payload)))
	packet = append(packet, length...)
	packet = append(packet, payload...)
	response, err := DecodeResponse(bytes.NewReader(packet))
	if err != nil {
		t.Fatal(err)
	}
	if response.Response != "success" {
		t.Fatalf("unexpected response: %#v", response)
	}
}
