package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type PacketReceiver func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error)

type PacketAcknowledgementReceiver func(ctx sdk.Context, packet channeltypes.Packet, ack PacketAcknowledgement) (*sdk.Result, error)

var ErrUnknownRequest = errors.New("unknown request error")

func ComposeHandlers(hs ...sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		for _, h := range hs {
			res, err := h(ctx, msg)
			if err == nil {
				return res, nil
			} else if err != ErrUnknownRequest {
				return res, err
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
	}
}

func ComposeQuerier(qs ...sdk.Querier) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		for _, q := range qs {
			res, err := q(ctx, path, req)
			if err == nil {
				return res, nil
			} else if err != ErrUnknownRequest {
				return res, err
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized request: %v %v", path, req)
	}
}

func ComposePacketReceivers(rs ...PacketReceiver) PacketReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error) {
		for _, r := range rs {
			res, packet, err := r(ctx, packet)
			if err == nil {
				return res, packet, nil
			} else if err != ErrUnknownRequest {
				return res, packet, err
			}
		}
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
	}
}

func ComposePacketAcknowledgementReceivers(rs ...PacketAcknowledgementReceiver) PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack PacketAcknowledgement) (*sdk.Result, error) {
		for _, r := range rs {
			res, err := r(ctx, packet, ack)
			if err == nil {
				return res, nil
			} else if err != ErrUnknownRequest {
				return res, err
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
	}
}
