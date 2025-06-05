package helper

import "testing"

func TestParseMode(t *testing.T) {
	tests := []struct {
		input       string
		expected    Mode
		expectError bool
	}{
		{"development", Development, false},
		{"production", Production, false},
		{"staging", 0, true},
		{"", 0, true},
		{"DEVELOPMENT", 0, true}, // case-sensitive test
	}

	for _, test := range tests {
		result, err := ParseMode(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("expected error for input %q, got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", test.input, err)
			} else if result != test.expected {
				t.Errorf("for input %q, expected mode %v, got %v", test.input, test.expected, result)
			}
		}
	}
}
