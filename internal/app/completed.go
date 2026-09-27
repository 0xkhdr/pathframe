package app

import (
	"path/filepath"

	"github.com/0xkhdr/pathframe/internal/artifacts"
	"github.com/0xkhdr/pathframe/internal/delegation"
	"github.com/0xkhdr/pathframe/internal/recovery"
)

func completedTasks(changeDir string, plan artifacts.Plan) ([]string, error) {
	identities := make(map[string]string, len(plan.Tasks))
	for _, task := range plan.Tasks {
		identities[task.ID] = recovery.TaskIdentity(task)
	}
	return delegation.CompletedTasksMatching(filepath.Join(changeDir, "runs"), identities)
}
