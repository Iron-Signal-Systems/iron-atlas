package zabbix

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var protocolHeader = []byte{'Z', 'B', 'X', 'D', 1}

type Metric struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int64  `json:"clock,omitempty"`
	NS    int    `json:"ns,omitempty"`
}

type request struct {
	Request string   `json:"request"`
	Data    []Metric `json:"data"`
}

type Response struct {
	Response string `json:"response"`
	Info     string `json:"info"`
}

type Client struct {
	Address   string
	Timeout   time.Duration
	TLSConfig *tls.Config
}

func Encode(metrics []Metric) ([]byte, error) {
	if len(metrics) == 0 {
		return nil, errors.New("at least one metric is required")
	}
	for i, metric := range metrics {
		if metric.Host == "" || metric.Key == "" {
			return nil, fmt.Errorf("metric %d requires host and key", i)
		}
	}
	payload, err := json.Marshal(request{Request: "sender data", Data: metrics})
	if err != nil {
		return nil, err
	}
	packet := make([]byte, 13+len(payload))
	copy(packet, protocolHeader)
	binary.LittleEndian.PutUint64(packet[5:13], uint64(len(payload)))
	copy(packet[13:], payload)
	return packet, nil
}

func DecodeResponse(reader io.Reader) (Response, error) {
	header := make([]byte, 13)
	if _, err := io.ReadFull(reader, header); err != nil {
		return Response{}, err
	}
	if string(header[:5]) != string(protocolHeader) {
		return Response{}, errors.New("invalid Zabbix response header")
	}
	size := binary.LittleEndian.Uint64(header[5:13])
	if size > 1<<20 {
		return Response{}, errors.New("Zabbix response exceeds safety limit")
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return Response{}, err
	}
	var response Response
	if err := json.Unmarshal(payload, &response); err != nil {
		return Response{}, err
	}
	return response, nil
}

func (c Client) Send(ctx context.Context, metrics []Metric) (Response, error) {
	packet, err := Encode(metrics)
	if err != nil {
		return Response{}, err
	}
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	dialer := &net.Dialer{Timeout: timeout}
	var conn net.Conn
	if c.TLSConfig != nil {
		conn, err = tls.DialWithDialer(dialer, "tcp", c.Address, c.TLSConfig.Clone())
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", c.Address)
	}
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))
	if _, err := conn.Write(packet); err != nil {
		return Response{}, err
	}
	return DecodeResponse(bufio.NewReader(conn))
}
