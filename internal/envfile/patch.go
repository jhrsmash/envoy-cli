package envfile

import "fmt"

// PatchOp represents a single patch operation on an env map.
type PatchOp struct {
	Action string // "set", "delete", "rename"
	Key    string
	Value  string // used by "set"
	NewKey string // used by "rename"
}

// PatchResult holds the outcome of applying a patch.
type PatchResult struct {
	Applied []PatchOp
	Skipped []PatchOp
	Errors  []string
}

// Patch applies a sequence of PatchOps to a copy of the given env map.
// It never mutates the input map.
func Patch(env map[string]string, ops []PatchOp) (map[string]string, PatchResult) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var result PatchResult

	for _, op := range ops {
		switch op.Action {
		case "set":
			out[op.Key] = op.Value
			result.Applied = append(result.Applied, op)

		case "delete":
			if _, exists := out[op.Key]; !exists {
				result.Skipped = append(result.Skipped, op)
				result.Errors = append(result.Errors, fmt.Sprintf("delete: key %q not found", op.Key))
				continue
			}
			delete(out, op.Key)
			result.Applied = append(result.Applied, op)

		case "rename":
			val, exists := out[op.Key]
			if !exists {
				result.Skipped = append(result.Skipped, op)
				result.Errors = append(result.Errors, fmt.Sprintf("rename: key %q not found", op.Key))
				continue
			}
			delete(out, op.Key)
			out[op.NewKey] = val
			result.Applied = append(result.Applied, op)

		default:
			result.Skipped = append(result.Skipped, op)
			result.Errors = append(result.Errors, fmt.Sprintf("unknown action %q for key %q", op.Action, op.Key))
		}
	}

	return out, result
}
