package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/libs/log"

	clienttypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// Keeper defines ibc-dns keeper
type Keeper struct {
	storeKey         sdk.StoreKey
	clientKeeper     types.ClientKeeper
	connectionKeeper types.ConnectionKeeper
	channelKeeper    types.ChannelKeeper
	scopedKeeper     capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new ibc-dns Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	clientKeeper types.ClientKeeper,
	connectionKeeper types.ConnectionKeeper,
	channelKeeper types.ChannelKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	return Keeper{
		storeKey:         storeKey,
		clientKeeper:     clientKeeper,
		connectionKeeper: connectionKeeper,
		channelKeeper:    channelKeeper,
		scopedKeeper:     scopedKeeper,
	}
}

// SendPacketRegisterDomain sends a packet to register a channel name on DNS server
func (k Keeper) SendPacketRegisterDomain(ctx sdk.Context, name string, sourcePort string, sourceChannel string, metadata []byte) (*channel.Packet, error) {
	data := servertypes.NewRegisterDomainPacketData(name, metadata)
	c, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return nil, fmt.Errorf("channel not found: port=%v channel=%v", sourcePort, sourceChannel)
	}
	return k.sendPacket(ctx, data.GetBytes(), sourcePort, sourceChannel, c.Counterparty.PortID, c.Counterparty.ChannelID, data.GetTimeoutHeight(), data.GetTimeoutTimestamp())
}

// ReceiveRegisterDomainPacketAcknowledgement receive an ack to set self-domain name
func (k Keeper) ReceiveRegisterDomainPacketAcknowledgement(ctx sdk.Context, status uint32, domain string, packet channel.Packet) error {
	if status != servertypes.STATUS_OK {
		k.Logger(ctx).Info("failed to register a channel domain", "status", status, "domain", domain)
		return nil
	}
	return k.setSelfDomainName(ctx, types.NewLocalDNSID(packet.SourcePort, packet.SourceChannel), domain)
}

// SendDomainAssociationCreatePacketData sends a packet to associate a domain with client on DNS server
func (k Keeper) SendDomainAssociationCreatePacketData(ctx sdk.Context, dnsID types.LocalDNSID, srcClient, dstClient types.ClientDomain) (*channel.Packet, error) {
	_, found := k.GetSelfDomainName(ctx, dnsID)
	if !found {
		return nil, fmt.Errorf("this channel does not have a domain name: dnsID=%v", dnsID.String())
	}
	_, found = k.clientKeeper.GetClientState(ctx, srcClient.ClientId)
	if !found {
		return nil, fmt.Errorf("client state not found: clientID=%v", srcClient.ClientId)
	}
	data := servertypes.NewDomainAssociationCreatePacketData(srcClient, dstClient)
	c, found := k.channelKeeper.GetChannel(ctx, dnsID.SourcePort, dnsID.SourceChannel)
	if !found {
		return nil, fmt.Errorf("channel not found: port=%v channel=%v", dnsID.SourcePort, dnsID.SourceChannel)
	}
	return k.sendPacket(ctx, data.GetBytes(), dnsID.SourcePort, dnsID.SourceChannel, c.Counterparty.PortID, c.Counterparty.ChannelID, data.GetTimeoutHeight(), data.GetTimeoutTimestamp())
}

// ReceiveDomainAssociationCreatePacketAcknowledgement receive an ack
func (k Keeper) ReceiveDomainAssociationCreatePacketAcknowledgement(ctx sdk.Context, status uint32) error {
	// TODO cleanup any state
	if status != servertypes.STATUS_OK {
		k.Logger(ctx).Info("failed to register a local channel", "status", status)
		return nil
	}
	return nil
}

// ReceiveDomainAssociationResultPacketData receives a packet to save a domain info on local state
func (k Keeper) ReceiveDomainAssociationResultPacketData(ctx sdk.Context, packet channel.Packet, data servertypes.DomainAssociationResultPacketData) error {
	dnsID := types.NewLocalDNSID(packet.DestinationPort, packet.DestinationChannel)
	_, found := k.GetSelfDomainName(ctx, dnsID)
	if !found {
		return fmt.Errorf("failed to GetSelfDomainName")
	}

	store := ctx.KVStore(k.storeKey)
	bz0, err := proto.Marshal(&data.CounterpartyDomain.DNSID)
	if err != nil {
		return err
	}
	store.Set(clienttypes.KeyLocalDNSID(dnsID, data.CounterpartyDomain.Name), bz0)
	store.Set(clienttypes.KeyClientDomain(dnsID, data.CounterpartyDomain.Name), []byte(data.ClientId))
	return nil
}

