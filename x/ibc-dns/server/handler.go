package server

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

func NewPacketReceiver(keeper Keeper) commontypes.PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error) {
		var data commontypes.PacketDataI
		if err := codec.UnmarshalAny(keeper.Codec(), &data, packet.GetData()); err != nil {
			return nil, nil, err
		}
		switch data := data.(type) {
		case *types.RegisterDomainPacketData:
			return handlePacketRegisterChannelDomain(ctx, keeper, packet, data)
		case *types.DomainAssociationCreatePacketData:
			return handleDomainAssociationCreatePacketData(ctx, keeper, packet, data)
		default:
			return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet data type: %T", data)
		}
	}
}

func handlePacketRegisterChannelDomain(ctx sdk.Context, keeper Keeper, packet channeltypes.Packet, data *types.RegisterDomainPacketData) (*sdk.Result, []byte, error) {
	var status uint32
	if err := keeper.ReceivePacketRegisterDomain(ctx, packet, data); err != nil {
		ctx.Logger().Info("failed to handle a packet 'PacketRegisterChannelDomain'", "err", err)
		status = types.STATUS_FAILED
	} else {
		status = types.STATUS_OK
	}
	ack := types.NewRegisterDomainPacketAcknowledgement(status, data.DomainName)
	if err := keeper.PacketExecuted(ctx, packet, ack.GetBytes()); err != nil {
		return nil, nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, ack.GetBytes(), nil
}

func handleDomainAssociationCreatePacketData(ctx sdk.Context, keeper Keeper, packet channeltypes.Packet, data *types.DomainAssociationCreatePacketData) (*sdk.Result, []byte, error) {
	ack, completed := keeper.ReceiveDomainAssociationCreatePacketData(ctx, packet, data)
	if completed {
		err := keeper.SendDomainAssociationResultPacketData(ctx, ack.Status, data.SrcClient, data.DstClient)
		if err != nil {
			return nil, nil, err
		}
	}
	if err := keeper.PacketExecuted(ctx, packet, ack.GetBytes()); err != nil {
		return nil, nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, ack.GetBytes(), nil
}

func NewPacketAcknowledgementReceiver(keeper Keeper) commontypes.PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack []byte) (*sdk.Result, error) {
		var ackData commontypes.PacketAcknowledgementI
		if err := codec.UnmarshalAny(keeper.Codec(), &ackData, ack); err != nil {
			return nil, err
		}
		switch ackData := ackData.(type) {
		case *types.DomainAssociationResultPacketAcknowledgement:
			return handleDomainAssociationResultPacketAcknowledgement(ctx, keeper, ackData)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet data type: %T", ackData)
		}
	}
}

func handleDomainAssociationResultPacketAcknowledgement(ctx sdk.Context, k Keeper, ackData *types.DomainAssociationResultPacketAcknowledgement) (*sdk.Result, error) {
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
