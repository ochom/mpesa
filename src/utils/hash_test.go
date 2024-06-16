package utils

import "testing"

func TestHashText(t *testing.T) {
	type args struct {
		certPath string
		text     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				certPath: "/Users/mac/projects/ochom/mpesa/certs/b2c_cert.cer",
				text:     "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashText(tt.args.certPath, tt.args.text); len(got) == 0 {
				t.Errorf("HashText() = %v, want len(got) > 0", got)
			}
		})
	}
}
