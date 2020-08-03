package types

import (
	"fmt"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

const (
	// StoreKey to be used when creating the KVStore
	StoreKey = "ibc-dns-server"
)

const (
	KeyForwardDomainPrefix uint8 = iota
	KeyReverseDomainPrefix
	KeyDomainAssociationPrefix
)

// KeyPrefixBytes return the key prefix bytes from a URL string format
func KeyPrefixBytes(prefix uint8) []byte {
	return []byte(fmt.Sprintf("%d/", prefix))
}

// KeyForwardDomain returns the key of record stores ibc-dns info
func KeyForwardDomain(name string) []byte {
	return append(
		KeyPrefixBytes(KeyForwardDomainPrefix),
		[]byte(name)...,
	)
}

// KeyReverseDomain returns the key of record stores domain name
func KeyReverseDomain(port, channel string) []byte {
	return append(
		KeyPrefixBytes(KeyReverseDomainPrefix),
		[]byte(fmt.Sprintf("%v/%v", port, channel))...,
	)
}

// KeyDomainAssociation returns the key of DomainAssociation
func KeyDomainAssociation(srcClientDomain, dstClientDomain types.ClientDomain) []byte {
	return append(
		KeyPrefixBytes(KeyDomainAssociationPrefix),
		[]byte(fmt.Sprintf("%v/%v/%v/%v", srcClientDomain.ClientId, srcClientDomain.DomainName, dstClientDomain.ClientId, dstClientDomain.DomainName))...,
	)
}
