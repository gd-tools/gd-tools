package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	DHParamsFile = "dhparams.pem"
	DHTimeout    = 5 * time.Minute
	DHMinBits    = 2048
)

// GenerateDHParams returns existing dhparams.pem or generates a new one using openssl.
func GenerateDHParams(bits int) ([]byte, error) {
	dhBytes, err := os.ReadFile(DHParamsFile)
	if err == nil && len(dhBytes) > 0 {
		return dhBytes, nil
	}

	if bits < DHMinBits {
		bits = DHMinBits
	}

	ctx, cancel := context.WithTimeout(context.Background(), DHTimeout)
	defer cancel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, "openssl", "dhparam", "-outform", "PEM", fmt.Sprintf("%d", bits))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("openssl dhparam timed out")
		}
		return nil, fmt.Errorf("failed to run openssl dhparam: %w\n%s", err, stderr.String())
	}

	dhBytes = stdout.Bytes()

	tmp := DHParamsFile + ".tmp"

	if err := os.WriteFile(tmp, dhBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write %s: %w", tmp, err)
	}

	if err := os.Rename(tmp, DHParamsFile); err != nil {
		return nil, fmt.Errorf("failed to replace %s: %w", DHParamsFile, err)
	}

	return dhBytes, nil
}
