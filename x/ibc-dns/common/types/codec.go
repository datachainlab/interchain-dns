package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*PacketDataI)(nil), nil)
	cdc.RegisterInterface((*PacketAcknowledgementI)(nil), nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"datachainlab.ibc.dns.v1beta1.PacketDataI",
		(*PacketDataI)(nil),
	)

	registry.RegisterInterface(
		"datachainlab.ibc.dns.v1beta1.PacketAcknowledgementI",
		(*PacketAcknowledgementI)(nil),
	)
}
