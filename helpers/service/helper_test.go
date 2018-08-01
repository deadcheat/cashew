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
			name:    "remove trailing ? for http://example.com?",
			args:    args{u: "http://example.com?"},
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "remove trailing / for http://example.com/",
			args:    args{u: "http://example.com/"},
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "remove trailingslash last '&' in query params for \"http://example.com/?ticket=12345&\"",
			args:    args{u: "http://example.com/?ticket=12345&"},
			want:    "http://example.com?ticket=12345",
			wantErr: false,
		},
		{
			name:    "unexcape query and url for \"http%3A%2F%2Fexample.com%2F%3Fticket%3D12345\"",
			args:    args{u: "http%3A%2F%2Fexample.com%2F%3Fticket%3D12345"},
			want:    "http://example.com?ticket=12345",
			wantErr: false,
		},
		{
			name:    "return escape error for \"%\"",
			args:    args{u: "http://%"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "parse error",
			args:    args{u: "://"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "do not parse when blank",
			args:    args{u: "  "},
			want:    "",
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
