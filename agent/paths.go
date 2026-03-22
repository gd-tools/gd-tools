package agent

import (
	"fmt"
	"net"
	"path/filepath"
)

// These paths are fundamental, can be faked for tests.
func (agt *Agent) RootPath(paths ...string) string {
	if agt != nil {
		if fn := agt.rootPath; fn != nil {
			return fn(path, perm)
		}
	}
	return filepath.Join("/root", paths...)
}

func (agt *Agent) VarPath(paths ...string) string {
	return "TODO"
}

func (agt *Agent) EtcPath(paths ...string) string {
	return "TODO"
}

func (agt *Agent) BinPath(paths ...string) string {
	return "TODO"
}

func (agt *Agent) RunPath(paths ...string) string {
	return "TODO"
}

// These paths are derived from the fundamental paths.
func (agt *Agent) DownloadsPath(paths ...string) string {
	return utils.JoinPath(agt.RootPath("Downloads"), paths...)
}

func (agt *Agent) ToolsPath(paths ...string) string {
	return utils.JoinPath(agt.VarPath("gd-tools"), paths...)
}

func (agt *Agent) DataPath(paths ...string) string {
	return utils.JoinPath(agt.ToolsPath("data"), paths...)
}

func (agt *Agent) CertsPath(paths ...string) string {
	return utils.JoinPath(agt.DataPath("certs"), paths...)
}

func (agt *Agent) ToolsApachePath(paths ...string) string {
	return utils.JoinPath(agt.DataPath("apache"), paths...)
}

func (agt *Agent) LogsPath(paths ...string) string {
	return utils.JoinPath(agt.ToolsPath("logs"), paths...)
}

func (agt *Agent) EtcApachePath(paths ...string) string {
	return utils.JoinPath(agt.EtcPath("apache2"), paths...)
}

func (agt *Agent) EtcPhpPath(paths ...string) string {
	if agt.Baseline == nil || agt.Baseline.PHP == "" {
		panic("missing Baseline or PHP version")
	}
	return utils.JoinPath(agt.EtcPath("php", agt.Baseline.PHP), paths...)
}
