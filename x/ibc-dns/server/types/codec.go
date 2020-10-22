package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	types.RegisterCodec(ModuleCdc)
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(RegisterDomainPacketData{}, "ibc/dns/server/RegisterDomainPacketData", nil)
	cdc.RegisterConcrete(RegisterDomainPacketAcknowledgement{}, "ibc/dns/server/RegisterDomainPacketAcknowledgement", nil)

	cdc.RegisterConcrete(DomainAssociationCreatePacketData{}, "ibc/dns/server/DomainAssociationCreatePacketData", nil)
	cdc.RegisterConcrete(DomainAssociationCreatePacketAcknowledgement{}, "ibc/dns/server/DomainAssociationCreatePacketAcknowledgement", nil)

	cdc.RegisterConcrete(DomainAssociationResultPacketData{}, "ibc/dns/server/DomainAssociationResultPacketData", nil)
	cdc.RegisterConcrete(DomainAssociationResultPacketAcknowledgement{}, "ibc/dns/server/DomainAssociationResultPacketAcknowledgement", nil)
}
