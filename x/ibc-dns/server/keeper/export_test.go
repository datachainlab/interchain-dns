package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

func SendPacket(k Keeper, ctx sdk.Context,
	packet *channeltypes.Packet) error {
	return k.sendPacket(ctx, packet)
}
