package metadata

import (
	"fmt"
	"net"
)

// IPNet is a struct that embeds net.IPNet and adds support for
// unmarshalling of CIDR-style IP ranges.
type IPNet struct {
	net.IPNet
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The input for this function is expected in CIDR format.
func (ipNet *IPNet) UnmarshalText(text []byte) error {
	_, parsedIP, err := net.ParseCIDR(string(text))
	if err != nil {
		return fmt.Errorf("failed to parse linode cidr ip: %w", err)
	}

	ipNet.IP = parsedIP.IP
	ipNet.Mask = parsedIP.Mask

	return nil
}
