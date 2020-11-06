package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	clienttypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/types"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
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
	commontypes.RegisterLegacyAminoCodec(cdc)
	clienttypes.RegisterLegacyAminoCodec(cdc)
	servertypes.RegisterLegacyAminoCodec(cdc)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	commontypes.RegisterInterfaces(registry)
	clienttypes.RegisterInterfaces(registry)
	servertypes.RegisterInterfaces(registry)
}
