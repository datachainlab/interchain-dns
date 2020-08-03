package types

const (
	// ModuleName is the name of the module
	ModuleName = "ibcdns"

	// Version defines the current version the DNS module supports
	Version = "ibcdns-1"

	// PortID that DNS module binds to
	PortID = "ibcdns"

	// RouterKey is the msg router key for the IBC module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for Cross
	QuerierRoute = ModuleName
)
