package dns

import (
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/keeper"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/types"
)

// nolint
const (
	ModuleName = commontypes.ModuleName
	RouterKey  = commontypes.RouterKey
	PortID     = commontypes.PortID
)

// nolint
type (
	Keeper                        = keeper.Keeper
	PacketReceiver                = commontypes.PacketReceiver
	PacketAcknowledgementReceiver = commontypes.PacketAcknowledgementReceiver
	ChannelKeeper                 = commontypes.ChannelKeeper
	PortKeeper                    = commontypes.PortKeeper
)

// nolint
var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
)
