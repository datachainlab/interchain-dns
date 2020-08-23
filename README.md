# Interchain DNS

## Introduction

Interchain DNS is a system that provides unique names for blockchains.

## Motivation

In IBC, it is possible to identify the other chain in the chain through client. However, it is a little difficult to pass a client identifier in one chain to another chain.

To achieve this, it is necessary to give a routes of channels that can be referenced by the other chain to which the identifier is passed. However, such a relative reference has the problem that there are multiple paths that represent the same chain or module.

Therefore, we propose DNS Module that provides a unique identifier for each chain.

## Implementation

Please check [here](https://github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns).

## Maintainers

- [Jun Kimura](https://github.com/bluele)
