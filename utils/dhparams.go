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

// DHParams returns existing dhparams.pem or generates a new one using openssl.
func DHParams(bits int) ([]byte, error) {
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

	if err := SaveFile(DHParamsFile, dhBytes); err != nil {
		return nil, err
	}

	return dhBytes, nil
}
