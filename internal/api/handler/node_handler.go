package handler

import (
	"github.com/gin-gonic/gin"
	"go-venice/internal/api/service"
	"net/http"
)

type NodeHandler interface {
	GetBalance(ctx *gin.Context)
}
type nodeHandler struct {
	svc service.NodeService
}

func NewNodeHandler(svc service.NodeService) NodeHandler {
	return &nodeHandler{svc: svc}
}

func (nh *nodeHandler) GetBalance(ctx *gin.Context) {
	// 예시: 쿼리 파라미터 받아서 서비스 호출
	address := ctx.Query("address")
	balance, err := nh.svc.GetBalance(ctx, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, balance)
}
