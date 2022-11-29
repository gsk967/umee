package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/umee-network/umee/v3/util/store"
	"github.com/umee-network/umee/v3/x/incentive"
)

// SetPendingRewards sets pending rewards for all reward denoms assocated with a given address
// using an sdk.Coins after clearing the addresses existing pending rewards
func (k Keeper) SetPendingRewards(ctx sdk.Context, addr sdk.AccAddress, rewards sdk.Coins) error {
	// clear all existing rewards
	existingRewards := k.GetPendingRewards(ctx, addr)
	for _, v := range existingRewards {
		// we only need to clear an entry if it will not be overwritten with a nonzero amount
		if rewards.AmountOf(v.Denom).IsZero() {
			if err := store.SetStoredInt(
				k.kvStore(ctx),
				incentive.CreatePendingRewardKey(addr, v.Denom),
				sdk.ZeroInt(),
				sdk.ZeroInt(),
				"pending reward",
			); err != nil {
				return err
			}
		}
	}

	// set all nonzero pending reward amounts
	for _, v := range rewards {
		if err := store.SetStoredInt(
			k.kvStore(ctx),
			incentive.CreatePendingRewardKey(addr, v.Denom),
			v.Amount,
			sdk.ZeroInt(),
			"pending reward",
		); err != nil {
			return err
		}
	}

	return nil
}
