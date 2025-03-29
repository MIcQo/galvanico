package resource

import "testing"

func TestSpecialResource_String(t *testing.T) {
	tests := []struct {
		name string
		b    SpecialResource
		want string
	}{
		{name: "Electricity", b: Electricity, want: "electricity"},
		{name: "waste", b: Waste, want: "waste"},
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
