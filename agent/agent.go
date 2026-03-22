package agent

import (
	"github.com/spf13/afero"
)

type Agent struct {
	fs afero.Fs
}

func NewAgent(fs afero.Fs) *Agent {
	return &Agent{fs: fs}
}
