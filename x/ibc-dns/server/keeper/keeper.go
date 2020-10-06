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

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
)

// Keeper defines ibc-dns keeper
type Keeper struct {
	storeKey      sdk.StoreKey
	channelKeeper types.ChannelKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new ibc-dns Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	channelKeeper types.ChannelKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	return Keeper{
		storeKey:      storeKey,
		channelKeeper: channelKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

// ReceivePacketRegisterDomain receives a PacketRegisterDomain to register a new domain record
func (k Keeper) ReceivePacketRegisterDomain(ctx sdk.Context, packet channel.Packet, data servertypes.RegisterDomainPacketData) error {
	c := types.NewChannel(packet.SourcePort, packet.SourceChannel, packet.DestinationPort, packet.DestinationChannel)
	return k.registerDomain(ctx, data.DomainName, c, data.Metadata)
}

// ReceiveDomainAssociationCreatePacketData receives a DomainAssociationCreatePacketData to associate domain with client
func (k Keeper) ReceiveDomainAssociationCreatePacketData(ctx sdk.Context, packet channel.Packet, data servertypes.DomainAssociationCreatePacketData) (ack servertypes.DomainAssociationCreatePacketAcknowledgement, completed bool) {
	// check if counterparty domain exists
	_, err := k.ForwardLookupDomain(ctx, data.DstClient.DomainName)
	if err != nil {
		return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), false
	}
	srcName, err := k.ReverseLookupDomain(ctx, packet.DestinationPort, packet.DestinationChannel)
	if err != nil {
		return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), false
	}
	if data.SrcClient.DomainName != srcName {
		return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_FAILED, fmt.Sprintf("unexpected domain name: actual=%v expected=%v", data.SrcClient.DomainName, srcName)), false
	}

	if k.ensureClientDomainExistence(ctx, data.SrcClient, data.DstClient) {
		return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_FAILED, fmt.Sprintf("this association is already created: src=%v dst=%v", data.SrcClient.String(), data.DstClient.String())), false
	}

	// check if opposite association already exists
	// if exists, try to confirm it
	// if not exists, try to create a new association
	exists := k.ensureClientDomainExistence(ctx, data.DstClient, data.SrcClient)
	if !exists {
		err := k.createDomainAssociation(ctx, data.SrcClient, data.DstClient)
		if err == nil {
			return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_OK, "ok"), false
		} else {
			return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), false
		}
	}

	if err := k.confirmDomainAssociation(ctx, data.SrcClient, data.DstClient); err != nil {
		k.Logger(ctx).Info("failed to confirm the domain association", "err", err)
		return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), true
	} else {
		return servertypes.NewDomainAssociationCreatePacketAcknowledgement(servertypes.STATUS_OK, "ok"), true
	}
}

// CreateDomainAssociationResultPacketData creates a packet 'DomainAssociationResultPacketData'
func (k Keeper) CreateDomainAssociationResultPacketData(ctx sdk.Context, status uint32, srcClientDomain, dstClientDomain types.ClientDomain) (srcPacket *channel.Packet, dstPacket *channel.Packet, err error) {
	srcLocalDNSID, err := k.GetLocalDNSID(ctx, srcClientDomain.DomainName)
	if err != nil {
		return
	}
	dstLocalDNSID, err := k.GetLocalDNSID(ctx, dstClientDomain.DomainName)
	if err != nil {
		return
	}
	srcDomainInfo, err := k.ForwardLookupDomain(ctx, srcClientDomain.DomainName)
	if err != nil {
		return
	}
	dstDomainInfo, err := k.ForwardLookupDomain(ctx, dstClientDomain.DomainName)
	if err != nil {
		return
	}

	srcChannel := srcDomainInfo.Channel
	dstChannel := dstDomainInfo.Channel

	toSrcData := servertypes.NewDomainAssociationResultPacketData(
		status,
		dstClientDomain.DomainName,
		*dstLocalDNSID,
		dstClientDomain.ClientId,
	)

	toDstData := servertypes.NewDomainAssociationResultPacketData(
		status,
		srcClientDomain.DomainName,
		*srcLocalDNSID,
		srcClientDomain.ClientId,
	)

	srcPacket, err = k.createPacket(
		ctx,
		toSrcData.GetBytes(),
		srcChannel.DestinationPort,
		srcChannel.DestinationChannel,
		srcChannel.SourcePort,
		srcChannel.SourceChannel,
		toSrcData.GetTimeoutHeight(),
		toSrcData.GetTimeoutTimestamp(),
	)
	if err != nil {
		return
	}

	dstPacket, err = k.createPacket(
		ctx,
		toDstData.GetBytes(),
		dstChannel.DestinationPort,
		dstChannel.DestinationChannel,
		dstChannel.SourcePort,
		dstChannel.SourceChannel,
		toDstData.GetTimeoutHeight(),
		toDstData.GetTimeoutTimestamp(),
	)
	if err != nil {
		return
	}

	return
}

