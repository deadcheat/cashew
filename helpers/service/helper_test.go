package service

import "testing"

func TestNormalizeURL(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "http://google.com?ticket=12345",
			args:    args{u: "http://google.com?ticket=12345"},
			want:    "http://google.com?ticket=12345",
			wantErr: false,
		},
		{
			name:    "http://google.com/?ticket=12345&",
			args:    args{u: "http://google.com/?ticket=12345&"},
			want:    "http://google.com?ticket=12345",
			wantErr: false,
		},
		{
			name:    "http%3A%2F%2Fgoogle.com%2F%3Fticket%3D12345",
			args:    args{u: "http%3A%2F%2Fgoogle.com%2F%3Fticket%3D12345"},
			want:    "http://google.com?ticket=12345",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
