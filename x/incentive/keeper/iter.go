package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/umee-network/umee/v3/util/store"
	"github.com/umee-network/umee/v3/x/incentive"
)

// GetPendingRewards returns an sdk.Coins object containing all pending rewards
// associated with an address.
func (k Keeper) GetPendingRewards(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	kvStore := ctx.KVStore(k.storeKey)
	prefix := incentive.CreatePendingRewardKeyNoDenom(addr)
	pendingReward := sdk.NewCoins()

	iterator := func(key, val []byte) error {
		// get reward denom from key
		denom := store.DenomFromKeyWithAddress(key, prefix)

		// get pending reward (panic on unmarshal fail)
		amount := store.GetStoredInt(kvStore, key, sdk.ZeroInt(), "pending reward")

		// add to pendingReward
		pendingReward = pendingReward.Add(sdk.NewCoin(denom, amount))
		return nil
	}

	store.Iterate(kvStore, ctx, prefix, iterator)

	return pendingReward
}

// GetRewardBasis returns an sdk.Coins object containing the reward basis associated with
// a given address and lock denom.
func (k Keeper) GetRewardBasis(ctx sdk.Context, addr sdk.AccAddress, lockDenom string) sdk.DecCoins {
	kvStore := ctx.KVStore(k.storeKey)
	prefix := incentive.CreateRewardBasisKeyNoRewardDenom(addr, lockDenom)
	rewardBasis := sdk.NewDecCoins()

	iterator := func(key, val []byte) error {
		// get reward denom from key
		denom := store.TrailingDenomFromKey(key, prefix)

		// get pending reward (panic on unmarshal fail)
		amount := store.GetStoredDec(kvStore, key, sdk.ZeroDec(), "reward basis")

		// add to pendingReward
		rewardBasis = rewardBasis.Add(sdk.NewDecCoinFromDec(denom, amount))
		return nil
	}

	store.Iterate(kvStore, ctx, prefix, iterator)

	return rewardBasis
}
