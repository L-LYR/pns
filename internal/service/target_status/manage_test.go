package target_status

import "testing"

func Test_StatusOffset(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "test",
			args: args{
				key: "1234:1",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := _StatusOffset(tt.args.key); got != tt.want {
				t.Errorf("_StatusOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}
