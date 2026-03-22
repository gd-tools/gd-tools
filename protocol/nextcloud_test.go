package protocol

import "testing"

func TestNextcloudFQDN(t *testing.T) {
	tests := []struct {
		name string
		nc   *Nextcloud
		want string
	}{
		{
			name: "nil",
			nc:   nil,
			want: "",
		},
		{
			name: "host and domain",
			nc: &Nextcloud{
				HostName:   "cloud",
				DomainName: "example.org",
			},
			want: "cloud.example.org",
		},
		{
			name: "domain only",
			nc: &Nextcloud{
				DomainName: "example.org",
			},
			want: "example.org",
		},
		{
			name: "host only",
			nc: &Nextcloud{
				HostName: "cloud",
			},
			want: "cloud",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.nc.FQDN()
			if got != tt.want {
				t.Fatalf("FQDN() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNextcloudSlashSubdir(t *testing.T) {
	tests := []struct {
		name string
		nc   *Nextcloud
		want string
	}{
		{
			name: "nil",
			nc:   nil,
			want: "",
		},
		{
			name: "empty",
			nc: &Nextcloud{
				Subdir: "",
			},
			want: "",
		},
		{
			name: "slash only",
			nc: &Nextcloud{
				Subdir: "/",
			},
			want: "",
		},
		{
			name: "plain",
			nc: &Nextcloud{
				Subdir: "nextcloud",
			},
			want: "/nextcloud",
		},
		{
			name: "with slashes",
			nc: &Nextcloud{
				Subdir: "/nextcloud/",
			},
			want: "/nextcloud",
		},
		{
			name: "nested",
			nc: &Nextcloud{
				Subdir: "apps/nextcloud",
			},
			want: "/apps/nextcloud",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.nc.SlashSubdir()
			if got != tt.want {
				t.Fatalf("SlashSubdir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequestHasNextcloudApp(t *testing.T) {
	tests := []struct {
		name string
		req  *Request
		want bool
	}{
		{
			name: "nil request",
			req:  nil,
			want: false,
		},
		{
			name: "empty request",
			req:  &Request{},
			want: false,
		},
		{
			name: "nextcloud only",
			req: &Request{
				NextcloudApp: NextcloudApp{
					Nextcloud: &Nextcloud{Name: "main"},
				},
			},
			want: true,
		},
		{
			name: "configs only",
			req: &Request{
				NextcloudApp: NextcloudApp{
					NextConfigs: []*NextConfig{
						{Key: "overwrite.cli.url", Type: "string", Value: "https://cloud.example.org"},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasNextcloudApp()
			if got != tt.want {
				t.Fatalf("HasNextcloudApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
