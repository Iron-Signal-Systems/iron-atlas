package config

import (
	"errors"
	"fmt"
	"net"
)

type File struct {
	ListenAddress string              `json:"listen_address"`
	SiteID        string              `json:"site_id"`
	DataDirectory string              `json:"data_directory"`
	EvidenceStore EvidenceStoreConfig `json:"evidence_store"`
	Zabbix        ZabbixConfig        `json:"zabbix"`
}

type EvidenceStoreConfig struct {
	Path string `json:"path"`
}

type ZabbixConfig struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"address"`
	Host    string `json:"host"`
}

func (f File) Validate() error {
	if f.ListenAddress == "" {
		return errors.New("listen_address is required")
	}
	if _, _, err := net.SplitHostPort(f.ListenAddress); err != nil {
		return fmt.Errorf("listen_address: %w", err)
	}
	if f.SiteID == "" {
		return errors.New("site_id is required")
	}
	if f.DataDirectory == "" {
		return errors.New("data_directory is required")
	}
	if f.EvidenceStore.Path == "" {
		return errors.New("evidence_store.path is required")
	}
	if f.Zabbix.Enabled {
		if f.Zabbix.Address == "" || f.Zabbix.Host == "" {
			return errors.New("zabbix.address and zabbix.host are required when enabled")
		}
		if _, _, err := net.SplitHostPort(f.Zabbix.Address); err != nil {
			return fmt.Errorf("zabbix.address: %w", err)
		}
	}
	return nil
}
