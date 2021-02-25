package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
)

// InitGenesis initializes the ibc-transfer state and binds to PortID.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {}

// ExportGenesis exports ibc-transfer module's portID and denom trace info into its genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
