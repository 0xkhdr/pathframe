package verification

import (
	"path"
	"sort"
	"strings"
)

func ScopeViolations(changed, scope []string) []string {
	var violations []string
	for _, name := range changed {
		name = strings.TrimPrefix(strings.ReplaceAll(name, "\\", "/"), "./")
		allowed := false
		for _, pattern := range scope {
			pattern = strings.TrimPrefix(strings.ReplaceAll(pattern, "\\", "/"), "./")
			if pattern == name || (strings.HasSuffix(pattern, "/**") && strings.HasPrefix(name, strings.TrimSuffix(pattern, "**"))) {
				allowed = true
				break
			}
			if matched, _ := path.Match(pattern, name); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			violations = append(violations, name)
		}
	}
	sort.Strings(violations)
	return violations
}
