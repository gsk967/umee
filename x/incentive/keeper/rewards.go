package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/umee-network/umee/v3/x/incentive"
)

func (k Keeper) AllocatePendingRewards(ctx sdk.Context, addr sdk.AccAddress) error {
	return incentive.ErrNotImplemented
}
