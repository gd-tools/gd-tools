package protocol

import "testing"

func TestBootstrapHasBootstrap(t *testing.T) {
	tests := []struct {
		name string
		in   Bootstrap
		want bool
	}{
		{
			name: "empty bootstrap",
			in:   Bootstrap{},
			want: false,
		},
		{
			name: "fqdn only",
			in: Bootstrap{
				FQDN: "host.example.org",
			},
			want: true,
		},
		{
			name: "time zone only",
			in: Bootstrap{
				TimeZone: "Europe/Berlin",
			},
			want: true,
		},
		{
			name: "ssh port only",
			in: Bootstrap{
				SSHPort: 2222,
			},
			want: true,
		},
		{
			name: "all fields set",
			in: Bootstrap{
				FQDN:     "host.example.org",
				TimeZone: "Europe/Berlin",
				SSHPort:  2222,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.HasBootstrap()
			if got != tt.want {
				t.Fatalf("HasBootstrap() = %v, want %v", got, tt.want)
			}
		})
	}
}
