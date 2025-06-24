package grpc

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/pkg/log"
	pb "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HealthHandler 定义健康检查的grpc请求处理接口
type HealthzHandler interface {
	Healthz(ctx context.Context, req *emptypb.Empty) (*pb.HealthzResponse, error)
}

type healthHandler struct {
}

// Healthz implements HealthzHandler.
func (h *healthHandler) Healthz(ctx context.Context, req *emptypb.Empty) (*pb.HealthzResponse, error) {
	log.Log.Info("健康检查。。")
	return &pb.HealthzResponse{
		Status:    pb.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "服务正常运行",
	}, nil
}

func NewHealzHandler(log *log.Logger) HealthzHandler {
	return &healthHandler{}
}
