package keeper_test

import (
	"testing"

	dns "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	dnsservertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestDNSKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DNSKeeperTestSuite))
}

type DNSKeeperTestSuite struct {
	KeeperTestSuite

	app0 *appContext
	app1 *appContext
	dns0 *appContext

	chA0toA1 ChannelInfo
	chA1toA0 ChannelInfo

	chA0toD0 ChannelInfo
	chD0toA0 ChannelInfo

	chA1toD0 ChannelInfo
	chD0toA1 ChannelInfo
}

func (suite *DNSKeeperTestSuite) SetupTest() {
	suite.dns0 = suite.createAppWithHeader(abci.Header{ChainID: "dns0"})
	suite.app0 = suite.createAppWithHeader(abci.Header{ChainID: "app0"})
	suite.app1 = suite.createAppWithHeader(abci.Header{ChainID: "app1"})

	suite.chA0toA1 = ChannelInfo{"testportzeroone", "testchannelzeroone"} // app0 -> app1
	suite.chA1toA0 = ChannelInfo{"testportonezero", "testchannelonezero"} // app1 -> app0

	suite.chA0toD0 = ChannelInfo{dns.PortID, "testchannelzerodns"} // app0 -> dns0
	suite.chD0toA0 = ChannelInfo{dns.PortID, "testchanneldnszero"} // dns0 -> app0

	suite.chA1toD0 = ChannelInfo{dns.PortID, "testchannelonedns"} // app1 -> dns0
	suite.chD0toA1 = ChannelInfo{dns.PortID, "testchanneldnsone"} // dns0 -> app1
}

