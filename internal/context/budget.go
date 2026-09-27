package context

import "fmt"

func ApplyBudget(entries []Entry, limit int) ([]Entry, []Omission, Budget, error) {
	if limit <= 0 {
		return nil, nil, Budget{}, fmt.Errorf("context budget must be positive")
	}
	budget := Budget{LimitBytes: limit}
	for _, entry := range entries {
		if entry.Required {
			budget.RequiredBytes += entry.Bytes
		}
	}
	if budget.RequiredBytes > limit {
		return nil, nil, budget, fmt.Errorf("required context is %d bytes, exceeding the %d-byte budget; split the task or explicitly raise the budget", budget.RequiredBytes, limit)
	}
	included := make([]Entry, 0, len(entries))
	omissions := []Omission{}
	optionalBytes := 0
	for _, entry := range entries {
		if entry.Required || budget.RequiredBytes+optionalBytes+entry.Bytes <= limit {
			included = append(included, entry)
			budget.IncludedBytes += entry.Bytes
			if !entry.Required {
				optionalBytes += entry.Bytes
			}
			continue
		}
		omissions = append(omissions, Omission{Layer: entry.Layer, Path: entry.Path, Reason: "optional context exceeds budget", Bytes: entry.Bytes})
		budget.OmittedBytes += entry.Bytes
	}
	return included, omissions, budget, nil
}
