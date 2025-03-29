package resource

import "testing"

func TestBaseResource_String(t *testing.T) {
	tests := []struct {
		name string
		b    BaseResource
		want string
	}{
		{name: "Wood", b: Wood, want: "wood"},
		{name: "Water", b: Water, want: "water"},
		{name: "unknown", b: 99, want: unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
