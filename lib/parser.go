package lib

import (
	"encoding/json"
	"log"
	"os"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func ParseConfig(path string) specs.Spec {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var config specs.Spec
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("failed to decode file: %v", err)
	}

	return config
}
