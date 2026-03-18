package utils

import (
	"os"
	"os/exec"
)

func RSAKeyPair(fqdn string) ([]byte, []byte, error) {
	priv, err1 := os.ReadFile("root_id_rsa")
	publ, err2 := os.ReadFile("root_id_rsa.pub")
	if err1 == nil && err2 == nil {
		return priv, publ, nil
	}

	cmd := exec.Command(
		"ssh-keygen",
		"-t", "rsa",
		"-C", "root@"+fqdn,
		"-f", "root_id_rsa",
		"-N", "",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}

	priv, err1 = os.ReadFile("root_id_rsa")
	if err1 != nil {
		return nil, nil, err1
	}
	publ, err2 = os.ReadFile("root_id_rsa.pub")
	if err2 != nil {
		return nil, nil, err2
	}

	return priv, publ, nil
}
