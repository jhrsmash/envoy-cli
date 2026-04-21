package envfile

import (
	"fmt"
	"strings"
)

// FormatPatchResult returns a human-readable summary of a PatchResult.
func FormatPatchResult(result PatchResult) string {
	var sb strings.Builder

	if len(result.Applied) == 0 && len(result.Skipped) == 0 {
		sb.WriteString("No patch operations provided.\n")
		return sb.String()
	}

	if len(result.Applied) > 0 {
		sb.WriteString(fmt.Sprintf("Applied (%d):\n", len(result.Applied)))
		for _, op := range result.Applied {
			switch op.Action {
			case "set":
				sb.WriteString(fmt.Sprintf("  SET    %s = %s\n", op.Key, op.Value))
			case "delete":
				sb.WriteString(fmt.Sprintf("  DELETE %s\n", op.Key))
			case "rename":
				sb.WriteString(fmt.Sprintf("  RENAME %s -> %s\n", op.Key, op.NewKey))
			}
		}
	}

	if len(result.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped (%d):\n", len(result.Skipped)))
		for _, msg := range result.Errors {
			sb.WriteString(fmt.Sprintf("  ! %s\n", msg))
		}
	}

	return sb.String()
}
