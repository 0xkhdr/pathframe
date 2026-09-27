package delegation

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	IntegrationSchema               = "pathframe.integration/v1"
	CapabilitySequentialSubagent    = "sequential_subagent"
	CapabilitySharedWorkspace       = "shared_workspace"
	CapabilityResultReturn          = "structured_result_return"
	CapabilityWriteScopeEnforcement = "write_scope_enforcement"
)

type CapabilityDeclaration struct {
	Schema       string   `json:"schema"`
	Version      string   `json:"version"`
	Host         string   `json:"host"`
	Transport    string   `json:"transport"`
	Capabilities []string `json:"capabilities"`
}

func LoadCapabilities(path, host string) (CapabilityDeclaration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return CapabilityDeclaration{}, err
	}
	var declaration CapabilityDeclaration
	if err := json.Unmarshal(data, &declaration); err != nil {
		return CapabilityDeclaration{}, err
	}
	if declaration.Schema != IntegrationSchema || declaration.Version != "1.0.0" || declaration.Host != host || declaration.Transport != "stdio" {
		return CapabilityDeclaration{}, fmt.Errorf("incompatible %s integration manifest", host)
	}
	return declaration, nil
}

func (c CapabilityDeclaration) Has(capability string) bool {
	for _, candidate := range c.Capabilities {
		if candidate == capability {
			return true
		}
	}
	return false
}
