package unit

import "testing"

func TestShip_String(t *testing.T) {
	tests := []struct {
		name string
		s    Ship
		want string
	}{
		{name: "gunboat", s: Gunboat, want: "gunboat"},
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

func TestShipFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Ship
		wantErr bool
	}{
		{name: "gunboat", args: args{s: "gunboat"}, want: Gunboat, wantErr: false},
		{name: "unknown", args: args{s: "blabla"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ShipFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShipFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ShipFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
