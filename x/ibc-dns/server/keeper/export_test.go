package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

func SendPacket(k Keeper, ctx sdk.Context,
	packet *channel.Packet) error {
	return k.sendPacket(ctx, packet)
}
