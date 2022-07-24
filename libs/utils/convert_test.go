package utils

import (
	"testing"
)

func TestParseFloat64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"test other type convert to  float64", args{int(64)}, float64(64), false},
		{"test other type convert to  float64", args{"64"}, float64(64), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloat64(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFloat64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	type args struct {
		v interface{}
	}
	// float64 variable
	f := 64.64
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"test other type convert to  int", args{"64"}, int(64), false},
		{"test other type convert to  int", args{f}, int(f), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInt(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	type args struct {
		v interface{}
	}
	// float64 variable
	f := 64.64
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test other type convert to  string", args{64}, "64", false},
		{"test other type convert to  string", args{f}, "64.64", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseString(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// fill float64 convert to timestamp up later
func TestParseTimestamp(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"test other type convert to  timestamp", args{"2022/07/04 15:04:05.999999999"}, int64(1656947045), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimestamp(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTimestamp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
