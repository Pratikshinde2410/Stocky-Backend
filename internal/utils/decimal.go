package utils

import "math/big"

func NewDecimalFromString(s string) (*big.Rat, bool) {
    r := new(big.Rat)
    _, ok := r.SetString(s)
    return r, ok
}


