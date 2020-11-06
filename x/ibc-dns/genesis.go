package dns

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/keeper"

	abci "github.com/tendermint/tendermint/abci/types"
)

// NewGenesisState is a constructor of GenesisState
func NewGenesisState(master string) commontypes.GenesisState {
	return commontypes.GenesisState{}
}

// ValidateGenesis checks the Genesis
func ValidateGenesis(data commontypes.GenesisState) error {
	return nil
}

// DefaultGenesisState returns default genesis state
func DefaultGenesisState() *commontypes.GenesisState {
	return &commontypes.GenesisState{}
}

// InitGenesis inits genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data commontypes.GenesisState) []abci.ValidatorUpdate {
	_, err := keeper.BindPort(ctx, commontypes.PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) commontypes.GenesisState {
	return commontypes.GenesisState{}
}
