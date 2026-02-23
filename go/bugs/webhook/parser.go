package webhook

import (
	"regexp"
	"strings"
)

// issueRefPattern matches patterns like:
// - "Fixes #42", "Closes #42", "Resolves #42"
// - "Fixes L8B-42", "Closes L8B-42"
// - "Fixes <uuid>", "Resolves <uuid>"
var issueRefPattern = regexp.MustCompile(
	`(?i)(?:fix(?:es|ed)?|close[sd]?|resolve[sd]?)\s+` +
		`(#\d+|[A-Za-z]+-\d+|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

// ExtractIssueRefs extracts issue references from text.
// Returns deduplicated list of reference strings (without the # prefix for numeric refs).
func ExtractIssueRefs(text string) []string {
	matches := issueRefPattern.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var refs []string
	for _, m := range matches {
		ref := strings.TrimPrefix(m[1], "#")
		if !seen[ref] {
			seen[ref] = true
			refs = append(refs, ref)
		}
	}
	return refs
}
