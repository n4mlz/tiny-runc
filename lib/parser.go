package lib

import (
	"encoding/json"
	"os"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func ParseConfig(path string) (specs.Spec, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return specs.Spec{}, err
	}

	var config specs.Spec
	if err := json.Unmarshal(file, &config); err != nil {
		return specs.Spec{}, err
	}

	return config, nil
}
