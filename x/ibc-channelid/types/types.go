package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeRegisterChannel string = "register_channel"

var _ sdk.Msg = (*MsgRegisterChannel)(nil)

// Route implements sdk.Msg
func (MsgRegisterChannel) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgRegisterChannel) Type() string {
	return TypeRegisterChannel
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterChannel) ValidateBasic() error {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgRegisterChannel) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgRegisterChannel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
