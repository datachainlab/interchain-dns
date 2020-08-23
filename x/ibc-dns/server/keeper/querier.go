package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	commontypes "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	"github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/server/types"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	queryDomain  = "domain"
	queryDomains = "domains"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case queryDomain:
			return handleQueryDomain(ctx, keeper, req)
		case queryDomains:
			return handleQueryDomains(ctx, keeper, req)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown DNS %s query endpoint", commontypes.ModuleName)
		}
	}
}

func handleQueryDomain(ctx sdk.Context, k Keeper, req abci.RequestQuery) ([]byte, error) {
	var query types.QueryDomainRequest
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &query); err != nil {
		return nil, err
	}
	res, err := k.QueryDomain(ctx, query)
	if err != nil {
		return nil, err
	}
	return types.ModuleCdc.MarshalJSON(res)
}

func handleQueryDomains(ctx sdk.Context, k Keeper, req abci.RequestQuery) ([]byte, error) {
	var query types.QueryDomainsRequest
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &query); err != nil {
		return nil, err
	}
	res, err := k.QueryDomains(ctx)
	if err != nil {
		return nil, err
	}
	return types.ModuleCdc.MarshalJSON(res)
}

// QueryDomain returns a domain info corresponding to a given name
func (k Keeper) QueryDomain(ctx sdk.Context, req types.QueryDomainRequest) (*types.QueryDomainResponse, error) {
	store := ctx.KVStore(k.storeKey)
	var info types.DomainChannelInfo
	bz := store.Get(types.KeyForwardDomain(req.Name))
	if bz == nil {
		return nil, fmt.Errorf("domain not found: %v", req.Name)
	}
	if err := proto.Unmarshal(bz, &info); err != nil {
		return nil, err
	}
	return &types.QueryDomainResponse{
		Domain: types.DomainInfo{
			Name:     req.Name,
			Metadata: info.Metadata,
			DnsId:    commontypes.NewLocalDNSID(info.Channel.SourcePort, info.Channel.SourceChannel),
			Channel:  info.Channel,
		},
	}, nil
}

// QueryDomains returns all domain info
func (k Keeper) QueryDomains(ctx sdk.Context) (*types.QueryDomainsResponse, error) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.KeyPrefixBytes(types.KeyForwardDomainPrefix)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	var res types.QueryDomainsResponse
	for ; iter.Valid(); iter.Next() {
		var info types.DomainChannelInfo
		if err := proto.Unmarshal(iter.Value(), &info); err != nil {
			return nil, err
		}
		name := strings.TrimPrefix(string(iter.Key()), string(prefix))
		res.Domains = append(res.Domains, &types.DomainInfo{
			Name:     name,
			Metadata: info.Metadata,
			DnsId:    commontypes.NewLocalDNSID(info.Channel.SourcePort, info.Channel.SourceChannel),
			Channel:  info.Channel,
		})
	}
	return &res, nil
}
