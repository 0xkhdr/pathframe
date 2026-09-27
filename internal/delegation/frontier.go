package delegation

import (
	"fmt"
	"sort"

	"github.com/0xkhdr/pathframe/internal/artifacts"
)

type Projection struct {
	Next        string     `json:"next,omitempty"`
	Frontier    []string   `json:"frontier"`
	Waves       [][]string `json:"waves"`
	Diagnostics []string   `json:"diagnostics"`
	Recovery    []string   `json:"recovery"`
}

func Project(tasks []artifacts.Task, completed []string) Projection {
	done := map[string]bool{}
	for _, id := range completed {
		done[id] = true
	}
	remaining := map[string]artifacts.Task{}
	for _, task := range tasks {
		if !done[task.ID] {
			remaining[task.ID] = task
		}
	}
	result := Projection{Frontier: []string{}, Waves: [][]string{}, Diagnostics: []string{}, Recovery: []string{}}
	for len(remaining) > 0 {
		wave := []string{}
		for id, task := range remaining {
			ready := true
			for _, dependency := range task.Dependencies {
				if !done[dependency] {
					ready = false
					break
				}
			}
			if ready {
				wave = append(wave, id)
			}
		}
		sort.Strings(wave)
		if len(wave) == 0 {
			ids := make([]string, 0, len(remaining))
			for id := range remaining {
				ids = append(ids, id)
			}
			sort.Strings(ids)
			result.Diagnostics = append(result.Diagnostics, fmt.Sprintf("dependency cycle blocks tasks: %v", ids))
			result.Recovery = []string{"replan task dependencies and request human reapproval"}
			return result
		}
		if len(result.Waves) == 0 {
			result.Frontier = append(result.Frontier, wave...)
			result.Next = wave[0]
		}
		result.Waves = append(result.Waves, wave)
		for _, id := range wave {
			delete(remaining, id)
			done[id] = true
		}
	}
	return result
}
