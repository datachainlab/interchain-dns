package types

import "github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"

const (
	DomainMappingStatusInit uint32 = iota + 1
	DomainMappingStatusConfirmed
)

func NewDomainMapping(status uint32, srcClient, dstClient types.ClientDomain) types.DomainMapping {
	return types.DomainMapping{Status: status, SrcClient: srcClient, DstClient: dstClient}
}
