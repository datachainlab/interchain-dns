package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
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
	cdc.RegisterConcrete(RegisterDomainPacketData{}, "ibc/dns/server/RegisterDomainPacketData", nil)
	cdc.RegisterConcrete(RegisterDomainPacketAcknowledgement{}, "ibc/dns/server/RegisterDomainPacketAcknowledgement", nil)
	cdc.RegisterConcrete(DomainAssociationCreatePacketData{}, "ibc/dns/server/DomainAssociationCreatePacketData", nil)
	cdc.RegisterConcrete(DomainAssociationCreatePacketAcknowledgement{}, "ibc/dns/server/DomainAssociationCreatePacketAcknowledgement", nil)
	cdc.RegisterConcrete(DomainAssociationResultPacketData{}, "ibc/dns/server/DomainAssociationResultPacketData", nil)
	cdc.RegisterConcrete(DomainAssociationResultPacketAcknowledgement{}, "ibc/dns/server/DomainAssociationResultPacketAcknowledgement", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*types.PacketDataI)(nil),
		&RegisterDomainPacketData{},
		&DomainAssociationCreatePacketData{},
		&DomainAssociationResultPacketData{},
	)

	registry.RegisterImplementations(
		(*types.PacketAcknowledgementI)(nil),
		&RegisterDomainPacketAcknowledgement{},
		&DomainAssociationCreatePacketAcknowledgement{},
		&DomainAssociationResultPacketAcknowledgement{},
	)
}
