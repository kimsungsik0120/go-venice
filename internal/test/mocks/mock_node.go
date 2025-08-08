package mocks

import (
	"context"
	"math/big"

	"github.com/stretchr/testify/mock"
)

type MockNode struct {
	mock.Mock
}

func (m *MockNode) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*big.Int), args.Error(1)
}
func (m *MockNode) CreateTransferTransaction(fromAddress, toAddress, ethAmount string) (string, error) {
	args := m.Called(fromAddress, toAddress, ethAmount)
	return args.Get(0).(string), args.Error(1)
}
