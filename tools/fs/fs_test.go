package fs_test

import "strings"

func isErrorResult(result string) bool {
	return strings.HasPrefix(result, "Cannot") || strings.Contains(result, "due to error")
}
