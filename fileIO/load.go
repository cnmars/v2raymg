package fileIO

import (
	"encoding/json"
	"io/ioutil"

	"github.com/lureiny/v2raymg/protocol"
)

// LoadConfig load config from file
func LoadConfig(file string) (*protocol.V2rayConfig, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config protocol.V2rayConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
