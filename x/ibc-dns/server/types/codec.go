package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

// RegisterInterfaces register the ibc transfer module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ibc.dns.server.v1.PacketData",
		(*types.PacketDataI)(nil),
		&RegisterDomainPacketData{},
		&DomainAssociationCreatePacketData{},
		&DomainAssociationResultPacketData{},
	)

	registry.RegisterInterface(
		"ibc.dns.server.v1.PacketAcknowledgement",
		(*types.PacketAcknowledgementI)(nil),
		&RegisterDomainPacketAcknowledgement{},
		&DomainAssociationCreatePacketAcknowledgement{},
		&DomainAssociationResultPacketAcknowledgement{},
	)
}

var (
	// ModuleCdc references the global x/ibc-transfer module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc-transfer and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)
