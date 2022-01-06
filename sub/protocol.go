package sub

import (
	"errors"
	"fmt"
	"sort"

	"github.com/v2fly/v2ray-core/v4/infra/conf"
)

type VlessProtocolConfig struct {
	Type       string
	Encryption string
}

func (c *VlessProtocolConfig) Build() string {
	// 参考https://github.com/XTLS/Xray-core/discussions/716，暂时省略"encryption"字段
	// 不支持vmess指定加密方式
	return fmt.Sprintf("type=%s", c.Type)
}

func inProtocols(p string) bool {
	protocols := []string{"grpc", "http", "quic", "tcp", "ws"}
	index := sort.SearchStrings(protocols, p)
	return index != len(protocols)
}

func newProtocolConfig(streamSetting *conf.StreamConfig) (*VlessProtocolConfig, error) {
	if inProtocols(string(*streamSetting.Network)) {
		return &VlessProtocolConfig{Type: string(*streamSetting.Network), Encryption: "none"}, nil
	}
	errMsg := fmt.Sprintf("unsupoort network config: %v", string(*streamSetting.Network))
	return nil, errors.New(errMsg)
}
