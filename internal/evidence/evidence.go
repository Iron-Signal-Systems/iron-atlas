package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"
)

type Classification string

const (
	ClassificationSensitiveInfrastructure Classification = "sensitive_infrastructure"
	ClassificationRedactedDocumentation   Classification = "redacted_documentation"
)

type Manifest struct {
	CollectionID   string         `json:"collection_id"`
	DeviceID       string         `json:"device_id"`
	CollectionType string         `json:"collection_type"`
	CollectedAt    time.Time      `json:"collected_at"`
	SHA256         string         `json:"sha256"`
	Classification Classification `json:"classification"`
	ParserVersion  string         `json:"parser_version"`
}

func SHA256(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
