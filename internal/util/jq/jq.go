package jq

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

// EvalBoolExpr evaluates a JQ expression against the given JSON input and returns a boolean result.
//
// inputJSON must be a valid JSON byte slice.
//
// jqExpr is the JQ expression to evaluate, which **must** return a boolean value (true or false).
// If the expression returns a value of a different type, the function returns an error.
//
// The function executes the expression on the entire JSON and returns the first boolean result produced.
//
// If evaluation fails or if no boolean value is produced by the expression, an error is returned.
//
// Example of a valid expression:
//
//	.status.conditions[] | select(.type == "Ready") | .status == "True"
//
// Returns:
//   - (bool, nil) if the expression returns true or false successfully
//   - (false, error) if there is a JSON parsing error, JQ parsing error, or if the expression does not return a boolean
func EvalBoolExpr(inputJSON []byte, jqExpr string) (bool, error) {
	var data any
	if err := json.Unmarshal(inputJSON, &data); err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}

	query, err := gojq.Parse(jqExpr)
	if err != nil {
		return false, fmt.Errorf("invalid JQ expression: %w", err)
	}

	iter := query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			return false, fmt.Errorf("evaluation error: %w", err)
		}
		if b, ok := v.(bool); ok {
			return b, nil
		} else {
			// Fail if result is not boolean
			return false, fmt.Errorf("expression did not return a boolean: got %T (%v)", v, v)
		}
	}

	return false, nil
}
