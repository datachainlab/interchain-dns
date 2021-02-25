package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	commontypes "github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
)

const (
	TypeRegisterChannelDomain string = "register_channel_domain"
	TypeDomainMappingCreate   string = "domain_mapping_create"
)

var _ sdk.Msg = (*MsgRegisterDomain)(nil)

// Route implements sdk.Msg
func (MsgRegisterDomain) Route() string {
	return commontypes.RouterKey
}

// Type implements sdk.Msg
func (MsgRegisterDomain) Type() string {
	return TypeRegisterChannelDomain
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterDomain) ValidateBasic() error {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgRegisterDomain) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgRegisterDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

var _ sdk.Msg = (*MsgDomainMappingCreate)(nil)

// Route implements sdk.Msg
func (MsgDomainMappingCreate) Route() string {
	return commontypes.RouterKey
}

// Type implements sdk.Msg
func (MsgDomainMappingCreate) Type() string {
	return TypeDomainMappingCreate
}

// ValidateBasic implements sdk.Msg
func (msg MsgDomainMappingCreate) ValidateBasic() error {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgDomainMappingCreate) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgDomainMappingCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
