package types

import (
	"errors"
	"math"

	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	"github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
)

const (
	PacketTypeRegisterChannelDomain                = "register_channel_domain"
	PacketTypeRegisterChannelDomainAcknowledgement = "register_channel_domain_acknowledgement"
	PacketTypeDomainMappingCreate                  = "domain_mapping_create"
	PacketTypeDomainMappingCreateAcknowledgement   = "domain_mapping_create_acknowledgement"
	PacketTypeDomainMappingResult                  = "domain_mapping_result"
	PacketTypeDomainMappingResultAcknowledgement   = "domain_mapping_result_acknowledgement"
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

func (p RegisterDomainPacketData) GetTimeoutHeight() ibcclienttypes.Height {
	return ibcclienttypes.NewHeight(0, math.MaxInt64)
}

func (p RegisterDomainPacketData) GetTimeoutTimestamp() uint64 {
	return 0
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
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p RegisterDomainPacketAcknowledgement) Type() string {
	return PacketTypeRegisterChannelDomainAcknowledgement
}

// Define DomainMappingCreatePacketData

var _ types.PacketDataI = (*DomainMappingCreatePacketData)(nil)

func NewDomainMappingCreatePacketData(srcClientDomain, dstClientDomain types.ClientDomain) DomainMappingCreatePacketData {
	return DomainMappingCreatePacketData{SrcClient: srcClientDomain, DstClient: dstClientDomain}
}

func (p DomainMappingCreatePacketData) ValidateBasic() error {
	return nil
}

func (p DomainMappingCreatePacketData) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p DomainMappingCreatePacketData) GetTimeoutHeight() ibcclienttypes.Height {
	return ibcclienttypes.NewHeight(0, math.MaxInt64)
}

func (p DomainMappingCreatePacketData) GetTimeoutTimestamp() uint64 {
	return 0
}

func (p DomainMappingCreatePacketData) Type() string {
	return PacketTypeDomainMappingCreate
}

// Define DomainMappingCreatePacketAcknowledgement

var _ types.PacketAcknowledgementI = (*DomainMappingCreatePacketAcknowledgement)(nil)

func NewDomainMappingCreatePacketAcknowledgement(status uint32, msg string) DomainMappingCreatePacketAcknowledgement {
	return DomainMappingCreatePacketAcknowledgement{Status: status, Msg: msg}
}

func (p DomainMappingCreatePacketAcknowledgement) ValidateBasic() error {
	return nil
}

func (p DomainMappingCreatePacketAcknowledgement) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p DomainMappingCreatePacketAcknowledgement) Type() string {
	return PacketTypeDomainMappingCreateAcknowledgement
}

// Define DomainMappingResultPacketData

var _ types.PacketDataI = (*DomainMappingResultPacketData)(nil)

func NewDomainMappingResultPacketData(
	status uint32,
	counterpartyDomainName string,
	dnsID types.LocalDNSID,
	clientID string,
) DomainMappingResultPacketData {
	return DomainMappingResultPacketData{
		Status:             status,
		CounterpartyDomain: types.NewLocalDomain(dnsID, counterpartyDomainName),
		ClientId:           clientID,
	}
}

func (p DomainMappingResultPacketData) ValidateBasic() error {
	return nil
}

func (p DomainMappingResultPacketData) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p DomainMappingResultPacketData) GetTimeoutHeight() ibcclienttypes.Height {
	return ibcclienttypes.NewHeight(0, math.MaxInt64)
}

func (p DomainMappingResultPacketData) GetTimeoutTimestamp() uint64 {
	return 0
}

func (p DomainMappingResultPacketData) Type() string {
	return PacketTypeDomainMappingResult
}

// Define DomainMappingResultPacketAcknowledgement

var _ types.PacketAcknowledgementI = (*DomainMappingResultPacketAcknowledgement)(nil)

func NewDomainMappingResultPacketAcknowledgement() DomainMappingResultPacketAcknowledgement {
	return DomainMappingResultPacketAcknowledgement{}
}

func (p DomainMappingResultPacketAcknowledgement) ValidateBasic() error {
	return nil
}

func (p DomainMappingResultPacketAcknowledgement) GetBytes() []byte {
	bz, err := types.SerializeJSONPacketData(PacketCdc(), &p)
	if err != nil {
		panic(err)
	}
	return bz
}

func (p DomainMappingResultPacketAcknowledgement) Type() string {
	return PacketTypeDomainMappingResultAcknowledgement
}
