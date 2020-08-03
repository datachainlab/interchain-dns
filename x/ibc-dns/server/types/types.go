package types

import "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"

const (
	DomainAssociationStatusInit uint32 = iota + 1
	DomainAssociationStatusConfirmed
)

func NewDomainAssociation(status uint32, srcClient, dstClient types.ClientDomain) types.DomainAssociation {
	return types.DomainAssociation{Status: status, SrcClient: srcClient, DstClient: dstClient}
}
