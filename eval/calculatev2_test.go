package eval

import "testing"

func TestCalculateV2(t *testing.T) {
	type args struct {
		mathematicalExpression string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"", args{""}, 0, true},
		{"1", args{"1"}, 1, false},
		{"1+1", args{"1+1"}, 2, false},
		{"-1+1", args{"-1+1"}, 0, false},
		{"-1+1+1", args{"-1+1+1"}, 1, false},
		{"-1+1+1-1", args{"-1+1+1-1"}, 0, false},
		{"-100+1", args{"-100+1"}, -99, false},
		{"-100*1", args{"-100*1"}, -100, false},
		{"-100/1", args{"-100/1"}, -100, false},
		{"-100*1/2", args{"-100*1/2"}, -50, false},
		{"-100*1/2", args{"-100*1/2"}, -50, false},
		{"2.5*10", args{"2.5*10"}, 25, false},
		{"2.5*10*10", args{"2.5*10*10"}, 250, false},
		{"3/0.5", args{"3/0.5"}, 6, false},
		{"0*-1", args{"0*-1"}, -0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateV2(tt.args.mathematicalExpression)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateV2() = %v, want %v", got, tt.want)
			}
		})
	}
}
