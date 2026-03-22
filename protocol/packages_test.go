package protocol

import "testing"

func TestRequestHasPackageList(t *testing.T) {
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
			name: "one package",
			req: &Request{
				PackageList: PackageList{
					Packages: []string{"apache"},
				},
			},
			want: true,
		},
		{
			name: "multiple packages",
			req: &Request{
				PackageList: PackageList{
					Packages: []string{"apache", "php8.3", "sqlite3"},
				},
			},
			want: true,
		},
		{
			name: "do upgrade only",
			req: &Request{
				PackageList: PackageList{
					DoUpgrade: true,
				},
			},
			want: true,
		},
		{
			name: "ubuntu pro only",
			req: &Request{
				PackageList: PackageList{
					UbuntuPro: "token",
				},
			},
			want: true,
		},
		{
			name: "all fields set",
			req: &Request{
				PackageList: PackageList{
					Packages:  []string{"apache"},
					DoUpgrade: true,
					UbuntuPro: "token",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasPackageList()
			if got != tt.want {
				t.Fatalf("HasPackageList() = %v, want %v", got, tt.want)
			}
		})
	}
}
