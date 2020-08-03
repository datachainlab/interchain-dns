package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	types.RegisterCodec(ModuleCdc)
	servertypes.RegisterCodec(ModuleCdc)
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterDomain{}, "ibc/dns/client/MsgRegisterDomain", nil)
	cdc.RegisterConcrete(MsgDomainAssociationCreate{}, "ibc/dns/client/MsgDomainAssociationCreate", nil)
}
