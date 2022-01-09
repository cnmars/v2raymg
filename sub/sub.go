package sub

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lureiny/v2raymg/fileIO"
)

// GetUserSubUri 获取某个指定用户的订阅uri
func GetUserSubUri(host string, user string, port uint32, configFile string) (string, error) {
	c, err := fileIO.LoadConfig(configFile)
	if err != nil {
		return "", err
	}

	for _, in := range c.InboundConfigs {
		switch strings.ToLower(in.Protocol) {
		case "vless":
			u, err := NewVlessShareConfig(&in, user, host, port)
			if err != nil {
				continue
			}
			return getVlessUri(u)

		case "vmess":
			u, err := NewVmessShareConfig(&in, user, host, port)
			if err != nil {
				continue
			}
			return getVmessUri(u)
		}
	}
	errMsg := "No User"
	return "", errors.New(errMsg)
}

func getVmessUri(u *VmessShareConfig) (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("vmess://%s", base64.StdEncoding.EncodeToString(b)), nil
}

func getVlessUri(v *VlessShareConfig) (string, error) {
	uri := fmt.Sprintf("vless://%s", v.Build())
	return uri, nil
}
