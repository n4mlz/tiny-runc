package lib

import (
	"encoding/json"
	"log"
	"os"
)

// https://github.com/opencontainers/runtime-spec/blob/main/config.md
type Config struct {
	Linux Linux `json:"linux"`
}

// https://github.com/opencontainers/runtime-spec/blob/main/config-linux.md
type Linux struct {
	UidMappings []struct {
		ContainerID int `json:"container_id"`
		HostID      int `json:"host_id"`
		Size        int `json:"size"`
	} `json:"uidMappings"`
	GidMappings []struct {
		ContainerID int `json:"container_id"`
		HostID      int `json:"host_id"`
		Size        int `json:"size"`
	} `json:"gidMappings"`
	Namespaces []struct {
		Type string `json:"type"`
	} `json:"namespaces"`
}

func ParseConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("failed to decode file: %v", err)
	}

	log.Printf("config: %+v", config)

	return config
}
