package getopt_test

import (
	"reflect"
	"testing"

	getoptutil "github.com/lucasepe/resto/internal/util/getopt"
	"github.com/lucasepe/x/getopt"
)

func TestHasOpt(t *testing.T) {
	tests := []struct {
		name     string
		opts     []getopt.OptArg
		lookup   []string
		expected bool
	}{
		{
			name: "Option found",
			opts: []getopt.OptArg{
				{Option: "-a"},
				{Option: "-b"},
			},
			lookup:   []string{"-b"},
			expected: true,
		},
		{
			name: "Option not found, fallback used",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{"-x"},
			expected: false,
		},
		{
			name: "Multiple matches, first one returned",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "first"},
				{Option: "-b", Argument: "second"},
				{Option: "-c", Argument: "third"},
			},
			lookup:   []string{"-b", "-c"},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getoptutil.HasOpt(tc.opts, tc.lookup)
			if result != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestFindOptVal(t *testing.T) {
	tests := []struct {
		name     string
		opts     []getopt.OptArg
		lookup   []string
		expected string
	}{
		{
			name: "Option found",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
				{Option: "-b", Argument: "valueB"},
			},
			lookup:   []string{"-b"},
			expected: "valueB",
		},
		{
			name: "Option not found",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{"-x"},
			expected: "",
		},
		{
			name:     "Empty opts list",
			opts:     []getopt.OptArg{},
			lookup:   []string{"-a"},
			expected: "",
		},
		{
			name: "Multiple matches, first one returned",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "first"},
				{Option: "-b", Argument: "second"},
				{Option: "-c", Argument: "third"},
			},
			lookup:   []string{"-b", "-c"},
			expected: "second",
		},
		{
			name: "Lookup list empty",
			opts: []getopt.OptArg{
				{Option: "-a", Argument: "valueA"},
			},
			lookup:   []string{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getoptutil.OptVal(tc.opts, tc.lookup)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestWantsHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"no args", []string{}, true},
		{"exact match", []string{"help"}, true},
		{"uppercase", []string{"HELP"}, true},
		{"mixed case", []string{"HeLp"}, true},
		{"not help", []string{"start"}, false},
		{"help not first", []string{"start", "help"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getoptutil.WantsHelp(tt.args)
			if result != tt.expected {
				t.Errorf("WantsHelp(%v) = %v; want %v", tt.args, result, tt.expected)
			}
		})
	}
}

func TestAllOptArgs(t *testing.T) {
	tests := []struct {
		name     string
		opts     []getopt.OptArg
		keys     []string
		expected []string
	}{
		{
			name:     "No matching options",
			opts:     []getopt.OptArg{{Option: "-a", Argument: "foo"}},
			keys:     []string{"-w", "--when"},
			expected: []string{},
		},
		{
			name: "Single match",
			opts: []getopt.OptArg{
				{Option: "-w", Argument: "now"},
			},
			keys:     []string{"-w", "--when"},
			expected: []string{"now"},
		},
		{
			name: "Multiple matches mixed short/long",
			opts: []getopt.OptArg{
				{Option: "-w", Argument: "now"},
				{Option: "--when", Argument: "later"},
				{Option: "-o", Argument: "something"},
			},
			keys:     []string{"-w", "--when"},
			expected: []string{"now", "later"},
		},
		{
			name: "Duplicate options",
			opts: []getopt.OptArg{
				{Option: "-w", Argument: "1"},
				{Option: "-w", Argument: "2"},
				{Option: "--when", Argument: "3"},
			},
			keys:     []string{"-w", "--when"},
			expected: []string{"1", "2", "3"},
		},
		{
			name: "Empty argument values",
			opts: []getopt.OptArg{
				{Option: "-w", Argument: ""},
				{Option: "--when", Argument: "x"},
			},
			keys:     []string{"-w", "--when"},
			expected: []string{"", "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getoptutil.AllOptArgs(tt.opts, tt.keys)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AllOptArgs() = %v, want %v", got, tt.expected)
			}
		})
	}
}
