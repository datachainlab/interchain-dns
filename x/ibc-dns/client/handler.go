package client

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/types"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// NewHandler returns a handler
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgRegisterDomain:
			return handleRegisterDomain(ctx, msg, keeper)
		case *types.MsgDomainAssociationCreate:
			return handleDomainAssociationCreate(ctx, msg, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
		}
	}
}

func handleRegisterDomain(ctx sdk.Context, msg *types.MsgRegisterDomain, keeper Keeper) (*sdk.Result, error) {
	_, err := keeper.SendPacketRegisterDomain(ctx, msg.Domain, msg.SourcePort, msg.SourceChannel, msg.Metadata)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to send a packet 'PacketRegisterChannelDomain': %v", err)
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleDomainAssociationCreate(ctx sdk.Context, msg *types.MsgDomainAssociationCreate, keeper Keeper) (*sdk.Result, error) {
	_, err := keeper.SendDomainAssociationCreatePacketData(ctx, msg.DnsId, msg.SrcClient, msg.DstClient)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to send a packet 'DomainAssociationCreatePacketData': %v", err)
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// NewPacketReceiver returns a new PacketReceiver
func NewPacketReceiver(keeper Keeper) commontypes.PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error) {
		data, err := commontypes.DeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData())
		if err != nil {
			return nil, nil, err
		}
		switch data := data.(type) {
		case *servertypes.DomainAssociationResultPacketData:
			return handleDomainAssociationResultPacketData(ctx, keeper, packet, data)
		default:
			return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet data type: %T", data)
		}
	}
}

func handleDomainAssociationResultPacketData(ctx sdk.Context, keeper Keeper, packet channeltypes.Packet, data *servertypes.DomainAssociationResultPacketData) (*sdk.Result, []byte, error) {
	switch data.Status {
	case servertypes.STATUS_OK:
		err := keeper.ReceiveDomainAssociationResultPacketData(
			ctx,
			packet,
			data,
		)
		if err != nil {
			return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to handle a packet 'DomainAssociationResultPacketData: %v'", err)
		}
	case servertypes.STATUS_FAILED:
		// TODO cleanup
	default:
		return nil, nil, fmt.Errorf("unknown status '%v'", data.Status)
	}

	ack := servertypes.NewDomainAssociationResultPacketAcknowledgement()
	if err := keeper.PacketExecuted(ctx, packet, ack.GetBytes()); err != nil {
		return nil, nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, ack.GetBytes(), nil
}

// NewPacketAcknowledgementReceiver returns a new PacketAcknowledgementReceiver
func NewPacketAcknowledgementReceiver(keeper Keeper) commontypes.PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack []byte) (*sdk.Result, error) {
		ackData, err := commontypes.DeserializeJSONPacketAck(servertypes.PacketCdc(), packet.GetData())
		if err != nil {
			return nil, err
		}
		switch ackData := ackData.(type) {
		case *servertypes.RegisterDomainPacketAcknowledgement:
			return handleRegisterDomainPacketAcknowledgement(ctx, keeper, ackData, packet)
		case *servertypes.DomainAssociationCreatePacketAcknowledgement:
			return handleDomainAssociationCreatePacketAcknowledgement(ctx, keeper, ackData)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet data type: %T", ackData)
		}
	}
}

func handleRegisterDomainPacketAcknowledgement(ctx sdk.Context, k Keeper, ack *servertypes.RegisterDomainPacketAcknowledgement, packet channeltypes.Packet) (*sdk.Result, error) {
	if err := k.ReceiveRegisterDomainPacketAcknowledgement(ctx, ack.Status, ack.DomainName, packet); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to handle a packet 'RegisterDomainPacketAcknowledgement: %v'", err)
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleDomainAssociationCreatePacketAcknowledgement(ctx sdk.Context, k Keeper, ack *servertypes.DomainAssociationCreatePacketAcknowledgement) (*sdk.Result, error) {
	if err := k.ReceiveDomainAssociationCreatePacketAcknowledgement(ctx, ack.Status); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to handle a packet 'DomainAssociationCreatePacketAcknowledgement: %v'", err)
	}
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