// SendDomainAssociationResultPacketData sends a result of the domain association
func (k Keeper) SendDomainAssociationResultPacketData(ctx sdk.Context, status uint32, srcClientDomain, dstClientDomain types.ClientDomain) error {
	srcPacket, dstPacket, err := k.CreateDomainAssociationResultPacketData(ctx, status, srcClientDomain, dstClientDomain)
	if err != nil {
		return err
	}
	if err := k.sendPacket(ctx, srcPacket); err != nil {
		return err
	}
	if err := k.sendPacket(ctx, dstPacket); err != nil {
		return nil
	}
	return nil
}

// GetLocalDNSID returns a local DNS-ID corresponding to given name
func (k Keeper) GetLocalDNSID(ctx sdk.Context, name string) (*types.LocalDNSID, error) {
	c, err := k.ForwardLookupDomain(ctx, name)
	if err != nil {
		return nil, err
	}
	return &types.LocalDNSID{SourcePort: c.Channel.SourcePort, SourceChannel: c.Channel.SourceChannel}, nil
}

// ForwardLookupDomain returns a local channel info corresponding to given name
func (k Keeper) ForwardLookupDomain(ctx sdk.Context, name string) (*servertypes.DomainChannelInfo, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(servertypes.KeyForwardDomain(name))
	if bz == nil {
		return nil, fmt.Errorf("Domain '%v' not found", name)
	}
	var info servertypes.DomainChannelInfo
	if err := proto.Unmarshal(bz, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ReverseLookupDomain returns a domain name corresponding to given channel info
func (k Keeper) ReverseLookupDomain(ctx sdk.Context, port, channel string) (string, error) {
	store := ctx.KVStore(k.storeKey)
	name := store.Get(servertypes.KeyReverseDomain(port, channel))
	if name == nil {
		return "", fmt.Errorf("failed to reverseLookup: port=%v channel=%v", port, channel)
	}
	return string(name), nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ibc/dns/server")
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

func (k Keeper) registerDomain(ctx sdk.Context, domain string, channel types.LocalChannel, metadata []byte) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(servertypes.KeyForwardDomain(domain)) {
		return fmt.Errorf("Domain name '%v' already exists", domain)
	}
	info := servertypes.DomainChannelInfo{
		Metadata: metadata,
		Channel:  channel,
	}
	bz, err := proto.Marshal(&info)
	if err != nil {
		return err
	}
	// for forward lookup
	store.Set(servertypes.KeyForwardDomain(domain), bz)
	// for reverse lookup
	store.Set(servertypes.KeyReverseDomain(channel.DestinationPort, channel.DestinationChannel), []byte(domain))
	return nil
}

func (k Keeper) ensureClientDomainExistence(ctx sdk.Context, srcClientDomain, dstClientDomain types.ClientDomain) (exists bool) {
	store := ctx.KVStore(k.storeKey)
	key := servertypes.KeyDomainAssociation(srcClientDomain, dstClientDomain)
	return store.Has(key)
}

func (k Keeper) createDomainAssociation(ctx sdk.Context, srcClientDomain, dstClientDomain types.ClientDomain) error {
	store := ctx.KVStore(k.storeKey)
	srcKey := servertypes.KeyDomainAssociation(srcClientDomain, dstClientDomain)

	r := servertypes.NewDomainAssociation(
		servertypes.DomainAssociationStatusConfirmed,
		srcClientDomain,
		dstClientDomain,
	)
	bz, err := proto.Marshal(&r)
	if err != nil {
		return err
	}
	store.Set(srcKey, bz)
	return nil
}

func (k Keeper) confirmDomainAssociation(ctx sdk.Context, srcClientDomain, dstClientDomain types.ClientDomain) error {
	store := ctx.KVStore(k.storeKey)
	dstKey := servertypes.KeyDomainAssociation(dstClientDomain, srcClientDomain)

	var da types.DomainAssociation
	if err := proto.Unmarshal(store.Get(dstKey), &da); err != nil {
		return err
	}

	if da.Status != servertypes.DomainAssociationStatusInit {
		return nil
	}

	da.Status = servertypes.DomainAssociationStatusConfirmed
	bz, err := proto.Marshal(&da)
	if err != nil {
		return err
	}
	store.Set(dstKey, bz)
	return nil
}

func (k Keeper) createPacket(
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
	return &packet, nil
}

func (k Keeper) sendPacket(
	ctx sdk.Context,
	packet *channel.Packet,
) error {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(packet.SourcePort, packet.SourceChannel))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	if err := k.channelKeeper.SendPacket(ctx, channelCap, packet); err != nil {
		return err
	}

	return nil
}
