package dns

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"

	dnsclient "github.com/datachainlab/interchain-dns/x/ibc-dns/client"
	clientkeeper "github.com/datachainlab/interchain-dns/x/ibc-dns/client/keeper"
	commontypes "github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/interchain-dns/x/ibc-dns/keeper"
	dnskeeper "github.com/datachainlab/interchain-dns/x/ibc-dns/keeper"
	dnsserver "github.com/datachainlab/interchain-dns/x/ibc-dns/server"
	serverkeeper "github.com/datachainlab/interchain-dns/x/ibc-dns/server/keeper"
	"github.com/datachainlab/interchain-dns/x/ibc-dns/types"
)

const (
	flagServer uint8 = 1 << iota
	flagClient
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ porttypes.IBCModule   = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is an app module Basics object
type AppModuleBasic struct {
	flags uint8
}

func NewAppModuleBasic(flags uint8) AppModuleBasic {
	return AppModuleBasic{flags: flags}
}

// Name implements AppModuleBasic interface
func (AppModuleBasic) Name() string {
	return commontypes.ModuleName
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// transfer module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(commontypes.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the ibc transfer module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var gs commontypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", commontypes.ModuleName, err)
	}
	return gs.Validate()
}

// RegisterRESTRoutes implements AppModuleBasic interface
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(ctx client.Context, serveMux *runtime.ServeMux) {
}

// GetTxCmd implements AppModuleBasic interface
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	// return cli.NewTxCmd()
	return nil
}

// GetQueryCmd implements AppModuleBasic interface
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	// return cli.GetQueryCmd()
	return nil
}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule struct
type AppModule struct {
	AppModuleBasic
	keeper                        keeper.Keeper
	handler                       sdk.Handler
	querier                       sdk.Querier
	packetReceiver                commontypes.PacketReceiver
	packetAcknowledgementReceiver commontypes.PacketAcknowledgementReceiver
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k dnskeeper.Keeper, ck *clientkeeper.Keeper, sk *serverkeeper.Keeper) AppModule {
	var (
		flags uint8
		hs    []sdk.Handler
		qs    []sdk.Querier
		rs    []commontypes.PacketReceiver
		as    []commontypes.PacketAcknowledgementReceiver
	)
	if ck != nil {
		flags |= flagClient
		hs = append(hs, dnsclient.NewHandler(*ck))
		rs = append(rs, dnsclient.NewPacketReceiver(*ck))
		as = append(as, dnsclient.NewPacketAcknowledgementReceiver(*ck))
	}
	if sk != nil {
		flags |= flagServer
		qs = append(qs, serverkeeper.NewQuerier(*sk))
		rs = append(rs, dnsserver.NewPacketReceiver(*sk))
		as = append(as, dnsserver.NewPacketAcknowledgementReceiver(*sk))
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

// RegisterInvariants implements the AppModule interface
func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO
}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(commontypes.RouterKey, am.handler)
}

// QuerierRoute returns module name
func (am AppModule) QuerierRoute() string {
	return commontypes.QuerierRoute
}

func (am AppModule) RegisterServices(configurator module.Configurator) {
	return
}

// LegacyQuerierHandler implements the AppModule interface
func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterQueryService(server grpc.Server) {
	// commontypes.RegisterQueryServer(server, am.keeper)
}

// InitGenesis performs genesis initialization for the ibc-transfer module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState commontypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the ibc-transfer
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the transfer module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	// simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized ibc-transfer param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	// return simulation.ParamChanges(r)
	return nil
}

// RegisterStoreDecoder registers a decoder for transfer module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	// sdr[types.StoreKey] = simulation.NewDecodeStore(am.keeper)
}

// WeightedOperations returns the all the transfer module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

//____________________________________________________________________________

// OnChanOpenInit implements the IBCModule interface
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

	if counterparty.GetPortID() != commontypes.PortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "counterparty has invalid portid. expected: %s, got %s", commontypes.PortID, counterparty.GetPortID())
	}

	if version != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, commontypes.Version)
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, err.Error()+"by cross chanOpenInit")
	}

	return nil
}

// OnChanOpenTry implements the IBCModule interface
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

	if counterparty.GetPortID() != commontypes.PortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "counterparty has invalid portid. expected: %s, got %s", commontypes.PortID, counterparty.GetPortID())
	}

	if version != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, commontypes.Version)
	}

	if counterpartyVersion != commontypes.Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, commontypes.Version)
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, err.Error()+"by cross chanOpenTry")
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
	return am.packetAcknowledgementReceiver(ctx, packet, acknowledgement)
}

func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	return nil, nil
}
