package types

import (
	"fmt"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
)

const (
	// StoreKey to be used when creating the KVStore
	StoreKey = "ibc-dns-client"
)

const (
	KeySelfDomainPrefix uint8 = iota
	KeyLocalDNSIDPrefix
	KeyClientDomainPrefix
	KeyDomainChannelPrefix
)

// KeyPrefixBytes return the key prefix bytes from a URL string format
func KeyPrefixBytes(prefix uint8) []byte {
	return []byte(fmt.Sprintf("%d/", prefix))
}

func KeyLocalDNSID(dnsID types.LocalDNSID, domain string) []byte {
	return append(
		KeyPrefixBytes(KeyLocalDNSIDPrefix),
		[]byte(fmt.Sprintf("%v/%v/%v", dnsID.SourcePort, dnsID.SourceChannel, domain))...,
	)
}

func KeySelfDomain(dnsID types.LocalDNSID) []byte {
	return append(
		KeyPrefixBytes(KeySelfDomainPrefix),
		[]byte(fmt.Sprintf("%v/%v", dnsID.SourcePort, dnsID.SourceChannel))...,
	)
}

// clientID => domain name
func KeyClientDomain(dnsID types.LocalDNSID, domainName string) []byte {
	return append(
		KeyPrefixBytes(KeyClientDomainPrefix),
		[]byte(fmt.Sprintf("%v/%v/%v", dnsID.SourcePort, dnsID.SourceChannel, domainName))...,
	)
}

// domain:port => channel
func KeyDomainChannel(dnsID types.LocalDNSID, domainName, port string) []byte {
	return append(
		KeyPrefixBytes(KeyDomainChannelPrefix),
		[]byte(fmt.Sprintf("%v/%v/%v/%v", dnsID.SourcePort, dnsID.SourceChannel, domainName, port))...,
	)
}
