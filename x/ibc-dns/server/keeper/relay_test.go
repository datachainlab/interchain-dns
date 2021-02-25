package keeper_test

import (
	"fmt"
	"testing"

	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	dnsservertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	ibctesting "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc/testing"
)

func TestDNSKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DNSKeeperTestSuite))
}

type DNSKeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	app0 *ibctesting.TestChain
	app1 *ibctesting.TestChain
	dns0 *ibctesting.TestChain

	chA0toA1 ibctesting.TestChannel
	chA1toA0 ibctesting.TestChannel

	chA0toD0 ibctesting.TestChannel
	chD0toA0 ibctesting.TestChannel

	chA1toD0 ibctesting.TestChannel
	chD0toA1 ibctesting.TestChannel
}

func (suite *DNSKeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 3)
	suite.dns0 = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.app0 = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.app1 = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func (suite *DNSKeeperTestSuite) TestDomainRegistration() {
	require := suite.Require()

	suite.openALLChannels(suite.app1.ChainID, suite.app0.ChainID)

	const app0Name = "domain-app0"

	//* app0: Domain registration *//

	suite.registerDomain(suite.app0, app0Name, suite.chA0toD0)

	//* Try to register a duplicated domain *//

	p0, err := suite.app1.App.DNSClientKeeper.SendPacketRegisterDomain(
		suite.app1.GetContext(),
		app0Name, // already used
		suite.chA1toD0.PortID,
		suite.chA1toD0.ID,
		[]byte("memo"),
	)
	require.NoError(err)
	data0 := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), p0.GetData()).(*dnsservertypes.RegisterDomainPacketData)
	require.Error(suite.dns0.App.DNSServerKeeper.ReceivePacketRegisterDomain(
		suite.dns0.GetContext(),
		*p0,
		data0,
	))
	require.NoError(suite.app1.App.DNSClientKeeper.ReceiveRegisterDomainPacketAcknowledgement(
		suite.app1.GetContext(),
		dnsservertypes.STATUS_FAILED,
		app0Name,
		*p0,
	))
	_, found := suite.app1.App.DNSClientKeeper.GetSelfDomainName(
		suite.app1.GetContext(),
		types.NewLocalDNSID(suite.chA1toD0.PortID, suite.chA1toD0.ID),
	)
	require.False(found)
}

