package building

import (
	"reflect"
	"testing"
)

func TestFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Building
		wantErr bool
	}{
		{name: "cathedral", args: args{s: "cathedral"}, want: Cathedral, wantErr: false},
		{name: "city_hall", args: args{s: "city_hall"}, want: CityHall, wantErr: false},
		{name: "random building", args: args{s: "blablalba"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilding_String(t *testing.T) {
	tests := []struct {
		name string
		b    Building
		want string
	}{
		{name: "cathedral", b: Cathedral, want: "cathedral"},
		{name: "city_hall", b: CityHall, want: "city_hall"},
		{name: "unknown", b: Building(99), want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilding_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		b       Building
		want    []byte
		wantErr bool
	}{
		{name: "exists", b: CityHall, want: []byte(`"city_hall"`), wantErr: false},
		{name: "not exists", b: 9999, want: []byte(`"unknown"`), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilding_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       Building
		args    args
		wantErr bool
	}{
		{name: "simple", b: 0, args: args{data: []byte(`"city_hall"`)}, wantErr: false},
		{name: "simple", b: 0, args: args{data: []byte(`"cathedral"`)}, wantErr: false},
		{name: "simple", b: 0, args: args{data: []byte(`"hello"`)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
