package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/desmos-labs/desmos/x/profiles/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.GetDTagTransferRequests(ctx),
		k.GetAllRelationships(ctx),
		k.GetAllUsersBlocks(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the chain state based on the given GenesisState
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) []abci.ValidatorUpdate {
	// Initialize the module params
	k.SetParams(ctx, data.Params)

	// Initialize the Profiles
	k.ak.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		if profile, ok := (account).(*types.Profile); ok {
			err := k.StoreProfile(ctx, profile)
			if err != nil {
				panic(err)
			}
		}
		return false
	})

	// Store the transfer requests
	for _, request := range data.DTagTransferRequest {
		err := k.SaveDTagTransferRequest(ctx, request)
		if err != nil {
			panic(err)
		}
	}

	// Store the relationships
	for _, relationship := range data.Relationships {
		err := k.SaveRelationship(ctx, relationship)
		if err != nil {
			panic(err)
		}
	}

	// Store the user blocks
	for _, userBlock := range data.Blocks {
		err := k.SaveUserBlock(ctx, userBlock)
		if err != nil {
			panic(err)
		}
	}

	return nil
}