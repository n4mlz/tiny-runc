package lib

import (
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func NewState(containerID string, bundlePath string) (specs.State, error) {
	absolutePath, err := filepath.Abs(bundlePath)
	if err != nil {
		return specs.State{}, err
	}

	return specs.State{
		// TODO: add Version field
		ID:     containerID,
		Status: "creating",
		Pid:    -1,
		Bundle: absolutePath,
		// TODO: add Annotations field
	}, nil
}
