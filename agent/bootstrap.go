package agent

import (
	"fmt"
	"os"
	"strings"
)

const (
	FirstRunMarker = ".gd-tools-first-run"
)

func (agt *Agent) BootstrapTest(req *Request) bool {
	path := agt.RootDir(FirstRunMarker)
	if _, err := agt.ReadFile(path); os.IsNotExist(err) {
		return true
	}
	if req == nil {
		return false
	}
	return req.FQDN != "" || req.TimeZone != "" || req.SwapSize != ""
}

func (agt *Agent) BootstrapHandler(req *Request, resp *Response) error {
	if req == nil || resp == nil {
		return nil
	}
	firstRun := false
	path := agt.RootDir(FirstRunMarker)
	if _, err := agt.ReadFile(path); os.IsNotExist(err) {
		firstRun = true
		cmds := []string{
			"apt-get update",
			"apt-get upgrade -y",
		}
		if _, err := agt.RunShell(cmds); err != nil {
			resp.Err = fmt.Sprintf("failed updating on first run: %v", err)
			return err
		}
		if err := agt.Create(path); err != nil {
			resp.Err = fmt.Sprintf("failed creating %s: %v", path, err)
			return err
		}
		resp.Info("FirstRunMarker was created")
	} else if req.FQDN != "" {
		resp.Info("✅ FirstRunMarker exists")
	}

	if req.FQDN != "" {
		if status, err := bootstrapHostName(req.FQDN); err != nil {
			resp.Err = fmt.Sprintf("failed setting hostname to %s: %v", req.FQDN, err)
			return err
		} else {
			resp.Info(status)
		}
	}

	if req.TimeZone != "" {
		if status, err := bootstrapTimeZone(req.TimeZone); err != nil {
			resp.Err = fmt.Sprintf("failed setting timezone to %s: %v", req.TimeZone, err)
			return err
		} else {
			resp.Info(status)
		}
	}

	if req.SwapSize != "" && req.SwapSize != "0" {
		if status, err := bootstrapSwapSize(req.SwapSize); err != nil {
			resp.Err = fmt.Sprintf("failed setting swapsize to %s: %v", req.SwapSize, err)
			return err
		} else {
			resp.Info(status)
		}
	}

	if firstRun {
		resp.Info("New System ----- please check reboot ...")
	}
	return nil
}

func bootstrapHostName(fqdn string) (string, error) {
	currName, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to read hostname: %w", err)
	}
	if currName == fqdn {
		return fmt.Sprintf("✅ hostname %s is set", fqdn), nil
	}
	if _, err := RunCommand("hostnamectl", "set-hostname", fqdn); err != nil {
		return "", fmt.Errorf("failed to set hostname: %w", err)
	}
	return fmt.Sprintf("hostname updated to: %s", fqdn), nil
}

func bootstrapTimeZone(timeZone string) (string, error) {
	link, err := os.Readlink("/etc/localtime")
	if err == nil {
		if strings.Contains(link, timeZone) {
			return fmt.Sprintf("✅ timezone %s is set", timeZone), nil
		}
	}

	if _, err := RunCommand("timedatectl", "set-timezone", timeZone); err != nil {
		return "", fmt.Errorf("timedatectl failed: %w", err)
	}

	return fmt.Sprintf("timezone updated to: %s", timeZone), nil
}

func bootstrapSwapSize(swapSize string) (string, error) {
	swapFile := "/swap.img"

	if _, err := RunCommand("swapon", "--show"); err == nil {
		out, _ := RunCommand("swapon", "--show")
		if strings.Contains(string(out), swapFile) {
			return fmt.Sprintf("✅ swapfile '%s' active\n", swapFile), nil
		}
	}

	if _, err := os.Stat(swapFile); err == nil {
		return fmt.Sprintf("✅ swapfile '%s' exists\n", swapFile), nil
	}

	genCmds := []string{
		fmt.Sprintf("fallocate -l %s %s", swapSize, swapFile),
		fmt.Sprintf("chmod 600 %s", swapFile),
		fmt.Sprintf("mkswap %s", swapFile),
		fmt.Sprintf("swapon %s", swapFile),
	}
	if _, err := RunShell(genCmds); err != nil {
		return "", fmt.Errorf("swapfile generation failed: %w", err)
	}

	tabCmds := []string{
		`sed -i -e '/swap/d' /etc/fstab`,
		`echo '/swap.img none swap sw 0 0' >> /etc/fstab`,
	}
	if _, err := RunShell(tabCmds); err != nil {
		return "", fmt.Errorf("failed to include swapfile in /etc/fstab: %w", err)
	}

	return fmt.Sprintf("swapfile '%s' with %s exists\n", swapFile, swapSize), nil
}
