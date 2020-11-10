package keeper_test

import (
	"testing"

	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/suite"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	dnsservertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	ibctesting "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/testing"
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
		suite.chA1toD0.Port,
		suite.chA1toD0.Channel,
		[]byte("memo"),
		suite.app1.GetTimeoutHeight(suite.dns0.ChainID, 1),
		0,
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
		types.NewLocalDNSID(suite.chA1toD0.Port, suite.chA1toD0.Channel),
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

	dnsID0 := types.NewLocalDNSID(suite.chA0toD0.Port, suite.chA0toD0.Channel)

	// app0: try to create a domain association
	{
		packet, err := suite.app0.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app0.GetContext(),
			dnsID0,
			types.NewClientDomain(app0Name, suite.app1.ChainID),
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			suite.app0.GetTimeoutHeight(suite.dns0.ChainID, 1),
			0,
		)
		require.NoError(err)
		require.Equal(suite.chA0toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA0toD0.Channel, packet.GetSourceChannel())
		data := types.MustDeserializeJSONPacketData(servertypes.PacketCdc(), packet.GetData()).(*dnsservertypes.DomainAssociationCreatePacketData)
		ack, completed := suite.dns0.App.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.GetContext(),
			*packet,
			data,
		)
		require.False(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	dnsID1 := types.NewLocalDNSID(suite.chA1toD0.Port, suite.chA1toD0.Channel)
	// app1: try to confirm the domain association
	{
		packet, err := suite.app1.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app1.GetContext(),
			dnsID1,
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			types.NewClientDomain(app0Name, suite.app1.ChainID),
			suite.app1.GetTimeoutHeight(suite.dns0.ChainID, 1),
			0,
		)
		require.NoError(err)
		require.Equal(suite.chA1toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA1toD0.Channel, packet.GetSourceChannel())
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
			suite.dns0.GetTimeoutHeight(suite.app0.ChainID, 1),
			0,
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
			suite.chA0toA1.Port,
		)
		require.False(found)
		require.NoError(suite.app0.App.DNSClientKeeper.SetDomainChannel(
			suite.app0.GetContext(),
			dnsID0,
			app1Name,
			types.NewChannel(suite.chA0toA1.Port, suite.chA0toA1.Channel, suite.chA1toA0.Port, suite.chA1toA0.Channel),
		))

		// app0: resolve a channel using DNS
		c0, found := suite.app0.App.DNSClientKeeper.ResolveChannel(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.Port,
		)
		require.True(found)
		exc0, found := suite.app0.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app0.GetContext(), suite.chA0toA1.Port, suite.chA0toA1.Channel)
		require.True(found)
		require.Equal(exc0, c0)

		// app1
		_, found = suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.False(found)
		require.NoError(suite.app1.App.DNSClientKeeper.SetDomainChannel(
			suite.app1.GetContext(),
			dnsID1,
			app0Name,
			types.NewChannel(suite.chA1toA0.Port, suite.chA1toA0.Channel, suite.chA0toA1.Port, suite.chA0toA1.Channel),
		))

		// app1: resolve a channel using DNS
		c1, found := suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.True(found)
		exc1, found := suite.app1.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app1.GetContext(), suite.chA1toA0.Port, suite.chA1toA0.Channel)
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
	err := suite.coordinator.CreateClient(suite.app0, suite.app1, app1ClientID, ibctesting.Tendermint)
	require.NoError(err)
	connA0toA1, connA1toA0 := suite.coordinator.CreateConnection(suite.app0, suite.app1, app1ClientID, suite.app0.ChainID)
	suite.chA0toA1, suite.chA1toA0 = suite.coordinator.CreateMockChannels(suite.app0, suite.app1, connA0toA1, connA1toA0, channeltypes.UNORDERED)

	// app0: try to create a domain association
	{
		packet, err := suite.app0.App.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app0.GetContext(),
			dnsID0,
			types.NewClientDomain(app0Name, app1ClientID),
			types.NewClientDomain(app1Name, suite.app0.ChainID),
			suite.app0.GetTimeoutHeight(suite.dns0.ChainID, 1),
			0,
		)
		require.NoError(err)
		require.Equal(suite.chA0toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA0toD0.Channel, packet.GetSourceChannel())
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
			suite.app1.GetTimeoutHeight(suite.dns0.ChainID, 1),
			0,
		)
		require.NoError(err)
		require.Equal(suite.chA1toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA1toD0.Channel, packet.GetSourceChannel())
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
			suite.dns0.GetTimeoutHeight(suite.app0.ChainID, 1),
			0,
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
			suite.chA0toA1.Port,
		)
		require.True(found)
		require.NoError(suite.app0.App.DNSClientKeeper.SetDomainChannel(
			suite.app0.GetContext(),
			dnsID0,
			app1Name,
			types.NewChannel(suite.chA0toA1.Port, suite.chA0toA1.Channel, suite.chA1toA0.Port, suite.chA1toA0.Channel),
		))

		// app0: resolve a channel using DNS
		c0, found := suite.app0.App.DNSClientKeeper.ResolveChannel(
			suite.app0.GetContext(),
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.Port,
		)
		require.True(found)
		exc0, found := suite.app0.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app0.GetContext(), suite.chA0toA1.Port, suite.chA0toA1.Channel)
		require.True(found)
		require.Equal(exc0, c0)

		// app1
		// resolve the domain name to the old channel
		_, found = suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.True(found)
		require.NoError(suite.app1.App.DNSClientKeeper.SetDomainChannel(
			suite.app1.GetContext(),
			dnsID1,
			app0Name,
			types.NewChannel(suite.chA1toA0.Port, suite.chA1toA0.Channel, suite.chA0toA1.Port, suite.chA0toA1.Channel),
		))

		// app1: resolve a channel using DNS
		c1, found := suite.app1.App.DNSClientKeeper.ResolveChannel(
			suite.app1.GetContext(),
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.True(found)
		exc1, found := suite.app1.App.IBCKeeper.ChannelKeeper.GetChannel(suite.app1.GetContext(), suite.chA1toA0.Port, suite.chA1toA0.Channel)
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
		srcch.Port,
		srcch.Channel,
		[]byte("memo"),
		chain.GetTimeoutHeight(suite.dns0.ChainID, 1),
		0,
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

	dnsID := types.NewLocalDNSID(srcch.Port, srcch.Channel)

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

func (suite *DNSKeeperTestSuite) openALLChannels(srcClientID, dstClientID string) {
	suite.chA0toA1, suite.chA1toA0 = suite.openAppChannels(srcClientID, dstClientID)

	_, _, connA0toD0, connD0toA0 := suite.coordinator.SetupClientConnections(
		suite.app0,
		suite.dns0,
		suite.dns0.ChainID,
		dstClientID,
		ibctesting.Tendermint,
	)
	suite.chA0toD0, suite.chD0toA0 = suite.coordinator.CreateDNSChannels(
		suite.app0,
		suite.dns0,
		connA0toD0,
		connD0toA0,
		channeltypes.UNORDERED,
	)

	_, _, connA1toD0, connD0toA1 := suite.coordinator.SetupClientConnections(
		suite.app1,
		suite.dns0,
		suite.dns0.ChainID,
		srcClientID,
		ibctesting.Tendermint,
	)
	suite.chA1toD0, suite.chD0toA1 = suite.coordinator.CreateDNSChannels(
		suite.app1,
		suite.dns0,
		connA1toD0,
		connD0toA1,
		channeltypes.UNORDERED,
	)
}

func (suite *DNSKeeperTestSuite) openAppChannels(srcClientID, dstClientID string) (ibctesting.TestChannel, ibctesting.TestChannel) {
	_, _, connA0toA1, connA1toA0 := suite.coordinator.SetupClientConnections(
		suite.app0,
		suite.app1,
		srcClientID,
		dstClientID,
		ibctesting.Tendermint,
	)
	return suite.coordinator.CreateMockChannels(suite.app0, suite.app1, connA0toA1, connA1toA0, channeltypes.UNORDERED)
}
