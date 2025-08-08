package service_test

import (
	"context"
	"github.com/bxcodec/faker/v4"
	"github.com/bxcodec/faker/v4/pkg/interfaces"
	"github.com/bxcodec/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-venice/configs"
	"go-venice/internal/api/dto"
	"go-venice/internal/api/service"
	"go-venice/internal/test/mocks"
	"math/big"
	"strconv"
	"testing"
)

type testEnv struct {
	svc      service.NodeService
	mockNode *mocks.MockNode
}

func setupTestEnv() *testEnv {
	mockNode := new(mocks.MockNode)
	cfg := fakeConfig() // 테스트용 config 생성

	svc := service.NewNodeService(cfg, mockNode)

	return &testEnv{
		svc:      svc,
		mockNode: mockNode,
	}
}

func fakeConfig() *configs.EnvConfig {
	return &configs.EnvConfig{
		RpcUrl:  "http://fake-rpc",
		ChainId: 1,
	}
}

func TestGetBalance(t *testing.T) {
	// given: mock Noder
	env := setupTestEnv()

	//given
	var fakerInt int
	_ = faker.FakeData(&fakerInt, options.WithRandomIntegerBoundaries(interfaces.RandomIntegerBoundary{Start: 1, End: 100}))
	if fakerInt < 0 {
		fakerInt = -fakerInt // 양수로 변환
	}
	testBalance := new(big.Int).Mul(big.NewInt(int64(fakerInt)), big.NewInt(1e18))

	env.mockNode.
		On("GetBalance", mock.Anything, "0x123").
		Return(testBalance, nil)

	// and: service
	cfg := configs.EnvConfig{ /* 필요한 필드 채우기 */ }
	svc := service.NewNodeService(&cfg, env.mockNode)

	// when
	resp, err := svc.GetBalance(context.Background(), "0x123")

	// then

	expectNumberString := strconv.Itoa(fakerInt)
	assert.NoError(t, err)
	assert.Equal(t, &dto.BalanceResponse{
		Amount: expectNumberString, // utils.DivideBy 적용 결과
		Symbol: "BASE",
	}, resp)

	env.mockNode.AssertExpectations(t)
}
