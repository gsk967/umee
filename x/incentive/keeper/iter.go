package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllPendingRewards(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.NewCoins()
}
