package errorx

import "strings"

func joinStringsIfNonEmpty(delimiter string, parts ...string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		if len(parts[0]) == 0 {
			return parts[1]
		} else if len(parts[1]) == 0 {
			return parts[0]
		} else {
			return parts[0] + delimiter + parts[1]
		}
	default:
		filteredParts := make([]string, 0, len(parts))
		for _, part := range parts {
			if len(part) > 0 {
				filteredParts = append(filteredParts, part)
			}
		}

		return strings.Join(filteredParts, delimiter)
	}
}

