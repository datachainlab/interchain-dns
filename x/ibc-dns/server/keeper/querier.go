package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	queryDomains = "domains"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case queryDomains:
			return handleQueryDomains(ctx, keeper, req)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown DNS %s query endpoint", commontypes.ModuleName)
		}
	}
}

func handleQueryDomains(ctx sdk.Context, k Keeper, req abci.RequestQuery) ([]byte, error) {
	var query types.QueryDomainsRequest
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &query); err != nil {
		return nil, err
	}
	res, err := k.queryDomains(ctx)
	if err != nil {
		return nil, err
	}
	return types.ModuleCdc.MarshalJSON(res)
}

func (k Keeper) queryDomains(ctx sdk.Context) (*types.QueryDomainsResponse, error) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixBytes(types.KeyForwardDomainPrefix))
	var res types.QueryDomainsResponse
	for ; iter.Valid(); iter.Next() {
		name := string(iter.Key())
		res.DomainNames = append(res.DomainNames, name)
	}
	return &res, nil
}
