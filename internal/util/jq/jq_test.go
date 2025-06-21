package jq

import (
	"testing"
)

var podJSON = []byte(`{
  "status": {
    "phase": "Running",
    "conditions": [
      {"type": "Initialized", "status": "True"},
      {"type": "Ready", "status": "True"},
      {"type": "ContainersReady", "status": "True"},
      {"type": "PodScheduled", "status": "True"}
    ],
    "containerStatuses": [
      {"name": "test-container", "ready": true, "started": true}
    ]
  }
}`)

func TestEvaluateBool(t *testing.T) {
	tests := []struct {
		name      string
		jsonInput []byte
		expr      string
		want      bool
		wantErr   bool
	}{
		{
			name:      "Valid condition Ready == True",
			jsonInput: podJSON,
			expr:      `any(.status.conditions[]; .type == "Ready" and .status == "True")`,
			want:      true,
			wantErr:   false,
		},
		{
			name:      "False condition Ready == False",
			jsonInput: podJSON,
			expr:      `any(.status.conditions[]; .type == "Ready" and .status == "False")`,
			want:      false,
			wantErr:   false,
		},
		{
			name:      "Container ready == true",
			jsonInput: podJSON,
			expr:      `.status.containerStatuses[0].ready`,
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Wrong type (string returned)",
			jsonInput: podJSON,
			expr:      `.status.phase`,
			want:      false,
			wantErr:   true, // not a boolean
		},
		{
			name:      "JQ syntax error",
			jsonInput: podJSON,
			expr:      `.status.conditions[`,
			want:      false,
			wantErr:   true,
		},
		{
			name:      "Invalid JSON",
			jsonInput: []byte(`{status:}`),
			expr:      `.status.phase == "Running"`,
			want:      false,
			wantErr:   true,
		},
		{
			name:      "Missing field still works",
			jsonInput: podJSON,
			expr:      `.status.nonexistentField == null`,
			want:      true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalBoolExpr(tt.jsonInput, tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateBool() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateBool() = %v, want = %v", got, tt.want)
			}
		})
	}
}
