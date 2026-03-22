package protocol

type MySQL struct {
	Stmts   []string `json:"stmts"`
	DbPath  string   `json:"db_path"`
	DbName  string   `json:"db_name"`
	Comment string   `json:"comment"`
}

type Database struct {
	RedisPort int      `json:"redis_port,omitempty"`
	MySQLs    []*MySQL `json:"mysqls,omitempty"`
}

func (req *Request) HasDatabase() bool {
	if req == nil {
		return false
	}
	return req.RedisPort > 0 || len(req.MySQLs) > 0
}