// GetSelfDomainName returns self-domain name
func (k Keeper) GetSelfDomainName(ctx sdk.Context, dnsID types.LocalDNSID) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	key := clienttypes.KeySelfDomain(dnsID)
	domain := store.Get(key)
	if domain == nil {
		return "", false
	}
	return string(domain), true
}

// ResolveDNSID returns local DNS-ID corresponding to given a domain info
func (k Keeper) ResolveDNSID(ctx sdk.Context, domain types.LocalDomain) (types.LocalDNSID, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(clienttypes.KeyLocalDNSID(domain.DNSID, domain.Name))
	if bz == nil {
		return types.LocalDNSID{}, false
	}
	var id types.LocalDNSID
	err := proto.Unmarshal(bz, &id)
	if err != nil {
		panic(err)
	}
	return id, true
}

// ResolveChannel returns the channel corresponding to given a domain info
func (k Keeper) ResolveChannel(ctx sdk.Context, domain types.LocalDomain, portID string) (channel.Channel, bool) {
	key := clienttypes.KeyDomainChannel(domain.DNSID, domain.Name, portID)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return channel.Channel{}, false
	}
	var lc types.LocalChannel
	err := proto.Unmarshal(bz, &lc)
	if err != nil {
		return channel.Channel{}, false
	}
	c, found := k.channelKeeper.GetChannel(ctx, lc.SourcePort, lc.SourceChannel)
	if !found {
		return channel.Channel{}, false
	}
	return c, true
}

// SetDomainChannel sets a channel corresponding with the domain name
func (k Keeper) SetDomainChannel(ctx sdk.Context, dnsID types.LocalDNSID, domainName string, channel types.LocalChannel) error {
	c, found := k.channelKeeper.GetChannel(ctx, channel.SourcePort, channel.SourceChannel)
	if !found {
		return fmt.Errorf("channel not found: port=%v channel=%v", channel.SourcePort, channel.SourceChannel)
	}
	connID := c.ConnectionHops[0]
	conn, found := k.connectionKeeper.GetConnection(ctx, connID)
	if !found {
		return fmt.Errorf("connection not found: connectionID=%v", connID)
	}

	store := ctx.KVStore(k.storeKey)
	counterpartyClientID := store.Get(clienttypes.KeyClientDomain(dnsID, domainName))
	if counterpartyClientID == nil {
		return fmt.Errorf("clientID not found: dnsID=%v domainName=%v", dnsID.String(), domainName)
	}

	if string(counterpartyClientID) != conn.Counterparty.ClientID {
		return fmt.Errorf("clientID mismatch: %v != %v", string(counterpartyClientID), conn.Counterparty.ClientID)
	}

	bz, err := proto.Marshal(&channel)
	if err != nil {
		return err
	}
	store.Set(clienttypes.KeyDomainChannel(dnsID, domainName, channel.SourcePort), bz)
	return nil
}

func (k Keeper) setSelfDomainName(ctx sdk.Context, dnsID types.LocalDNSID, domain string) error {
	if _, found := k.GetSelfDomainName(ctx, dnsID); found {
		return fmt.Errorf("Domain name is already set")
	}
	store := ctx.KVStore(k.storeKey)
	key := clienttypes.KeySelfDomain(dnsID)
	store.Set(key, []byte(domain))
	return nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ibc/dns/client")
}

// PacketExecuted defines a wrapper function for the channel Keeper's function
// in order to expose it to the cross handler.
// Keeper retreives channel capability and passes it into channel keeper for authentication
func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(packet.GetDestPort(), packet.GetDestChannel()))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "channel capability could not be retrieved for packet")
	}
	return k.channelKeeper.PacketExecuted(ctx, chanCap, packet, acknowledgement)
}

func (k Keeper) sendPacket(
	ctx sdk.Context,
	data []byte,
	sourcePort,
	sourceChannel,
	destinationPort,
	destinationChannel string,
	timeoutHeight uint64,
	timeoutTimestamp uint64,
) (*channel.Packet, error) {
	// get the next sequence
	seq, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return nil, channel.ErrSequenceSendNotFound
	}
	packet := channel.NewPacket(
		data,
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return nil, sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	if err := k.channelKeeper.SendPacket(ctx, channelCap, packet); err != nil {
		return nil, err
	}

	return &packet, nil
}
