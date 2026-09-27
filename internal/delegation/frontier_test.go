package delegation

import (
	"reflect"
	"testing"

	"github.com/0xkhdr/pathframe/internal/artifacts"
)

func TestProjectDeterministicWavesAndCycle(t *testing.T) {
	tasks := []artifacts.Task{{ID: "T3", Dependencies: []string{"T1", "T2"}}, {ID: "T2", Dependencies: []string{"T1"}}, {ID: "T1"}}
	got := Project(tasks, nil)
	if got.Next != "T1" || !reflect.DeepEqual(got.Waves, [][]string{{"T1"}, {"T2"}, {"T3"}}) || !reflect.DeepEqual(got.Frontier, []string{"T1"}) {
		t.Fatalf("Project() = %#v", got)
	}
	cycle := Project([]artifacts.Task{{ID: "T1", Dependencies: []string{"T2"}}, {ID: "T2", Dependencies: []string{"T1"}}}, nil)
	if len(cycle.Diagnostics) != 1 || len(cycle.Recovery) != 1 {
		t.Fatalf("cycle = %#v", cycle)
	}
}
