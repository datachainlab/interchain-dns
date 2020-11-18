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
