package telemetry

import "time"

type Severity string

const (
	SeverityCritical      Severity = "critical"
	SeverityHigh          Severity = "high"
	SeverityModerate      Severity = "moderate"
	SeverityLow           Severity = "low"
	SeverityInformational Severity = "informational"
)

type Metric struct {
	Name      string            `json:"name"`
	Value     string            `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels,omitempty"`
	DeviceID  string            `json:"device_id,omitempty"`
	SiteID    string            `json:"site_id,omitempty"`
	ModuleID  string            `json:"module_id,omitempty"`
}

type DeliveryState string

const (
	DeliveryPending   DeliveryState = "pending"
	DeliveryDelivered DeliveryState = "delivered"
	DeliveryFailed    DeliveryState = "failed"
)