func (suite *DNSKeeperTestSuite) TestDomainAssociation() {
	require := suite.Require()

	const (
		app0Name = "domain-app0"
		app1Name = "domain-app1"
	)

	suite.openALLChannels(suite.app1.ChainID, suite.app0.ChainID)
	suite.registerDomain(suite.app0, app0Name, suite.chA0toD0)
	suite.registerDomain(suite.app1, app1Name, suite.chA1toD0)

	/// case0: normal ///

	dnsID0 := types.NewLocalDNSID(suite.chA0toD0.PortID, suite.chA0toD0.ID)

	// app0: try to create a domain association
	{
		packet, err := suite.app0.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app0.GetContext(),
			dnsID0,
			types.NewClientDomain(app0Name, suite.app1.ChainID),
			types.NewClientDomain(app1Name, suite.app0.ChainID),
		)
		require.NoError(err)
		require.Equal(suite.chA0toD0.PortID, packet.GetSourcePort())
		require.Equal(suite.chA0toD0.ID, packet.GetSourceChannel())
		data := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData()).(*dnsservertypes.DomainAssociationCreatePacketData)
		ack, completed := suite.dns0.App.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.GetContext(),
			*packet,
			data,
		)
		require.False(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	dnsID1 := types.NewLocalDNSID(suite.chA1toD0.PortID, suite.chA1toD0.ID)
	// app1: try to confirm the domain association
	{
		packet, err := suite.app1.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app1.GetContext(),
			dnsID1,
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			types.NewClientDomain(app0Name, suite.app1.ChainID),
		)
		require.NoError(err)
		require.Equal(suite.chA1toD0.PortID, packet.GetSourcePort())
		require.Equal(suite.chA1toD0.ID, packet.GetSourceChannel())
		data := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData()).(*dnsservertypes.DomainAssociationCreatePacketData)
		ack, completed := suite.dns0.App.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.GetContext(),
			*packet,
			data,
		)
		require.True(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	// dns0: create a domain association
	{
		srcPacket, dstPacket, err := suite.dns0.App.DNSServerKeeper.CreateDomainAssociationResultPacketData(
			suite.dns0.GetContext(),
			servertypes.STATUS_OK,
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			types.NewClientDomain(app0Name, suite.app1.ChainID),
		)
		require.NoError(err)

		srcData := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), srcPacket.GetData()).(*dnsservertypes.DomainAssociationResultPacketData)
		require.Equal(servertypes.STATUS_OK, srcData.Status)
		require.Equal(suite.app1.ChainID, srcData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID0, app0Name), srcData.CounterpartyDomain)

		dstData := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), dstPacket.GetData()).(*dnsservertypes.DomainAssociationResultPacketData)
		require.Equal(servertypes.STATUS_OK, dstData.Status)
		require.Equal(suite.app0.ChainID, dstData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID1, app1Name), dstData.CounterpartyDomain)

		// receive the result of domain association
		require.NoError(
			suite.app0.App.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app0.GetContext(),
				*dstPacket,
				dstData,
			),
		)
		require.NoError(
			suite.app1.App.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app1.GetContext(),
				*srcPacket,
				srcData,
			),
		)

		// app0: get a local DNS-ID using DNS
		id1, found := suite.app0.App.DNSClientKeeper.ResolveDNSID(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
		)
		require.True(found)
		require.Equal(dnsID1, id1)

		// app1: get a local DNS-ID using DNS
		id0, found := suite.app1.App.DNSClientKeeper.ResolveDNSID(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
		)
		require.True(found)
		require.Equal(dnsID0, id0)

		// app0
		_, found = suite.app0.App.DNSClientKeeper.ResolveChannel(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.PortID,
		)
		require.False(found)
		require.NoError(suite.app0.App.DNSClientKeeper.SetDomainChannel(
			suite.app0.GetContext(),
			dnsID0,
			app1Name,
			types.NewChannel(suite.chA0toA1.PortID, suite.chA0toA1.ID, suite.chA1toA0.PortID, suite.chA1toA0.ID),
		))

		// app0: resolve a channel using DNS
		c0, found := suite.app0.App.DNSClientKeeper.ResolveChannel(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.PortID,
		)
		require.True(found)
		exc0, found := suite.app0.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app0.GetContext(), suite.chA0toA1.PortID, suite.chA0toA1.ID)
		require.True(found)
		require.Equal(exc0, c0)

		// app1
		_, found = suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.PortID,
		)
		require.False(found)
		require.NoError(suite.app1.App.DNSClientKeeper.SetDomainChannel(
			suite.app1.GetContext(),
			dnsID1,
			app0Name,
			types.NewChannel(suite.chA1toA0.PortID, suite.chA1toA0.ID, suite.chA0toA1.PortID, suite.chA0toA1.ID),
		))

		// app1: resolve a channel using DNS
		c1, found := suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.PortID,
		)
		require.True(found)
		exc1, found := suite.app1.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app1.GetContext(), suite.chA1toA0.PortID, suite.chA1toA0.ID)
		require.True(found)
		require.Equal(exc1, c1)

		res, err := suite.dns0.App.DNSServerKeeper.QueryDomains(suite.dns0.GetContext())
		require.NoError(err)
		require.Equal(2, len(res.Domains))
		require.Equal(app0Name, res.Domains[0].Name)
		require.Equal(app1Name, res.Domains[1].Name)
	}

	/// case1: A client referring to app1 in app0 is frozen, but DNS-ID is not changed ///

	var (
		app1ClientID = suite.app1.ChainID + "new"
	)

	// update channel info
	err := suite.createClient(suite.app0, suite.app1, app1ClientID, ibctesting.Tendermint)
	require.NoError(err)
	newConnA0toA1, newConnA1toA0 := suite.coordinator.CreateConnection(suite.app0, suite.app1, app1ClientID, suite.app0.ChainID)
	newChA0toA1, newChA1toA0 := suite.coordinator.CreateMockChannels(suite.app0, suite.app1, newConnA0toA1, newConnA1toA0, channeltypes.UNORDERED)

	// app0: try to create a domain association
	{
		packet, err := suite.app0.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app0.GetContext(),
			dnsID0,
			types.NewClientDomain(app0Name, app1ClientID),
			types.NewClientDomain(app1Name, suite.app0.ChainID),
		)
		require.NoError(err)
		require.Equal(suite.chA0toD0.PortID, packet.GetSourcePort())
		require.Equal(suite.chA0toD0.ID, packet.GetSourceChannel())
		data := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData()).(*dnsservertypes.DomainAssociationCreatePacketData)
		ack, completed := suite.dns0.App.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.GetContext(),
			*packet,
			data,
		)
		require.False(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	// app1: try to confirm the domain association
	{
		packet, err := suite.app1.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app1.GetContext(),
			dnsID1,
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			types.NewClientDomain(app0Name, app1ClientID),
		)
		require.NoError(err)
		require.Equal(suite.chA1toD0.PortID, packet.GetSourcePort())
		require.Equal(suite.chA1toD0.ID, packet.GetSourceChannel())
		data := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData()).(*dnsservertypes.DomainAssociationCreatePacketData)
		ack, completed := suite.dns0.App.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.GetContext(),
			*packet,
			data,
		)
		require.True(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	// dns0: create a domain association
	{
		srcPacket, dstPacket, err := suite.dns0.App.DNSServerKeeper.CreateDomainAssociationResultPacketData(
			suite.dns0.GetContext(),
			servertypes.STATUS_OK,
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			types.NewClientDomain(app0Name, app1ClientID),
		)
		require.NoError(err)

		srcData := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), srcPacket.GetData()).(*dnsservertypes.DomainAssociationResultPacketData)
		require.Equal(servertypes.STATUS_OK, srcData.Status)
		require.Equal(app1ClientID, srcData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID0, app0Name), srcData.CounterpartyDomain)

		dstData := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), dstPacket.GetData()).(*dnsservertypes.DomainAssociationResultPacketData)
		require.Equal(servertypes.STATUS_OK, dstData.Status)
		require.Equal(suite.app0.ChainID, dstData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID1, app1Name), dstData.CounterpartyDomain)

		// receive the result of domain association
		require.NoError(
			suite.app0.App.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app0.GetContext(),
				*dstPacket,
				dstData,
			),
		)
		require.NoError(
			suite.app1.App.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app1.GetContext(),
				*srcPacket,
				srcData,
			),
		)

		// app0: get a local DNS-ID using DNS
		id1, found := suite.app0.App.DNSClientKeeper.ResolveDNSID(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
		)
		require.True(found)
		require.Equal(dnsID1, id1)

		// app1: get a local DNS-ID using DNS
		id0, found := suite.app1.App.DNSClientKeeper.ResolveDNSID(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
		)
		require.True(found)
		require.Equal(dnsID0, id0)

		// app0
		// resolve the domain name to the old channel
		_, found = suite.app0.App.DNSClientKeeper.ResolveChannel(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
			newChA0toA1.PortID,
		)
		require.True(found)
		require.NoError(suite.app0.App.DNSClientKeeper.SetDomainChannel(
			suite.app0.GetContext(),
			dnsID0,
			app1Name,
			types.NewChannel(newChA0toA1.PortID, newChA0toA1.ID, newChA1toA0.PortID, newChA1toA0.ID),
		))

		// app0: resolve a channel using DNS
		c0, found := suite.app0.App.DNSClientKeeper.ResolveChannel(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
			newChA0toA1.PortID,
		)
		require.True(found)
		exc0, found := suite.app0.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app0.GetContext(), newChA0toA1.PortID, newChA0toA1.ID)
		require.True(found)
		require.Equal(exc0, c0)

		// app1
		// resolve the domain name to the old channel
		_, found = suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			newChA1toA0.PortID,
		)
		require.True(found)
		require.NoError(suite.app1.App.DNSClientKeeper.SetDomainChannel(
			suite.app1.GetContext(),
			dnsID1,
			app0Name,
			types.NewChannel(newChA1toA0.PortID, newChA1toA0.ID, newChA0toA1.PortID, newChA0toA1.ID),
		))

		// app1: resolve a channel using DNS
		c1, found := suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			newChA1toA0.PortID,
		)
		require.True(found)
		exc1, found := suite.app1.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app1.GetContext(), newChA1toA0.PortID, newChA1toA0.ID)
		require.True(found)
		require.Equal(exc1, c1)
	}
}

