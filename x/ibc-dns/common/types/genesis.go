package types

// NewGenesisState is a constructor of GenesisState
func NewGenesisState(master string) *GenesisState {
	return &GenesisState{}
}

// DefaultGenesisState returns default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs *GenesisState) Validate() error {
	return nil
}

// // InitGenesis inits genesis
// func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
// 	_, err := keeper.BindPort(ctx, commontypes.PortID)
// 	if err != nil {
// 		panic(fmt.Sprintf("could not claim port capability: %v", err))
// 	}
// 	return []abci.ValidatorUpdate{}
// }

// // ExportGenesis exports genesis
// func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
// 	return GenesisState{}
// }
