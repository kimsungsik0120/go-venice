package service

import (
	"go-venice/configs"
	"go-venice/internal/adapter/nodes"
	"go-venice/internal/api/dto"
)

type NodeService interface {
	GetBalance(address string) (*dto.BalanceResponse, error)
}

type nodeService struct {
	config *configs.EnvConfig
	node   nodes.Noder
}

func NewNodeService(config *configs.EnvConfig, node nodes.Noder) NodeService {
	return &nodeService{config, node}
}

func (s *nodeService) GetBalance(address string) (*dto.BalanceResponse, error) {
	balance, err := s.node.GetBalance(address)
	if err != nil {
		return nil, err
	}

	return &dto.BalanceResponse{Amount: balance.String(), Symbol: "BASE"}, nil
}
