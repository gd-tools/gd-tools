package config

import (
	"github.com/gd-tools/gd-tools/agent"
	"github.com/gd-tools/gd-tools/assets"
)

func (cfg *Config) DeployDatabase() error {
	cfg.Debug("Enter config/database.go")

	req := cfg.NewRequest()

	req.RedisPort = 6379

	path := "mysql/mariadb.conf.d/90-server-port.cnf"
	content, err := assets.Render(path, nil)
	if err != nil {
		return err
	}

	file := agent.File{
		Task:    "write",
		Path:    assets.GetEtcDir(path),
		Content: content,
		Mode:    "0644",
		User:    "root",
		Group:   "root",
		Service: "mariadb",
	}
	req.AddFile(&file)

	if err := req.Send(); err != nil {
		return err
	}

	cfg.Debug("Leave config/database.go")
	return nil
}
