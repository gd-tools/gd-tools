package config

func (cfg *Config) DeployMounts() error {
	cfg.Debug("Enter config/mounts.go")

	req := cfg.NewRequest()
	req.Mounts = cfg.Mounts

	if err := req.Send(); err != nil {
		return err
	}

	cfg.Debug("Leave config/mounts.go")
	return nil
}
