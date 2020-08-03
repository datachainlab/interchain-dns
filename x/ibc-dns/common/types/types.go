package types

func NewLocalDNSID(srcPort, srcChannel string) LocalDNSID {
	return LocalDNSID{SourcePort: srcPort, SourceChannel: srcChannel}
}

func NewChannel(srcPort, srcChannel, dstPort, dstChannel string) LocalChannel {
	return LocalChannel{SourcePort: srcPort, SourceChannel: srcChannel, DestinationPort: dstPort, DestinationChannel: dstChannel}
}

// NewLocalDomain returns a new LocalDomain
func NewLocalDomain(dnsID LocalDNSID, name string) LocalDomain {
	return LocalDomain{
		Name:  name,
		DNSID: dnsID,
	}
}

func NewClientDomain(domainName, clientID string) ClientDomain {
	return ClientDomain{DomainName: domainName, ClientId: clientID}
}
