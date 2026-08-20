package types

import (
	"slices"
	"strings"
)

func normalizeFieldMaskPaths(paths []string) string {
	if slices.Contains(paths, "*") {
		return "*"
	}
	paths = slices.Clone(paths)
	slices.Sort(paths)
	paths = slices.Compact(paths)
	normalized := paths[:0]
	for _, path := range paths {
		if len(normalized) > 0 && strings.HasPrefix(path, normalized[len(normalized)-1]+".") {
			continue
		}
		normalized = append(normalized, path)
	}
	return strings.Join(normalized, ",")
}
