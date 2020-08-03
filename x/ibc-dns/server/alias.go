package server

import (
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/keeper"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

type (
	Keeper = keeper.Keeper
)

var (
	StoreKey  = types.StoreKey
	NewKeeper = keeper.NewKeeper
)
