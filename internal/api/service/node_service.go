package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"go-venice/configs"
	"go-venice/internal/adapter/nodes"
	"go-venice/internal/api/dto"
)

type NodeService interface {
	GetBalance(ctx context.Context, address string) (*dto.BalanceResponse, error)
}

type nodeService struct {
	config *configs.EnvConfig
	node   nodes.Node
}

func NewNodeService(config *configs.EnvConfig, node nodes.Node) NodeService {
	return &nodeService{config, node}
}

func (s *nodeService) GetBalance(ctx context.Context, address string) (*dto.BalanceResponse, error) {
	balance, err := s.node.GetBalance(ctx, address)
	if err != nil {
		return nil, err
	}

	return &dto.BalanceResponse{Amount: balance.String(), Symbol: "BASE"}, nil
}

func (s *nodeService) CreateTransaction(ctx context.Context, fromAddress string, toAddress string, amount string) (*dto.TransactionResponse, error) {
	reqID, _ := ctx.Value("request_id").(string)
	log.Info().
		Str("request", reqID).
		Str("from", fromAddress).
		Str("to", toAddress).
		Str("amount", amount).
		Msg("Creating transaction")

	unsigned, err := s.node.CreateTransferTransaction(fromAddress, toAddress, amount)
	if err != nil {
		return nil, err
	}
	return &dto.TransactionResponse{Tx: unsigned}, nil
}
