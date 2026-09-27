package context

import (
	"fmt"
	"strings"

	"github.com/0xkhdr/pathframe/internal/artifacts"
)

func Assemble(change string, task artifacts.Task, role Role, resolved []Entry, unresolved []Omission, runtime string, limit int) (Packet, error) {
	entries := append([]Entry{}, resolved...)
	taskContent := fmt.Sprintf("%s: %s\n%s", task.ID, task.Title, task.Sections["Objective"])
	entries = append(entries,
		Entry{Layer: Task, Path: "tasks/" + task.ID, Required: true, Content: taskContent, Bytes: len(taskContent)},
		Entry{Layer: Runtime, Path: "workflow", Required: true, Content: runtime, Bytes: len(runtime)},
	)
	included, omissions, budget, err := ApplyBudget(entries, limit)
	if err != nil {
		return Packet{}, err
	}
	return Packet{
		Schema: PacketSchema, Change: change, Task: task.ID, Title: task.Title, ExecutionPolicy: task.ExecutionPolicy,
		Role: role, Objective: task.Sections["Objective"], Acceptance: task.Sections["Acceptance"], References: nonnil(task.References),
		Context: included, Dependencies: nonnil(task.Dependencies), CompletedPrerequisites: []string{}, WriteScope: bullets(task.Sections["Write Scope"]),
		WriteScopeAssurance: "advisory", Constraints: bullets(task.Sections["Constraints"]), Verification: task.Verification,
		ResultSchema: "pathframe.task-result/v1", HostAssurance: "not_evaluated", Omissions: append(unresolved, omissions...), Budget: budget,
	}, nil
}

func bullets(body string) []string {
	values := []string{}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if value != "" && value != "none" {
				values = append(values, value)
			}
		}
	}
	return values
}

func nonnil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
