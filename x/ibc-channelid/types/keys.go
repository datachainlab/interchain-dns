package types

const (
	// ModuleName is the name of the module
	ModuleName = "ibc-channel-id"

	// Version defines the current version the Cross
	// module supports
	Version = "ibc-channel-id-1"

	// PortID that Cross module binds to
	PortID = "channel-id"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the msg router key for the IBC module
	RouterKey string = ModuleName

	// QuerierRoute is the querier route for Cross
	QuerierRoute = ModuleName
)
