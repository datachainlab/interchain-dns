package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/datachainlab/interchain-dns/x/ibc-dns/common/types"
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
