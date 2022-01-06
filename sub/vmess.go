package sub

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lureiny/v2raymg/fileIO"
	"github.com/v2fly/v2ray-core/v4/common/protocol"
	"github.com/v2fly/v2ray-core/v4/infra/conf"
)

type VmessShareConfig struct {
	V    string `json:"v"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	Sni  string `json:"sni"`
}

func NewVmessShareConfig(in *fileIO.InboundDetourConfig, email string, host string, port uint32) (*VmessShareConfig, error) {
	// 获取UUID
	id, aid, err := getVmessUserUUID(in, email)
	if err != nil {
		return nil, err
	}

	if host == "" {
		host = in.ListenOn
	}

	if port == 0 {
		port = in.PortRange
	}

	v := NewDefaultVmessShareConfig()
	v.Add = host
	v.Port = fmt.Sprint(port)
	v.ID = id
	v.Aid = fmt.Sprint(aid)
	v.TLS = "tls"

	if err := insertVmessStreamSetting(v, in.StreamSetting); err != nil {
		return nil, err
	}
	return v, nil
}

// TODO(lureiny): 针对不同协议底层适配
func insertVmessStreamSetting(v *VmessShareConfig, streamSetting *conf.StreamConfig) error {
	switch string(*streamSetting.Network) {
	case "tcp":
		v.Net = "tcp"
	case "kcp":
		v.Net = "kcp"
		v.Path = *streamSetting.KCPSettings.Seed
	case "ws":
		v.Net = "ws"
		v.Path = streamSetting.WSSettings.Path
	case "http":
		v.Net = "http"
	case "quic":
		v.Net = "quic"
	case "grpc":
		v.Net = "grpc"
	default:
		return errors.New(fmt.Sprintf("Unsupport transport protocol %s", *streamSetting.Network))
	}
	return nil
}

func NewDefaultVmessShareConfig() *VmessShareConfig {
	return &VmessShareConfig{V: "2"}
}

func getVmessUserUUID(in *fileIO.InboundDetourConfig, email string) (string, int, error) {
	vmessConfig := new(conf.VMessInboundConfig)

	err := json.Unmarshal([]byte(*(in.Settings)), vmessConfig)
	if err != nil {
		return "", 0, err
	}

	for _, rawData := range vmessConfig.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil || user.Email != email {
			continue
		}
		account := new(conf.VMessAccount)
		json.Unmarshal(rawData, account)
		return account.ID, int(account.AlterIds), nil

	}

	return "", 0, errors.New("%d not in")
}
