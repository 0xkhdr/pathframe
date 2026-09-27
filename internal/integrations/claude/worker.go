package claude

import (
	pfcontext "github.com/0xkhdr/pathframe/internal/context"
	"github.com/0xkhdr/pathframe/internal/delegation"
)

func WorkerInstructions(packet pfcontext.Packet) string {
	return delegation.WorkerInstructions(packet)
}