func (suite *DNSKeeperTestSuite) registerDomain(
	chain *ibctesting.TestChain,
	name string,
	srcch ibctesting.TestChannel,
) {
	require := suite.Require()

	p0, err := chain.App.DNSClientKeeper.SendPacketRegisterDomain(
		chain.GetContext(),
		name,
		srcch.PortID,
		srcch.ID,
		[]byte("memo"),
	)
	require.NoError(err)
	data0 := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), p0.GetData()).(*dnsservertypes.RegisterDomainPacketData)
	require.NoError(suite.dns0.App.DNSServerKeeper.ReceivePacketRegisterDomain(
		suite.dns0.GetContext(),
		*p0,
		data0,
	))
	require.NoError(chain.App.DNSClientKeeper.ReceiveRegisterDomainPacketAcknowledgement(
		chain.GetContext(),
		dnsservertypes.STATUS_OK,
		name,
		*p0,
	))

	dnsID := types.NewLocalDNSID(srcch.PortID, srcch.ID)

	res, err := suite.dns0.App.DNSServerKeeper.QueryDomain(suite.dns0.GetContext(), servertypes.QueryDomainRequest{Name: name})
	require.NoError(err)
	require.Equal(dnsID, res.Domain.DnsId)

	name, found := chain.App.DNSClientKeeper.GetSelfDomainName(
		chain.GetContext(),
		dnsID,
	)
	require.True(found)
	require.Equal(name, name)
}

