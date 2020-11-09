package types

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/golang/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
)

type PacketReceiver func(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, []byte, error)

type PacketAcknowledgementReceiver func(ctx sdk.Context, packet channeltypes.Packet, ack []byte) (*sdk.Result, error)

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
			res, ack, err := r(ctx, packet)
			if err == nil {
				return res, ack, nil
			} else if err != ErrUnknownRequest {
				return res, ack, err
			}
		}
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC packet type: %T", packet)
	}
}

func ComposePacketAcknowledgementReceivers(rs ...PacketAcknowledgementReceiver) PacketAcknowledgementReceiver {
	return func(ctx sdk.Context, packet channeltypes.Packet, ack []byte) (*sdk.Result, error) {
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

func SerializeJSONPacketData(cdc codec.JSONMarshaler, msg proto.Message) ([]byte, error) {
	bz, err := MarshalJSONAny(cdc, msg)
	if err != nil {
		return nil, err
	}
	return sdk.SortJSON(bz)
}

func DeserializeJSONPacketData(cdc codec.Marshaler, data []byte) (PacketDataI, error) {
	var pd PacketDataI
	err := UnmarshalJSONAny(cdc, &pd, data)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func MustSerializeJSONPacketData(cdc codec.JSONMarshaler, msg proto.Message) []byte {
	bz, err := SerializeJSONPacketData(cdc, msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func MustDeserializeJSONPacketData(cdc codec.Marshaler, data []byte) PacketDataI {
	pd, err := DeserializeJSONPacketData(cdc, data)
	if err != nil {
		panic(err)
	}
	return pd
}

func SerializeJSONPacketAck(cdc codec.JSONMarshaler, msg proto.Message) ([]byte, error) {
	bz, err := MarshalJSONAny(cdc, msg)
	if err != nil {
		return nil, err
	}
	return sdk.SortJSON(bz)
}

func DeserializeJSONPacketAck(cdc codec.Marshaler, data []byte) (PacketAcknowledgementI, error) {
	var pd PacketAcknowledgementI
	err := UnmarshalJSONAny(cdc, &pd, data)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func MustSerializeJSONPacketAck(cdc codec.JSONMarshaler, msg proto.Message) []byte {
	bz, err := SerializeJSONPacketAck(cdc, msg)
	if err != nil {
		panic(err)
	}
	return bz
}

func MustDeserializeJSONPacketAck(cdc codec.Marshaler, data []byte) PacketAcknowledgementI {
	ack, err := DeserializeJSONPacketAck(cdc, data)
	if err != nil {
		panic(err)
	}
	return ack
}
