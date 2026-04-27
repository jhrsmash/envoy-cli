package envfile

import "fmt"

// ChainStep represents a single transformation step in a chain pipeline.
type ChainStep struct {
	Name   string
	Apply  func(env map[string]string) (map[string]string, error)
}

// ChainResult holds the output of a Chain execution.
type ChainResult struct {
	Final map[string]string
	Steps []ChainStepResult
}

// ChainStepResult records the before/after state of a single step.
type ChainStepResult struct {
	Name   string
	Before map[string]string
	After  map[string]string
	Err    error
}

// Chain executes a sequence of transformation steps against an env map,
// passing the output of each step as the input to the next.
// Execution halts at the first step that returns an error.
func Chain(env map[string]string, steps []ChainStep) ChainResult {
	current := copyMap(env)
	results := make([]ChainStepResult, 0, len(steps))

	for _, step := range steps {
		before := copyMap(current)
		next, err := step.Apply(current)
		sr := ChainStepResult{
			Name:   step.Name,
			Before: before,
			After:  next,
			Err:    err,
		}
		results = append(results, sr)
		if err != nil {
			return ChainResult{Final: before, Steps: results}
		}
		current = next
	}

	return ChainResult{Final: current, Steps: results}
}

// FormatChainResult returns a human-readable summary of a chain execution.
func FormatChainResult(r ChainResult) string {
	if len(r.Steps) == 0 {
		return "chain: no steps defined\n"
	}

	out := "chain pipeline:\n"
	for i, sr := range r.Steps {
		status := "ok"
		if sr.Err != nil {
			status = fmt.Sprintf("error: %s", sr.Err)
		}
		delta := len(sr.After) - len(sr.Before)
		sign := "+"
		if delta < 0 {
			sign = ""
		}
		out += fmt.Sprintf("  [%d] %-24s keys: %d (%s%d)  status: %s\n",
			i+1, sr.Name, len(sr.After), sign, delta, status)
		if sr.Err != nil {
			break
		}
	}
	out += fmt.Sprintf("result: %d key(s)\n", len(r.Final))
	return out
}

// copyMap returns a shallow copy of m.
func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