func (suite *DNSKeeperTestSuite) setupClients(
	chainA, chainB *ibctesting.TestChain,
	chainAClientID, chainBClientID string,
	clientType string,
) (string, string) {

	err := suite.createClient(chainA, chainB, chainAClientID, clientType)
	require.NoError(suite.T(), err)

	err = suite.createClient(chainB, chainA, chainBClientID, clientType)
	require.NoError(suite.T(), err)

	return chainAClientID, chainBClientID
}

func (suite *DNSKeeperTestSuite) setupClientConnections(
	chainA, chainB *ibctesting.TestChain,
	chainAClientID, chainBClientID string,
	clientType string,
) (string, string, *ibctesting.TestConnection, *ibctesting.TestConnection) {

	clientA, clientB := suite.setupClients(chainA, chainB, chainAClientID, chainBClientID, clientType)

	connA, connB := suite.coordinator.CreateConnection(chainA, chainB, clientA, clientB)

	return clientA, clientB, connA, connB
}

func (suite *DNSKeeperTestSuite) createClient(
	source, counterparty *ibctesting.TestChain,
	clientID string,
	clientType string,
) (err error) {
	suite.coordinator.CommitBlock(source, counterparty)

	switch clientType {
	case ibctesting.Tendermint:
		err = source.CreateTMClient(counterparty, clientID)

	default:
		err = fmt.Errorf("client type %s is not supported", clientType)
	}

	if err != nil {
		return err
	}

	suite.coordinator.IncrementTime()

	return nil
}

func (suite *DNSKeeperTestSuite) openALLChannels(srcClientID, dstClientID string) {
	suite.chA0toA1, suite.chA1toA0 = suite.openAppChannels(srcClientID, dstClientID)

	_, _, connA0toD0, connD0toA0 := suite.setupClientConnections(
		suite.app0,
		suite.dns0,
		suite.dns0.ChainID,
		dstClientID,
		ibctesting.Tendermint,
	)
	suite.dns0.CreatePortCapability(types.PortID)
	suite.chA0toD0, suite.chD0toA0 = suite.coordinator.CreateChannel(
		suite.app0,
		suite.dns0,
		connA0toD0,
		connD0toA0,
		types.PortID,
		types.PortID,
		channeltypes.UNORDERED,
	)

	_, _, connA1toD0, connD0toA1 := suite.setupClientConnections(
		suite.app1,
		suite.dns0,
		suite.dns0.ChainID,
		srcClientID,
		ibctesting.Tendermint,
	)
	suite.chA1toD0, suite.chD0toA1 = suite.coordinator.CreateChannel(
		suite.app1,
		suite.dns0,
		connA1toD0,
		connD0toA1,
		types.PortID,
		types.PortID,
		channeltypes.UNORDERED,
	)
}

func (suite *DNSKeeperTestSuite) openAppChannels(srcClientID, dstClientID string) (ibctesting.TestChannel, ibctesting.TestChannel) {
	_, _, connA0toA1, connA1toA0 := suite.setupClientConnections(
		suite.app0,
		suite.app1,
		srcClientID,
		dstClientID,
		ibctesting.Tendermint,
	)
	return suite.coordinator.CreateMockChannels(suite.app0, suite.app1, connA0toA1, connA1toA0, channeltypes.UNORDERED)
}
