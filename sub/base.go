package sub

import (
	"fmt"
	"net/url"
)

type VlessBaseConfig struct {
	UUID       string
	RemoteHost string
	RemotePort uint32
}

func (c *VlessBaseConfig) Build() string {
	return fmt.Sprintf("%s@%s:%d", url.QueryEscape(c.UUID), c.RemoteHost, c.RemotePort)
}
