package workflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestStageOneSchemasAreValidJSON(t *testing.T) {
	for _, name := range []string{"workflow-v1.schema.json", "transition-event-v1.schema.json"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "schemas", name))
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if schema["$schema"] == nil || schema["$id"] == nil {
			t.Fatalf("%s lacks schema identity", name)
		}
	}
}
