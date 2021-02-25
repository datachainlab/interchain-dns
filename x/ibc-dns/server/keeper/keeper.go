package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
	servertypes "github.com/datachainlab/interchain-dns/x/ibc-dns/server/types"
)

// Keeper defines ibc-dns keeper
type Keeper struct {
	cdc           codec.BinaryMarshaler
	storeKey      sdk.StoreKey
	channelKeeper types.ChannelKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new ibc-dns Keeper
func NewKeeper(
	cdc codec.BinaryMarshaler,
	storeKey sdk.StoreKey,
	channelKeeper types.ChannelKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		channelKeeper: channelKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

// ReceivePacketRegisterDomain receives a PacketRegisterDomain to register a new domain record
func (k Keeper) ReceivePacketRegisterDomain(ctx sdk.Context, packet channeltypes.Packet, data *servertypes.RegisterDomainPacketData) error {
	c := types.NewChannel(packet.SourcePort, packet.SourceChannel, packet.DestinationPort, packet.DestinationChannel)
	return k.registerDomain(ctx, data.DomainName, c, data.Metadata)
}

// ReceiveDomainMappingCreatePacketData receives a DomainMappingCreatePacketData to associate domain with client
func (k Keeper) ReceiveDomainMappingCreatePacketData(ctx sdk.Context, packet channeltypes.Packet, data *servertypes.DomainMappingCreatePacketData) (ack servertypes.DomainMappingCreatePacketAcknowledgement, completed bool) {
	// check if counterparty domain exists
	_, err := k.ForwardLookupDomain(ctx, data.DstClient.DomainName)
	if err != nil {
		return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), false
	}
	srcName, err := k.ReverseLookupDomain(ctx, packet.DestinationPort, packet.DestinationChannel)
	if err != nil {
		return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), false
	}
	if data.SrcClient.DomainName != srcName {
		return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_FAILED, fmt.Sprintf("unexpected domain name: actual=%v expected=%v", data.SrcClient.DomainName, srcName)), false
	}

	if k.ensureClientDomainExistence(ctx, data.SrcClient, data.DstClient) {
		return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_FAILED, fmt.Sprintf("this association is already created: src=%v dst=%v", data.SrcClient.String(), data.DstClient.String())), false
	}

	// check if opposite association already exists
	// if exists, try to confirm it
	// if not exists, try to create a new association
	exists := k.ensureClientDomainExistence(ctx, data.DstClient, data.SrcClient)
	if !exists {
		err := k.createDomainMapping(ctx, data.SrcClient, data.DstClient)
		if err == nil {
			return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_OK, "ok"), false
		} else {
			return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), false
		}
	}

	if err := k.confirmDomainMapping(ctx, data.SrcClient, data.DstClient); err != nil {
		k.Logger(ctx).Info("failed to confirm the domain mapping", "err", err)
		return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_FAILED, err.Error()), true
	} else {
		return servertypes.NewDomainMappingCreatePacketAcknowledgement(servertypes.STATUS_OK, "ok"), true
	}
}

// CreateDomainMappingResultPacketData creates a packet 'DomainMappingResultPacketData'
func (k Keeper) CreateDomainMappingResultPacketData(ctx sdk.Context, status uint32, srcClientDomain, dstClientDomain types.ClientDomain) (srcPacket *channeltypes.Packet, dstPacket *channeltypes.Packet, err error) {
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

	toSrcData := servertypes.NewDomainMappingResultPacketData(
		status,
		dstClientDomain.DomainName,
		*dstLocalDNSID,
		dstClientDomain.ClientId,
	)

	toDstData := servertypes.NewDomainMappingResultPacketData(
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

// SendDomainMappingResultPacketData sends a result of the domain mapping
func (k Keeper) SendDomainMappingResultPacketData(ctx sdk.Context, status uint32, srcClientDomain, dstClientDomain types.ClientDomain) error {
	srcPacket, dstPacket, err := k.CreateDomainMappingResultPacketData(ctx, status, srcClientDomain, dstClientDomain)
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

func (k Keeper) Codec() codec.BinaryMarshaler {
	return k.cdc
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ibc/dns/server")
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
	key := servertypes.KeyDomainMapping(srcClientDomain, dstClientDomain)
	return store.Has(key)
}

func (k Keeper) createDomainMapping(ctx sdk.Context, srcClientDomain, dstClientDomain types.ClientDomain) error {
	store := ctx.KVStore(k.storeKey)
	srcKey := servertypes.KeyDomainMapping(srcClientDomain, dstClientDomain)

	r := servertypes.NewDomainMapping(
		servertypes.DomainMappingStatusConfirmed,
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

func (k Keeper) confirmDomainMapping(ctx sdk.Context, srcClientDomain, dstClientDomain types.ClientDomain) error {
	store := ctx.KVStore(k.storeKey)
	dstKey := servertypes.KeyDomainMapping(dstClientDomain, srcClientDomain)

	var da types.DomainMapping
	if err := proto.Unmarshal(store.Get(dstKey), &da); err != nil {
		return err
	}

	if da.Status != servertypes.DomainMappingStatusInit {
		return nil
	}

	da.Status = servertypes.DomainMappingStatusConfirmed
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
	timeoutHeight ibcclienttypes.Height,
	timeoutTimestamp uint64,
) (*channeltypes.Packet, error) {
	// get the next sequence
	seq, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return nil, channeltypes.ErrSequenceSendNotFound
	}
	packet := channeltypes.NewPacket(
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
	packet *channeltypes.Packet,
) error {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(packet.SourcePort, packet.SourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	if err := k.channelKeeper.SendPacket(ctx, channelCap, packet); err != nil {
		return err
	}

	return nil
}
