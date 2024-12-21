package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestParseDowMatcher(t *testing.T) {
	tests := []struct {
		spec     string
		time     string
		expected bool
	}{
		{"*", "Thu Jan 2 2025", true},
		{"*", "Fri Jan 3 2025", true},
		{"?", "Thu Jan 2 2025", true},
		{"?", "Fri Jan 3 2025", true},

		{"5", "Fri Jan 3 2025", true},
		{"2", "Wed Jan 8 2025", false},
		{"6", "Sat Jan 4 2025", true},
		{"Sat", "Sat Jan 4 2025", true},
		{"SAT", "Sat Jan 4 2025", true},

		{"1-3", "Mon Jan 6 2025", true},
		{"Mon-Wed", "Mon Jan 6 2025", true},
		{"1-3", "Sun Jan 5 2025", false},
		{"1-3", "Tue Jan 7 2025", true},
		{"0-2", "Wed Jan 8 2025", false},

		{"1/2", "Mon Jan 6 2025", true},
		{"1/2", "Wed Jan 8 2025", true},
		{"1/2", "Sun Jan 5 2025", false},
		{"*/2", "Sun Jan 5 2025", true},
		{"*/2", "Mon Jan 6 2025", false},

		{"1-4/3", "Mon Jan 6 2025", true},
		{"1-4/3", "Thu Jan 9 2025", true},
		{"1-4/3", "Fri Jan 10 2025", false},

		{"1-4/3,5", "Fri Jan 10 2025", true},
		{"1-4/3,6", "Fri Jan 10 2025", false},

		{"1,2,3", "Mon Jan 6 2025", true},
		{"1,2,3", "Thu Jan 9 2025", false},
		{"1,02,3", "Tue Jan 7 2025", true},
	}

	const layout = "Mon Jan 2 2006"
	for _, test := range tests {
		t.Run(fmt.Sprintf("spec=%s,now=%s", strings.Replace(test.spec, "/", "|", -1), test.time), func(t *testing.T) {
			time, err := time.Parse(layout, test.time)
			if err != nil {
				t.Fatalf("expected nil error, got %s", err)
			}

			matcher, err := ParseDow(test.spec)
			if err != nil {
				t.Fatalf("expected nil error, got %s", err)
			}

			actual := matcher(time)
			if actual != test.expected {
				t.Fatalf("spec=%s, time=%s, expected=%t, got=%t",
					test.spec, test.time, test.expected, actual)
			}
		})
	}
}

func TestParseDowErrors(t *testing.T) {
	tests := []struct {
		spec     string
		expected string
	}{
		{"#", "failed to parse"},
		{"*-5", "invalid expression"},
		{"*-5/2", "invalid expression"},
		{"5-*", "failed to parse"},
		{"5-*/22", "failed to parse"},
		{"-1", "failed to parse"},
		{"-22", "failed to parse"},
		{"8", "value 8 out of valid range [0, 6]"},
		{"4-7", "value 7 out of valid range [0, 6]"},
		{"4-7/36", "value 7 out of valid range [0, 6]"},
		{"4-5/ABC", "failed to parse"},
		{"*//2", "invalid expression"},
		{"*/2/", "invalid expression"},
		{"*-2-", "invalid expression"},
		{"2-", "failed to parse"},
		{"1-2-", "invalid expression"},
	}

	for _, test := range tests {
		t.Run(strings.Replace(test.spec, "/", "|", -1), func(t *testing.T) {
			_, err := ParseDow(test.spec)
			if err == nil {
				t.Fatal("expected non-nil error, got nil")
			}

			actual := err.Error()
			if !strings.Contains(actual, test.expected) {
				t.Fatalf("spec=%s, expectedError=%s, gotError=%s",
					test.spec, test.expected, actual)
			}
		})
	}
}