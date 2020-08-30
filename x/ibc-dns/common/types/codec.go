package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterInterfaces registers the client interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ibc.dns.v1.common.PacketData",
		(*PacketDataI)(nil),
	)
	registry.RegisterInterface(
		"ibc.dns.v1.common.PacketAcknowledgement",
		(*PacketAcknowledgementI)(nil),
	)
}
