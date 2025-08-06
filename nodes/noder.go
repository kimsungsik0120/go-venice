package nodes

import "math/big"

type Noder interface {
	GetBalance(address string) (*big.Int, error)
	GetBalanceToken(address string) (*big.Int, error)
}
