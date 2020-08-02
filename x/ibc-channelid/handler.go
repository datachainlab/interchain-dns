package channelid

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/datachainlab/cosmos-sdk-interchain-channel-id/x/ibc-channelid/types"
)

// NewHandler returns a handler
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgRegisterChannel:
			return handleRegisterChannel(msg, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
		}
	}
}

func handleRegisterChannel(msg types.MsgRegisterChannel, keeper Keeper) (*sdk.Result, error) {
	panic("not implemented error")
}
