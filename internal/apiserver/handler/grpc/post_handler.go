package grpc

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	pb "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
)

// PostHandler 定义处理博客相关的grpc请求的接口
type PostHandler interface {
	// CreatePost 创建文章
	CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error)
	// UpdatePost 更新文章
	UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error)
	// DeletePost 删除文章
	DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error)
	// GetPost 获取文章信息
	GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error)
	// ListPost 列出所有文章
	ListPost(ctx context.Context, req *pb.ListPostRequest) (*pb.ListPostResponse, error)
}

// postHandler 作为PostHandler的接口实现

type postHandler struct {
	biz biz.IBiz
}

// CreatePost implements PostHandler.
func (p *postHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	panic("unimplemented")
}

// DeletePost implements PostHandler.
func (p *postHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	panic("unimplemented")
}

// GetPost implements PostHandler.
func (p *postHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	panic("unimplemented")
}

// ListPost implements PostHandler.
func (p *postHandler) ListPost(ctx context.Context, req *pb.ListPostRequest) (*pb.ListPostResponse, error) {
	panic("unimplemented")
}

// UpdatePost implements PostHandler.
func (p *postHandler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	panic("unimplemented")
}

// NewPostHandler
func NewPostHandler(biz biz.IBiz) PostHandler {
	return &postHandler{
		biz: biz,
	}
}
