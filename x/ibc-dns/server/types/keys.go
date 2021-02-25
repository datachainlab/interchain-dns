package types

import (
	"fmt"

	"github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
)

const (
	// StoreKey to be used when creating the KVStore
	StoreKey = "ibc-dns-server"
)

const (
	KeyForwardDomainPrefix uint8 = iota
	KeyReverseDomainPrefix
	KeyDomainMappingPrefix
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

// KeyDomainMapping returns the key of DomainMapping
func KeyDomainMapping(srcClientDomain, dstClientDomain types.ClientDomain) []byte {
	return append(
		KeyPrefixBytes(KeyDomainMappingPrefix),
		[]byte(fmt.Sprintf("%v/%v/%v/%v", srcClientDomain.ClientId, srcClientDomain.DomainName, dstClientDomain.ClientId, dstClientDomain.DomainName))...,
	)
}
