package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

type Keeper struct {
	portKeeper   types.PortKeeper
	scopedKeeper capabilitykeeper.ScopedKeeper
}

func NewKeeper(portKeeper types.PortKeeper, scopedKeeper capabilitykeeper.ScopedKeeper) Keeper {
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

// IsBound checks if the transfer module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}
