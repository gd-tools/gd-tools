package agent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"github.com/gd-tools/gd-tools/assets"
)

type WordPress struct {
	Instance  string
	Directory string
}

func (wp *WordPress) BaseDir() string {
	return assets.GetToolsDir("data", "wordpress", wp.Instance, wp.Directory)
}

func (wp *WordPress) LogsDir() string {
	return assets.GetToolsDir("logs", "wordpress", wp.Instance)
}

func (wp *WordPress) RunWPCLI(resp *Response, cmd string, args ...string) (string, error) {
	usr, err := user.Lookup("www-data")
	if err != nil {
		return "", fmt.Errorf("user lookup for www-data failed: %w", err)
	}
	uid, _ := strconv.Atoi(usr.Uid)
	gid, _ := strconv.Atoi(usr.Gid)

	wpPath := assets.GetBinDir("wp-cli")
	wpArgs := []string{cmd}
	wpArgs = append(wpArgs, args...)
	wpArgs = append(wpArgs, "--path="+wp.BaseDir())

	command := exec.Command(wpPath, wpArgs...)
	command.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	command.Env = append(os.Environ(),
		"HOME="+wp.BaseDir(),
		"PATH=/usr/local/bin:/usr/bin:/bin",
	)

	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("RunWPCLI failed: %w\n%s", err, output.String())
	}

	return output.String(), nil
}
