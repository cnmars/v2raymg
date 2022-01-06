package sub

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/v2fly/v2ray-core/v4/infra/conf"
)

// VlessTransportConfig 传输层最上层配置结构
type VlessTransportConfig struct {
	Security        string
	TransportConfig URIAdapter
}

func (c *VlessTransportConfig) Build() string {
	params := []string{}
	if len(c.Security) > 0 {
		params = append(params, "security="+c.Security)
	}
	params = append(params, c.TransportConfig.Build())
	return strings.Join(params, "&")
}

type VlessTcpConfig struct {
}

func (c *VlessTcpConfig) Build() string {
	return ""
}

// VlessHttp2Config http2配置结构
type VlessHttp2Config struct {
	Path string
	Host string
}

func (c *VlessHttp2Config) Build() string {
	params := []string{}
	if len(c.Path) > 0 {
		params = append(params, "path="+c.Path)
	}
	if len(c.Host) > 0 {
		params = append(params, "host="+c.Host)
	}
	return strings.Join(params, "&")
}

// VlessWebSocketConfig ws配置结构
type VlessWebSocketConfig struct {
	Path string
	Host string
}

func (c *VlessWebSocketConfig) Build() string {
	params := []string{}
	if len(c.Path) > 0 {
		params = append(params, "path="+c.Path)
	}
	if len(c.Host) > 0 {
		params = append(params, "host="+c.Host)
	}
	return strings.Join(params, "&")
}

// VlessMkcpConfig mkcp配置结构
type VlessMkcpConfig struct {
	HeaderType string
	Seed       string
}

func (c *VlessMkcpConfig) Build() string {
	params := []string{}
	if len(c.HeaderType) > 0 {
		params = append(params, "headerType="+c.HeaderType)
	}
	if len(c.Seed) > 0 {
		params = append(params, "seed="+c.Seed)
	}
	return strings.Join(params, "&")
}

// VlessQuicConfig quic配置结构
type VlessQuicConfig struct {
	QuicSecurity string
	Key          string
}

func (c *VlessQuicConfig) Build() string {
	if c.QuicSecurity == "none" || len(c.QuicSecurity) == 0 {
		return ""
	}
	return fmt.Sprintf("quicSecurity=%s&key=%s", c.QuicSecurity, c.Key)
}

// VlessGrpcConfig grpc配置结构
type VlessGrpcConfig struct {
	ServiceName string
	Mode        string
}

func (c *VlessGrpcConfig) Build() string {
	params := []string{}
	if len(c.ServiceName) > 0 {
		params = append(params, "serviceName="+c.ServiceName)
	}
	if len(c.Mode) > 0 {
		params = append(params, "mode="+c.Mode)
	}
	return strings.Join(params, "&")
}

func newTransportConfig(streamSetting *conf.StreamConfig) (*VlessTransportConfig, error) {
	transportConfig := VlessTransportConfig{
		Security: streamSetting.Security,
	}
	// TODO: 支持http2
	switch string(*streamSetting.Network) {
	case "tcp":
		transportConfig.TransportConfig = &VlessTcpConfig{}
	case "kcp":
		kcpConfig := streamSetting.KCPSettings
		var clientMkcpConfig VlessMkcpConfig
		var kcpHeader map[string]string
		if len(kcpConfig.HeaderConfig) > 0 {
			err := json.Unmarshal(kcpConfig.HeaderConfig, &kcpHeader)
			if err != nil {
				return nil, errors.New("invalid mKCP header config.")
			}
			if err != nil {
				return nil, errors.New("invalid mKCP header config")
			}
			if headerType, ok := kcpHeader["type"]; ok {
				clientMkcpConfig.HeaderType = headerType
			}
		}

		if kcpConfig.Seed != nil {
			clientMkcpConfig.Seed = *kcpConfig.Seed
		}
		transportConfig.TransportConfig = &clientMkcpConfig
	case "ws":
		var wsConfig VlessWebSocketConfig
		wsConfig.Host = streamSetting.WSSettings.Headers["Host"]
		wsConfig.Path = streamSetting.WSSettings.Path
		transportConfig.TransportConfig = &wsConfig
	case "quic":
		var quicConfig VlessQuicConfig
		quicConfig.Key = streamSetting.QUICSettings.Key
		quicConfig.QuicSecurity = streamSetting.QUICSettings.Security
		transportConfig.TransportConfig = &quicConfig
	case "grpc":
		var grpcConfig VlessGrpcConfig
		grpcConfig.ServiceName = streamSetting.GRPCSettings.ServiceName
		// 暂时不考虑xray-core
		grpcConfig.Mode = "gun"
		transportConfig.TransportConfig = &grpcConfig
	}
	transportConfig.Security = streamSetting.Security
	return &transportConfig, nil
}
