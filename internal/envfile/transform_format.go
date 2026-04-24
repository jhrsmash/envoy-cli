package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatTransformResult returns a human-readable summary of a TransformResult.
func FormatTransformResult(r TransformResult, op TransformOp) string {
	var sb strings.Builder

	if len(r.Changed) == 0 {
		sb.WriteString(fmt.Sprintf("transform(%s): no values changed\n", op))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("transform(%s): %d value(s) changed\n", op, len(r.Changed)))

	keys := make([]string, len(r.Changed))
	copy(keys, r.Changed)
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  ~ %s = %q\n", k, r.Output[k]))
	}

	return sb.String()
}
