package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
)

type Container struct {
	State specs.State
	Root  string
}

func NewContainer(containerID, bundlePath string) (Container, error) {
	root := filepath.Join(BASE_PATH, CONTAINER_PATH, containerID)
	if _, err := os.Stat(root); err == nil {
		return Container{}, fmt.Errorf("container %s already exists", containerID)
	}

	state, err := NewState(containerID, bundlePath)
	if err != nil {
		return Container{}, err
	}

	container := Container{
		State: state,
		Root:  root,
	}

	if err := container.save(); err != nil {
		return Container{}, err
	}

	return container, nil
}

func LoadContainer(containerID string) (Container, error) {
	root := filepath.Join(BASE_PATH, CONTAINER_PATH, containerID)
	statePath := filepath.Join(root, "state.json")

	stateData, err := os.ReadFile(statePath)
	if err != nil {
		return Container{}, err
	}

	var state specs.State
	if err := json.Unmarshal(stateData, &state); err != nil {
		return Container{}, err
	}

	return Container{
		State: state,
		Root:  root,
	}, nil
}

func (c Container) save() error {
	state, err := json.Marshal(c.State)
	if err != nil {
		return err
	}

	filePath := filepath.Join(c.Root, "state.json")

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filePath, state, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
