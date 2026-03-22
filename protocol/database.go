package protocol

// MySQL describes SQL statements to be executed against a MySQL database.
//
// DbPath specifies the MySQL access point, typically a Unix socket path.
// DbName specifies the target database.
// Stmts contains SQL statements executed in order.
// Comment is optional and used for logging/debugging purposes.
type MySQL struct {
	Stmts   []string `json:"stmts"`
	DbPath  string   `json:"db_path"`
	DbName  string   `json:"db_name"`
	Comment string   `json:"comment,omitempty"`
}

// Database contains database-related setup instructions.
//
// RedisPort defines the port for a Redis instance, if required.
// MySQLs contains MySQL execution blocks.
type Database struct {
	RedisPort int      `json:"redis_port,omitempty"`
	MySQLs    []*MySQL `json:"mysqls,omitempty"`
}

// AddMySQL adds a MySQL execution block to the request.
func (req *Request) AddMySQL(sql *MySQL) {
	if req == nil || sql == nil {
		return
	}
	req.MySQLs = append(req.MySQLs, sql)
}

// HasDatabase reports whether the request contains database-related instructions.
func (req *Request) HasDatabase() bool {
	if req == nil {
		return false
	}
	return req.RedisPort > 0 || len(req.MySQLs) > 0
}
