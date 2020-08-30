package types

import (
	"testing"

	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/stretchr/testify/require"
)

func TestPacketSerialization(t *testing.T) {
	require := require.New(t)
	cdc := PacketCdc()

	var p = &RegisterDomainPacketData{DomainName: "name"}
	bz := p.GetBytes()

	pd, err := types.DeserializeJSONPacketData(cdc, bz)
	require.NoError(err)
	require.Equal(p, pd)
}
