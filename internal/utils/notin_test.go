package utils

import (
	"reflect"
	"testing"
)

func TestFilterNotIn(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected []string
	}{
		{
			name:     "Basic case",
			a:        []string{"apple", "banana", "cherry"},
			b:        []string{"banana"},
			expected: []string{"apple", "cherry"},
		},
		{
			name:     "No overlap",
			a:        []string{"a", "b", "c"},
			b:        []string{"x", "y", "z"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "All overlap",
			a:        []string{"a", "b"},
			b:        []string{"a", "b"},
			expected: []string{},
		},
		{
			name:     "Empty input A",
			a:        []string{},
			b:        []string{"a", "b"},
			expected: []string{},
		},
		{
			name:     "Empty input B",
			a:        []string{"a", "b"},
			b:        []string{},
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterNotIn(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("FilterNotIn() = %v, want %v", got, tt.expected)
			}
		})
	}
}
