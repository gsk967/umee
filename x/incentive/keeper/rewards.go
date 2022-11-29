package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/umee-network/umee/v3/x/incentive"
)

// AllocatePendingRewards calculates and then sets pending rewards for an account,
// increasing its reward basis for any bonded uTojen denoms that have distributed
// reward tokens since the account's last allocation.
func (k Keeper) AllocatePendingRewards(ctx sdk.Context, addr sdk.AccAddress) error {
	return incentive.ErrNotImplemented
}
