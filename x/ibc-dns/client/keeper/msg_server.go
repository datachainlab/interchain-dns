package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/datachainlab/interchain-dns/x/ibc-dns/client/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) RegisterDomain(goCtx context.Context, msg *types.MsgRegisterDomain) (*types.MsgRegisterDomainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.SendPacketRegisterDomain(
		ctx,
		msg.Domain,
		msg.SourcePort,
		msg.SourceChannel,
		msg.Metadata,
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to send a packet 'PacketRegisterChannelDomain': %v", err)
	}

	k.Logger(ctx).Info(
		"IBC dns domain register",
		"domain", msg.Domain,
		"metadata", fmt.Sprintf("%s", msg.Metadata),
		"sender", msg.Sender,
	)

	return &types.MsgRegisterDomainResponse{}, nil

}

func (k Keeper) DomainMappingCreate(goCtx context.Context, msg *types.MsgDomainMappingCreate) (*types.MsgDomainMappingCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.SendDomainMappingCreatePacketData(
		ctx,
		msg.DnsId,
		msg.SrcClient,
		msg.DstClient,
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to send a packet 'DomainMappingCreatePacketData': %v", err)
	}

	k.Logger(ctx).Info(
		"IBC dns domain mapping create",
		"dnsId", msg.DnsId,
		"src client", msg.SrcClient,
		"dst client", msg.DstClient,
	)

	return &types.MsgDomainMappingCreateResponse{}, nil
}
