package resource

import "testing"

func TestUserResource_String(t *testing.T) {
	tests := []struct {
		name string
		u    UserResource
		want string
	}{
		{name: "Gold", u: Gold, want: "gold"},
		{name: "Ship", u: Ship, want: "ship"},
		{name: "Unknown", u: 99, want: unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
