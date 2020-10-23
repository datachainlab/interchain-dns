package server

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

func NewPacketReceiver(keeper Keeper) commontypes.PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
		var data commontypes.PacketData
		if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
		}
		switch data := data.(type) {
		case types.RegisterDomainPacketData:
			return handlePacketRegisterChannelDomain(ctx, keeper, packet, data)
		case types.DomainAssociationCreatePacketData:
			return handleDomainAssociationCreatePacketData(ctx, keeper, packet, data)
		default:
			return nil, commontypes.ErrUnknownRequest
		}
	}
}

func handlePacketRegisterChannelDomain(ctx sdk.Context, keeper Keeper, packet channeltypes.Packet, data types.RegisterDomainPacketData) (*sdk.Result, error) {
	var status uint32
	if err := keeper.ReceivePacketRegisterDomain(ctx, packet, data); err != nil {
		ctx.Logger().Info("failed to handle a packet 'PacketRegisterChannelDomain'", "err", err)
		status = types.STATUS_FAILED
	} else {
		status = types.STATUS_OK
	}
	ack := types.NewRegisterDomainPacketAcknowledgement(status, data.DomainName)
	if err := keeper.PacketExecuted(ctx, packet, ack.GetBytes()); err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleDomainAssociationCreatePacketData(ctx sdk.Context, keeper Keeper, packet channeltypes.Packet, data types.DomainAssociationCreatePacketData) (*sdk.Result, error) {
	ack, completed := keeper.ReceiveDomainAssociationCreatePacketData(ctx, packet, data)
	if completed {
		err := keeper.SendDomainAssociationResultPacketData(ctx, ack.Status, data.SrcClient, data.DstClient)
		if err != nil {
			return nil, err
		}
	}
	if err := keeper.PacketExecuted(ctx, packet, ack.GetBytes()); err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func NewPacketAcknowledgementReceiver(keeper Keeper) commontypes.PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack commontypes.PacketAcknowledgement) (*sdk.Result, error) {
		var data commontypes.PacketData
		if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
		}
		switch ack := ack.(type) {
		case types.DomainAssociationResultPacketAcknowledgement:
			return handleDomainAssociationResultPacketAcknowledgement(ctx, keeper, ack)
		default:
			return nil, commontypes.ErrUnknownRequest
		}
	}
}

func handleDomainAssociationResultPacketAcknowledgement(ctx sdk.Context, k Keeper, ack types.DomainAssociationResultPacketAcknowledgement) (*sdk.Result, error) {
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
