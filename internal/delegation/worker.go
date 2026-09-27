package delegation

import (
	"encoding/json"
	"fmt"

	pfcontext "github.com/0xkhdr/pathframe/internal/context"
)

func WorkerInstructions(packet pfcontext.Packet) string {
	data, _ := json.MarshalIndent(packet, "", "  ")
	return fmt.Sprintf("You are Pinky. Execute only this delegated packet. Stay within its role and write scope; scope is advisory unless marked host_enforced. Do not change Pathframe plan state, accept your work, or broaden execution_policy. Return only pathframe.task-result/v1 JSON for Brain submission.\n\n%s", data)
}
