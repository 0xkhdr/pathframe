package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/0xkhdr/pathframe/internal/artifacts"
)

func TaskIdentity(task artifacts.Task) string {
	data, _ := json.Marshal(task)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