func (suite *DNSKeeperTestSuite) TestDomainRegistration() {
	require := suite.Require()

	suite.openALLChannels(suite.app1.chainID, suite.app0.chainID)

	const app0Name = "domain-app0"

	//* app0: Domain registration *//

	suite.registerDomain(suite.app0, app0Name, suite.chA0toD0, suite.chD0toA0)

	//* Try to register a duplicated domain *//

	p0, err := suite.app1.app.DNSClientKeeper.SendPacketRegisterDomain(
		suite.app1.ctx,
		app0Name, // already used
		suite.chA1toD0.Port,
		suite.chA1toD0.Channel,
		[]byte("memo"),
	)
	require.NoError(err)
	var data0 dnsservertypes.RegisterDomainPacketData
	require.NoError(servertypes.ModuleCdc.UnmarshalJSON(p0.Data, &data0))
	require.Error(suite.dns0.app.DNSServerKeeper.ReceivePacketRegisterDomain(
		suite.dns0.ctx,
		*p0,
		data0,
	))
	require.NoError(suite.app1.app.DNSClientKeeper.ReceiveRegisterDomainPacketAcknowledgement(
		suite.app1.ctx,
		dnsservertypes.STATUS_FAILED,
		app0Name,
		*p0,
	))
	_, found := suite.app1.app.DNSClientKeeper.GetSelfDomainName(
		suite.app1.ctx,
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

	suite.openALLChannels(suite.app1.chainID, suite.app0.chainID)
	suite.registerDomain(suite.app0, app0Name, suite.chA0toD0, suite.chD0toA0)
	suite.registerDomain(suite.app1, app1Name, suite.chA1toD0, suite.chD0toA1)

	/// case0: normal ///

	dnsID0 := types.NewLocalDNSID(suite.chA0toD0.Port, suite.chA0toD0.Channel)

	// app0: try to create a domain association
	{
		packet, err := suite.app0.app.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app0.ctx,
			dnsID0,
			types.NewClientDomain(app0Name, suite.app1.chainID),
			types.NewClientDomain(app1Name, suite.app0.chainID),
		)
		require.NoError(err)
		require.Equal(suite.chA0toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA0toD0.Channel, packet.GetSourceChannel())
		var data dnsservertypes.DomainAssociationCreatePacketData
		require.NoError(servertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
		ack, completed := suite.dns0.app.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.ctx,
			*packet,
			data,
		)
		require.False(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	dnsID1 := types.NewLocalDNSID(suite.chA1toD0.Port, suite.chA1toD0.Channel)
	// app1: try to confirm the domain association
	{
		packet, err := suite.app1.app.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app1.ctx,
			dnsID1,
			types.NewClientDomain(app1Name, suite.app0.chainID),
			types.NewClientDomain(app0Name, suite.app1.chainID),
		)
		require.NoError(err)
		require.Equal(suite.chA1toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA1toD0.Channel, packet.GetSourceChannel())
		var data dnsservertypes.DomainAssociationCreatePacketData
		require.NoError(servertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
		ack, completed := suite.dns0.app.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.ctx,
			*packet,
			data,
		)
		require.True(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	// dns0: create a domain association
	{
		srcPacket, dstPacket, err := suite.dns0.app.DNSServerKeeper.CreateDomainAssociationResultPacketData(
			suite.dns0.ctx,
			servertypes.STATUS_OK,
			types.NewClientDomain(app1Name, suite.app0.chainID),
			types.NewClientDomain(app0Name, suite.app1.chainID),
		)
		require.NoError(err)

		var srcData, dstData servertypes.DomainAssociationResultPacketData
		servertypes.ModuleCdc.MustUnmarshalJSON(srcPacket.Data, &srcData)
		require.Equal(servertypes.STATUS_OK, srcData.Status)
		require.Equal(suite.app1.chainID, srcData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID0, app0Name), srcData.CounterpartyDomain)

		servertypes.ModuleCdc.MustUnmarshalJSON(dstPacket.Data, &dstData)
		require.Equal(servertypes.STATUS_OK, dstData.Status)
		require.Equal(suite.app0.chainID, dstData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID1, app1Name), dstData.CounterpartyDomain)

		// receive the result of domain association
		require.NoError(
			suite.app0.app.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app0.ctx,
				*dstPacket,
				dstData,
			),
		)
		require.NoError(
			suite.app1.app.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app1.ctx,
				*srcPacket,
				srcData,
			),
		)

		// app0: get a local DNS-ID using DNS
		id1, found := suite.app0.app.DNSClientKeeper.ResolveDNSID(
			suite.app0.ctx,
			types.NewLocalDomain(dnsID0, app1Name),
		)
		require.True(found)
		require.Equal(dnsID1, id1)

		// app1: get a local DNS-ID using DNS
		id0, found := suite.app1.app.DNSClientKeeper.ResolveDNSID(
			suite.app1.ctx,
			types.NewLocalDomain(dnsID1, app0Name),
		)
		require.True(found)
		require.Equal(dnsID0, id0)

		// app0
		_, found = suite.app0.app.DNSClientKeeper.ResolveChannel(
			suite.app0.ctx,
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.Port,
		)
		require.False(found)
		require.NoError(suite.app0.app.DNSClientKeeper.SetDomainChannel(
			suite.app0.ctx,
			dnsID0,
			app1Name,
			types.NewChannel(suite.chA0toA1.Port, suite.chA0toA1.Channel, suite.chA1toA0.Port, suite.chA1toA0.Channel),
		))

		// app0: resolve a channel using DNS
		c0, found := suite.app0.app.DNSClientKeeper.ResolveChannel(
			suite.app0.ctx,
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.Port,
		)
		require.True(found)
		exc0, found := suite.app0.app.IBCKeeper.ChannelKeeper.GetChannel(suite.app0.ctx, suite.chA0toA1.Port, suite.chA0toA1.Channel)
		require.True(found)
		require.Equal(exc0, c0)

		// app1
		_, found = suite.app1.app.DNSClientKeeper.ResolveChannel(
			suite.app1.ctx,
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.False(found)
		require.NoError(suite.app1.app.DNSClientKeeper.SetDomainChannel(
			suite.app1.ctx,
			dnsID1,
			app0Name,
			types.NewChannel(suite.chA1toA0.Port, suite.chA1toA0.Channel, suite.chA0toA1.Port, suite.chA0toA1.Channel),
		))

		// app1: resolve a channel using DNS
		c1, found := suite.app1.app.DNSClientKeeper.ResolveChannel(
			suite.app1.ctx,
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.True(found)
		exc1, found := suite.app1.app.IBCKeeper.ChannelKeeper.GetChannel(suite.app1.ctx, suite.chA1toA0.Port, suite.chA1toA0.Channel)
		require.True(found)
		require.Equal(exc1, c1)
	}

	/// case1: A client referring to app1 in app0 is frozen, but DNS-ID is not changed ///

	var (
		app1ClientID = suite.app1.chainID + "new"
	)

	// update channel info
	suite.chA0toA1 = ChannelInfo{suite.chA0toA1.Port, "testchannelzeroone" + "new"} // app0 -> app1
	suite.chA1toA0 = ChannelInfo{suite.chA1toA0.Port, "testchannelonezero" + "new"} // app1 -> app0
	suite.openAppChannels(app1ClientID, suite.app0.chainID, true)

	// app0: try to create a domain association
	{
		packet, err := suite.app0.app.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app0.ctx,
			dnsID0,
			types.NewClientDomain(app0Name, app1ClientID),
			types.NewClientDomain(app1Name, suite.app0.chainID),
		)
		require.NoError(err)
		require.Equal(suite.chA0toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA0toD0.Channel, packet.GetSourceChannel())
		var data dnsservertypes.DomainAssociationCreatePacketData
		require.NoError(servertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
		ack, completed := suite.dns0.app.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.ctx,
			*packet,
			data,
		)
		require.False(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	// app1: try to confirm the domain association
	{
		packet, err := suite.app1.app.DNSClientKeeper.SendDomainAssociationCreatePacketData(
			suite.app1.ctx,
			dnsID1,
			types.NewClientDomain(app1Name, suite.app0.chainID),
			types.NewClientDomain(app0Name, app1ClientID),
		)
		require.NoError(err)
		require.Equal(suite.chA1toD0.Port, packet.GetSourcePort())
		require.Equal(suite.chA1toD0.Channel, packet.GetSourceChannel())
		var data dnsservertypes.DomainAssociationCreatePacketData
		require.NoError(servertypes.ModuleCdc.UnmarshalJSON(packet.Data, &data))
		ack, completed := suite.dns0.app.DNSServerKeeper.ReceiveDomainAssociationCreatePacketData(
			suite.dns0.ctx,
			*packet,
			data,
		)
		require.True(completed)
		require.Equal(ack.Status, servertypes.STATUS_OK)
	}

	// dns0: create a domain association
	{
		srcPacket, dstPacket, err := suite.dns0.app.DNSServerKeeper.CreateDomainAssociationResultPacketData(
			suite.dns0.ctx,
			servertypes.STATUS_OK,
			types.NewClientDomain(app1Name, suite.app0.chainID),
			types.NewClientDomain(app0Name, app1ClientID),
		)
		require.NoError(err)

		var srcData, dstData servertypes.DomainAssociationResultPacketData
		servertypes.ModuleCdc.MustUnmarshalJSON(srcPacket.Data, &srcData)
		require.Equal(servertypes.STATUS_OK, srcData.Status)
		require.Equal(app1ClientID, srcData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID0, app0Name), srcData.CounterpartyDomain)

		servertypes.ModuleCdc.MustUnmarshalJSON(dstPacket.Data, &dstData)
		require.Equal(servertypes.STATUS_OK, dstData.Status)
		require.Equal(suite.app0.chainID, dstData.ClientId)
		require.Equal(types.NewLocalDomain(dnsID1, app1Name), dstData.CounterpartyDomain)

		// receive the result of domain association
		require.NoError(
			suite.app0.app.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app0.ctx,
				*dstPacket,
				dstData,
			),
		)
		require.NoError(
			suite.app1.app.DNSClientKeeper.ReceiveDomainAssociationResultPacketData(
				suite.app1.ctx,
				*srcPacket,
				srcData,
			),
		)

		// app0: get a local DNS-ID using DNS
		id1, found := suite.app0.app.DNSClientKeeper.ResolveDNSID(
			suite.app0.ctx,
			types.NewLocalDomain(dnsID0, app1Name),
		)
		require.True(found)
		require.Equal(dnsID1, id1)

		// app1: get a local DNS-ID using DNS
		id0, found := suite.app1.app.DNSClientKeeper.ResolveDNSID(
			suite.app1.ctx,
			types.NewLocalDomain(dnsID1, app0Name),
		)
		require.True(found)
		require.Equal(dnsID0, id0)

		// app0
		// resolve the domain name to the old channel
		_, found = suite.app0.app.DNSClientKeeper.ResolveChannel(
			suite.app0.ctx,
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.Port,
		)
		require.True(found)
		require.NoError(suite.app0.app.DNSClientKeeper.SetDomainChannel(
			suite.app0.ctx,
			dnsID0,
			app1Name,
			types.NewChannel(suite.chA0toA1.Port, suite.chA0toA1.Channel, suite.chA1toA0.Port, suite.chA1toA0.Channel),
		))

		// app0: resolve a channel using DNS
		c0, found := suite.app0.app.DNSClientKeeper.ResolveChannel(
			suite.app0.ctx,
			types.NewLocalDomain(dnsID0, app1Name),
			suite.chA0toA1.Port,
		)
		require.True(found)
		exc0, found := suite.app0.app.IBCKeeper.ChannelKeeper.GetChannel(suite.app0.ctx, suite.chA0toA1.Port, suite.chA0toA1.Channel)
		require.True(found)
		require.Equal(exc0, c0)

		// app1
		// resolve the domain name to the old channel
		_, found = suite.app1.app.DNSClientKeeper.ResolveChannel(
			suite.app1.ctx,
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.True(found)
		require.NoError(suite.app1.app.DNSClientKeeper.SetDomainChannel(
			suite.app1.ctx,
			dnsID1,
			app0Name,
			types.NewChannel(suite.chA1toA0.Port, suite.chA1toA0.Channel, suite.chA0toA1.Port, suite.chA0toA1.Channel),
		))

		// app1: resolve a channel using DNS
		c1, found := suite.app1.app.DNSClientKeeper.ResolveChannel(
			suite.app1.ctx,
			types.NewLocalDomain(dnsID1, app0Name),
			suite.chA1toA0.Port,
		)
		require.True(found)
		exc1, found := suite.app1.app.IBCKeeper.ChannelKeeper.GetChannel(suite.app1.ctx, suite.chA1toA0.Port, suite.chA1toA0.Channel)
		require.True(found)
		require.Equal(exc1, c1)
	}
}

func (suite *DNSKeeperTestSuite) registerDomain(app *appContext, name string, srcci, dstci ChannelInfo) {
	require := suite.Require()

	p0, err := app.app.DNSClientKeeper.SendPacketRegisterDomain(
		app.ctx,
		name,
		srcci.Port,
		srcci.Channel,
		[]byte("memo"),
	)
	require.NoError(err)
	var data0 dnsservertypes.RegisterDomainPacketData
	require.NoError(servertypes.ModuleCdc.UnmarshalJSON(p0.Data, &data0))
	require.NoError(suite.dns0.app.DNSServerKeeper.ReceivePacketRegisterDomain(
		suite.dns0.ctx,
		*p0,
		data0,
	))
	require.NoError(app.app.DNSClientKeeper.ReceiveRegisterDomainPacketAcknowledgement(
		app.ctx,
		dnsservertypes.STATUS_OK,
		name,
		*p0,
	))
	name, found := app.app.DNSClientKeeper.GetSelfDomainName(
		app.ctx,
		types.NewLocalDNSID(srcci.Port, srcci.Channel),
	)
	require.True(found)
	require.Equal(name, name)
}

func (suite *DNSKeeperTestSuite) openALLChannels(srcClientID, dstClientID string) {
	suite.openAppChannels(srcClientID, dstClientID, false)

	suite.openChannels(
		suite.dns0.chainID,
		dstClientID+suite.dns0.chainID,
		suite.chA0toD0,
		suite.app0,

		dstClientID,
		suite.dns0.chainID+dstClientID,
		suite.chD0toA0,
		suite.dns0,

		false,
	)

	suite.openChannels(
		suite.dns0.chainID,
		srcClientID+suite.dns0.chainID,
		suite.chA1toD0,
		suite.app1,

		srcClientID,
		suite.dns0.chainID+srcClientID,
		suite.chD0toA1,
		suite.dns0,

		false,
	)
}

func (suite *DNSKeeperTestSuite) openAppChannels(srcClientID, dstClientID string, skipIfClientExists bool) {
	suite.openChannels(
		srcClientID,
		dstClientID+srcClientID,
		suite.chA0toA1,
		suite.app0,

		dstClientID,
		srcClientID+dstClientID,
		suite.chA1toA0,
		suite.app1,

		skipIfClientExists,
	)
}
