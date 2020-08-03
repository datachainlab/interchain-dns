package client

import (
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/keeper"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/types"
)

type (
	Keeper = keeper.Keeper
)

var (
	StoreKey  = types.StoreKey
	NewKeeper = keeper.NewKeeper
)
