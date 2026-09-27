package delegation

import (
	"fmt"

	pfcontext "github.com/0xkhdr/pathframe/internal/context"
)

type Preflight struct {
	Ready               bool     `json:"ready"`
	Host                string   `json:"host"`
	WriteScopeAssurance string   `json:"write_scope_assurance"`
	Issues              []string `json:"issues"`
	Recovery            []string `json:"recovery"`
}

func Check(packet pfcontext.Packet, capabilities CapabilityDeclaration, leaseActive bool) Preflight {
	result := Preflight{Host: capabilities.Host, WriteScopeAssurance: "advisory", Issues: []string{}, Recovery: []string{}}
	if packet.ExecutionPolicy != "delegated" {
		result.Issues = append(result.Issues, "task execution_policy is not delegated")
		result.Recovery = append(result.Recovery, "execute the task as Brain or amend the task and request human reapproval")
	}
	if packet.Role.ID == "" || packet.Role.ID == "none" {
		result.Issues = append(result.Issues, "delegated task requires an explicit role")
		result.Recovery = append(result.Recovery, "assign an existing role and request human reapproval")
	}
	for _, capability := range []string{CapabilitySequentialSubagent, CapabilitySharedWorkspace, CapabilityResultReturn} {
		if !capabilities.Has(capability) {
			result.Issues = append(result.Issues, fmt.Sprintf("host does not declare required capability %s", capability))
			result.Recovery = append(result.Recovery, "repair or reinstall the host integration, then retry delegation")
		}
	}
	if len(packet.WriteScope) == 0 {
		result.Issues = append(result.Issues, "delegated task requires a non-empty write scope")
		result.Recovery = append(result.Recovery, "add write scope and request human reapproval")
	}
	if len(packet.Verification) == 0 {
		result.Issues = append(result.Issues, "delegated task requires verification argv")
		result.Recovery = append(result.Recovery, "add verification and request human reapproval")
	}
	if leaseActive {
		result.Issues = append(result.Issues, "a task lease is already active")
		result.Recovery = append(result.Recovery, "submit the active result or explicitly release the lost lease before retrying")
	}
	if capabilities.Has(CapabilityWriteScopeEnforcement) {
		result.WriteScopeAssurance = "host_enforced"
	}
	result.Ready = len(result.Issues) == 0
	return result
}
