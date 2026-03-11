package config

func (cfg *Config) DeployBootstrap() error {
	cfg.Debug("Enter config/bootstrap.go")

	req := cfg.NewRequest()
	req.FQDN = cfg.FQDN()
	req.TimeZone = cfg.TimeZone
	req.SwapSize = cfg.SwapSize

	if err := req.Send(); err != nil {
		return err
	}

	cfg.Debug("Leave config/bootstrap.go")
	return nil
}
