package dns

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	dnsclient "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/client"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/types"
)

const (
	flagServer uint8 = 1 << iota
	flagClient
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ port.IBCModule        = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is an app module Basics object
type AppModuleBasic struct {
	flags uint8
}

func NewAppModuleBasic(flags uint8) AppModuleBasic {
	return AppModuleBasic{flags: flags}
}

// Name returns module name
func (AppModuleBasic) Name() string {
	return commontypes.ModuleName
}

// RegisterCodec returns RegisterCodec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state
func (AppModuleBasic) DefaultGenesis(m codec.JSONMarshaler) json.RawMessage {
	return m.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis checks the Genesis
func (AppModuleBasic) ValidateGenesis(m codec.JSONMarshaler, bz json.RawMessage) error {
	var data GenesisState
	err := m.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

// RegisterRESTRoutes returns rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
}

// GetQueryCmd returns the root query command of this module
func (AppModuleBasic) GetQueryCmd(clientCtx client.Context) *cobra.Command {
	// return cli.GetQueryCmd(cdc)
	return nil
}

// GetTxCmd returns the root tx command of this module
func (AppModuleBasic) GetTxCmd(clientCtx client.Context) *cobra.Command {
	// return cli.GetTxCmd(cdc)
	return nil
}

// AppModule struct
type AppModule struct {
	AppModuleBasic
	keeper                        Keeper
	handler                       sdk.Handler
	querier                       sdk.Querier
	packetReceiver                PacketReceiver
	packetAcknowledgementReceiver PacketAcknowledgementReceiver
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k Keeper, ck *dnsclient.Keeper, sk *server.Keeper) AppModule {
	var (
		flags uint8
		hs    []sdk.Handler
		qs    []sdk.Querier
		rs    []PacketReceiver
		as    []PacketAcknowledgementReceiver
	)
	if ck != nil {
		flags |= flagClient
		hs = append(hs, client.NewHandler(*ck))
		rs = append(rs, client.NewPacketReceiver(*ck))
		as = append(as, client.NewPacketAcknowledgementReceiver(*ck))
	}
	if sk != nil {
		flags |= flagServer
		qs = append(qs, server.NewQuerier(*sk))
		rs = append(rs, server.NewPacketReceiver(*sk))
		as = append(as, server.NewPacketAcknowledgementReceiver(*sk))
	}
	return AppModule{
		AppModuleBasic:                NewAppModuleBasic(flags),
		keeper:                        k,
		handler:                       commontypes.ComposeHandlers(hs...),
		querier:                       commontypes.ComposeQuerier(qs...),
		packetReceiver:                commontypes.ComposePacketReceivers(rs...),
		packetAcknowledgementReceiver: commontypes.ComposePacketAcknowledgementReceivers(as...),
	}
}

// Name returns module name
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants is empty
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns RouterKey
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(RouterKey, am.handler)
}

// NewHandler returns new Handler
func (am AppModule) NewHandler() sdk.Handler {
	return am.handler
}

// QuerierRoute returns module name
func (am AppModule) QuerierRoute() string {
	return ModuleName
}

// NewQuerierHandler returns new Querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return am.querier
}

// BeginBlock is a callback function
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock is a callback function
func (am AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// InitGenesis inits genesis
func (am AppModule) InitGenesis(ctx sdk.Context, m codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	m.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis exports genesis
func (am AppModule) ExportGenesis(ctx sdk.Context, m codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return m.MustMarshalJSON(gs)
}

func (am AppModule) RegisterQueryService(grpc.Server) {

}

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	if counterparty.PortID != commontypes.PortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "counterparty has invalid portid. expected: %s, got %s", commontypes.PortID, counterparty.PortID)
	}

	if version != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, commontypes.Version)
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error()+"by cross chanOpenInit")
	}

	return nil
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	if counterparty.PortID != commontypes.PortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "counterparty has invalid portid. expected: %s, got %s", commontypes.PortID, counterparty.PortID)
	}

	if version != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, commontypes.Version)
	}

	if counterpartyVersion != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, commontypes.Version)
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error()+"by cross chanOpenTry")
	}

	// TODO: escrow
	return nil
}

func (am AppModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyVersion string,
) error {
	if counterpartyVersion != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, commontypes.Version)
	}
	return nil
}

func (am AppModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {
	return am.packetReceiver(ctx, packet)
}

func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {
	var ack commontypes.PacketAcknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	return am.packetAcknowledgementReceiver(ctx, packet, ack)
}

func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	return nil, nil
}
