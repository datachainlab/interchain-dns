package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	portkeeper "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/keeper"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
)

type Keeper struct {
	portKeeper   portkeeper.Keeper
	scopedKeeper capabilitykeeper.ScopedKeeper
}

func NewKeeper(portKeeper portkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper) Keeper {
	return Keeper{
		portKeeper:   portKeeper,
		scopedKeeper: scopedKeeper,
	}
}

// BindPort defines a wrapper function for the ort TPCKeeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) (*capabilitytypes.Capability, error) {
	cap := k.portKeeper.BindPort(ctx, portID)
	if err := k.ClaimCapability(ctx, cap, host.PortPath(portID)); err != nil {
		return nil, err
	}
	return cap, nil
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
