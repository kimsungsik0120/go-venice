package nodes

import (
	"context"
	"math/big"
)

type Node interface {
	GetBalance(ctx context.Context, address string) (*big.Int, error)
	CreateTransferTransaction(fromAddress, toAddress, ethAmount string) (string, error)
}
