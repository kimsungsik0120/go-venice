package service

import (
	"go-venice/configs"
	"go-venice/internal/adapter/nodes"
	"go-venice/internal/api/dto"
)

type NodeService struct {
	node nodes.Noder
}

func NewNodeService() *NodeService {
	config := configs.NewEnvConfig()
	node := nodes.NewEvm(config.RpcUrl, config.ChainId)
	return &NodeService{node}
}

func (s *NodeService) GetBalance(address string) (*dto.BalanceRequest, error) {
	balance, err := s.node.GetBalance(address)
	if err != nil {
		return nil, err
	}

	return &dto.BalanceRequest{Amount: balance.String()}, nil
}
