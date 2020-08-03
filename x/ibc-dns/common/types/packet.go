package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

type PacketReceiver func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error)

type PacketAcknowledgementReceiver func(ctx sdk.Context, packet channeltypes.Packet, ack PacketAcknowledgement) (*sdk.Result, error)

var ErrUnknownRequest = errors.New("unknown request error")

func ComposeHandlers(hs ...sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		for _, h := range hs {
			res, err := h(ctx, msg)
			if err != ErrUnknownRequest {
				return res, err
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
	}
}

func ComposePacketReceivers(rs ...PacketReceiver) PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
		for _, r := range rs {
			res, err := r(ctx, packet)
			if err != ErrUnknownRequest {
				return res, err
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
	}
}

func ComposePacketAcknowledgementReceivers(rs ...PacketAcknowledgementReceiver) PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack PacketAcknowledgement) (*sdk.Result, error) {
		for _, r := range rs {
			res, err := r(ctx, packet, ack)
			if err != ErrUnknownRequest {
				return res, err
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
	}
}
