package protocol

import "testing"

func TestRequestAddMySQL(t *testing.T) {
	req := &Request{}
	sql := &MySQL{
		DbPath: "/run/mysqld/mysqld.sock",
		DbName: "nextcloud",
		Stmts: []string{
			"CREATE DATABASE nextcloud;",
		},
	}

	req.AddMySQL(sql)

	if len(req.MySQLs) != 1 {
		t.Fatalf("expected 1 MySQL entry, got %d", len(req.MySQLs))
	}
	if req.MySQLs[0] != sql {
		t.Fatalf("expected appended MySQL pointer to match input")
	}
}

func TestRequestAddMySQLIgnoresNilReceiver(t *testing.T) {
	var req *Request
	sql := &MySQL{DbName: "test"}

	req.AddMySQL(sql)
}

func TestRequestAddMySQLIgnoresNilSQL(t *testing.T) {
	req := &Request{}

	req.AddMySQL(nil)

	if len(req.MySQLs) != 0 {
		t.Fatalf("expected no MySQL entries, got %d", len(req.MySQLs))
	}
}

func TestRequestHasDatabase(t *testing.T) {
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
			name: "redis only",
			req: &Request{
				Database: Database{
					RedisPort: 6379,
				},
			},
			want: true,
		},
		{
			name: "mysql only",
			req: &Request{
				Database: Database{
					MySQLs: []*MySQL{
						{DbName: "app"},
					},
				},
			},
			want: true,
		},
		{
			name: "both",
			req: &Request{
				Database: Database{
					RedisPort: 6379,
					MySQLs: []*MySQL{
						{DbName: "app"},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.HasDatabase()
			if got != tt.want {
				t.Fatalf("HasDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}
