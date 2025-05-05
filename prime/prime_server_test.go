package prime

import "testing"

func Test_isPrime(t *testing.T) {
	type args struct {
		n float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1 is not prime",
			args: args{n: 1},
			want: false,
		},
		{
			name: "2 is prime",
			args: args{n: 2},
			want: true,
		},
		{
			name: "3 is prime",
			args: args{n: 3},
			want: true,
		},
		{
			name: "4 is not prime",
			args: args{n: 4},
			want: false,
		},
		{
			name: "5 is prime",
			args: args{n: 5},
			want: true,
		},
		{
			name: "6 is not prime",
			args: args{n: 6},
			want: false,
		},
		{
			name: "7 is prime",
			args: args{n: 7},
			want: true,
		},
		{
			name: "11 is prime",
			args: args{n: 11},
			want: true,
		},
		{
			name: "non-integer is not prime",
			args: args{n: 2.5},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrime(tt.args.n); got != tt.want {
				t.Errorf("isPrime() = %v, want %v", got, tt.want)
			}
		})
	}
}
