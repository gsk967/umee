package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Iterate through all keys in a KVStore with a given prefix using a provided function.
// If the provided function returns an error, iteration stops and the error is returned.
func Iterate(store sdk.KVStore, ctx sdk.Context, prefix []byte, cb func(key, val []byte) error) error {
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, val := iter.Key(), iter.Value()

		if err := cb(key, val); err != nil {
			return err
		}
	}

	return nil
}
