package types

import (
	"errors"
	math "math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

const (
	PacketTypeRegisterChannelDomain                  = "register_channel_domain"
	PacketTypeRegisterChannelDomainAcknowledgement   = "register_channel_domain_acknowledgement"
	PacketTypeDomainAssociationCreate                = "domain_association_create"
	PacketTypeDomainAssociationCreateAcknowledgement = "domain_association_create_acknowledgement"
	PacketTypeDomainAssociationResult                = "domain_association_result"
	PacketTypeDomainAssociationResultAcknowledgement = "domain_association_result_acknowledgement"
)

// Define RegisterDomainPacketData

var _ types.PacketDataI = (*RegisterDomainPacketData)(nil)

func NewRegisterDomainPacketData(name string, metadata []byte) RegisterDomainPacketData {
	return RegisterDomainPacketData{DomainName: name, Metadata: metadata}
}

func (p RegisterDomainPacketData) ValidateBasic() error {
	if p.DomainName == "" {
		return errors.New("Domain name must be set")
	}
	return nil
}

func (p RegisterDomainPacketData) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p RegisterDomainPacketData) Type() string {
	return PacketTypeRegisterChannelDomain
}

// Define RegisterDomainPacketAcknowledgement

const (
	STATUS_OK uint32 = iota + 1
	STATUS_FAILED
)

var _ types.PacketAcknowledgementI = (*RegisterDomainPacketAcknowledgement)(nil)

func NewRegisterDomainPacketAcknowledgement(status uint32, name string) RegisterDomainPacketAcknowledgement {
	return RegisterDomainPacketAcknowledgement{Status: status, DomainName: name}
}

func (p RegisterDomainPacketAcknowledgement) ValidateBasic() error {
	return nil
}

func (p RegisterDomainPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&p))
}

func (p RegisterDomainPacketAcknowledgement) Type() string {
	return PacketTypeRegisterChannelDomainAcknowledgement
}

// Define DomainAssociationCreatePacketData

var _ types.PacketDataI = (*DomainAssociationCreatePacketData)(nil)

func NewDomainAssociationCreatePacketData(srcClientDomain, dstClientDomain types.ClientDomain) DomainAssociationCreatePacketData {
	return DomainAssociationCreatePacketData{SrcClient: srcClientDomain, DstClient: dstClientDomain}
}

func (p DomainAssociationCreatePacketData) ValidateBasic() error {
	return nil
}

func (p DomainAssociationCreatePacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&p))
}

func (p DomainAssociationCreatePacketData) GetTimeoutHeight() uint64 {
	return math.MaxUint64
}

func (p DomainAssociationCreatePacketData) GetTimeoutTimestamp() uint64 {
	return 0
}

func (p DomainAssociationCreatePacketData) Type() string {
	return PacketTypeDomainAssociationCreate
}

// Define DomainAssociationCreatePacketAcknowledgement

var _ types.PacketAcknowledgementI = (*DomainAssociationCreatePacketAcknowledgement)(nil)

func NewDomainAssociationCreatePacketAcknowledgement(status uint32, msg string) DomainAssociationCreatePacketAcknowledgement {
	return DomainAssociationCreatePacketAcknowledgement{Status: status, Msg: msg}
}

func (p DomainAssociationCreatePacketAcknowledgement) ValidateBasic() error {
	return nil
}

func (p DomainAssociationCreatePacketAcknowledgement) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p DomainAssociationCreatePacketAcknowledgement) Type() string {
	return PacketTypeDomainAssociationCreateAcknowledgement
}

// Define DomainAssociationResultPacketData

var _ types.PacketDataI = (*DomainAssociationResultPacketData)(nil)

func NewDomainAssociationResultPacketData(
	status uint32,
	counterpartyDomainName string,
	dnsID types.LocalDNSID,
	clientID string,
) DomainAssociationResultPacketData {
	return DomainAssociationResultPacketData{
		Status:             status,
		CounterpartyDomain: types.NewLocalDomain(dnsID, counterpartyDomainName),
		ClientId:           clientID,
	}
}

func (p DomainAssociationResultPacketData) ValidateBasic() error {
	return nil
}

func (p DomainAssociationResultPacketData) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz

}

func (p DomainAssociationResultPacketData) GetTimeoutHeight() uint64 {
	return math.MaxUint64
}

func (p DomainAssociationResultPacketData) GetTimeoutTimestamp() uint64 {
	return 0
}

func (p DomainAssociationResultPacketData) Type() string {
	return PacketTypeDomainAssociationResult
}

// Define DomainAssociationResultPacketAcknowledgement

var _ types.PacketAcknowledgementI = (*DomainAssociationResultPacketAcknowledgement)(nil)

func NewDomainAssociationResultPacketAcknowledgement() DomainAssociationResultPacketAcknowledgement {
	return DomainAssociationResultPacketAcknowledgement{}
}

func (p DomainAssociationResultPacketAcknowledgement) ValidateBasic() error {
	return nil
}

func (p DomainAssociationResultPacketAcknowledgement) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p DomainAssociationResultPacketAcknowledgement) Type() string {
	return PacketTypeDomainAssociationResultAcknowledgement
}
