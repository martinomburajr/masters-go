package evolution

import "testing"

func TestRandString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandString(tt.args.n); got != tt.want {
				t.Errorf("RandString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compString(str string, str2 string) bool {
	if str == str2 {
		return true
	}
	return false
}

func compInt(str int, str2 int) bool {
	if str == str2 {
		return true
	}
	return false
}

var response = false
func BenchmarkCompareString(b *testing.B) {
	slices := []string{"1234", "2345", "3456", "5678"}
	for i := 0; i < b.N; i++ {
		for e := range slices {
			response = compString(slices[e], "1234")
		}
	}
}


func BenchmarkCompareInt(b *testing.B) {
	slices := []int{1234, 2345, 3456, 5678}
	for i := 0; i < b.N; i++ {
		for e := range slices {
			response = compInt(slices[e], 1234)
		}
	}
}

func BenchmarkRandString(b *testing.B) {
	n := 5;
	for i := 0; i < b.N; i++ {
		RandString(n)
	}
}