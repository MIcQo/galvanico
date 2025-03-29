package unit

import "testing"

func TestUnit_String(t *testing.T) {
	tests := []struct {
		name string
		s    Unit
		want string
	}{
		{name: "rifleman", s: Rifleman, want: "rifleman"},
		{name: "unknown", s: 99, want: unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Unit
		wantErr bool
	}{
		{name: "rifleman", args: args{s: "rifleman"}, want: Rifleman, wantErr: false},
		{name: "unknown", args: args{s: "blabla"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnitFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnitFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnitFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
