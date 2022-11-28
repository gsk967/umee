package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddressFromKey extracts address from a key with the form
// prefix | lengthPrefixed(addr) | ...
func AddressFromKey(key, prefix []byte) sdk.AccAddress {
	addrLength := int(key[len(prefix)])
	return key[len(prefix)+1 : len(prefix)+1+addrLength]
}

// DenomFromKeyWithAddress extracts denom from a key with the form
// prefix | lengthPrefixed(addr) | denom | 0x00
func DenomFromKeyWithAddress(key, prefix []byte) string {
	addrLength := int(key[len(prefix)])
	return string(key[len(prefix)+addrLength+1 : len(key)-1])
}

// DenomFromKey extracts denom from a key with the form
// prefix | denom | 0x00
func DenomFromKey(key, prefix []byte) string {
	return string(key[len(prefix) : len(key)-1])
}
