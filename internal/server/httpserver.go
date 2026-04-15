package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	handler "github.com/lichenglife/easyblog/internal/apiserver/handler/http"
	"github.com/lichenglife/easyblog/internal/app"
	"github.com/lichenglife/easyblog/internal/pkg/auth"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// HttpServer http 服务
type HTTPServer struct {
	config  *viper.Viper
	app     app.IApp
	engine  *gin.Engine
	http    *http.Server
	handler handler.Handler
}

func NewHttpServer(config *viper.Viper, app app.IApp) (*HTTPServer, error) {
	server := &HTTPServer{
		config: config,
		app:    app,
	}

	factory := app.GetStoreFactory()

	// 创建 JWT 实例
	jwt := auth.NewJWT(config)

	// 创建业务处理器 handler
	handler := handler.NewHandler(app.GetLogger(), factory, jwt)

	server.handler = handler

	return server, nil
}

// Init 初始化 httpServer
func (s *HTTPServer) Init() error {
	err := s.initEngine()
	if err != nil {
		return fmt.Errorf("初始化 Engine 失败:%v", err)
	}
	if err := s.registerRoutes(); err != nil {
		return fmt.Errorf("注册路由规则失败:%v", err)
	}
	if err := s.initHttpServer(); err != nil {
		return fmt.Errorf("初始化 httpserver 失败:%v", err)
	}
	return nil
}

func (s *HTTPServer) initEngine() error {
	gin.SetMode(s.config.GetString("server.http.mode"))
	engine := gin.New()

	engine.Use(
		middleware.RequestID(),
		middleware.Logger(s.app.GetLogger()),
		middleware.Recovery(s.app.GetLogger()),
		middleware.CORS(),
	)

	// 添加限流中间件（如果启用）
	if s.config.GetBool("rateLimit.enabled") {
		rate := s.config.GetInt("rateLimit.rate")
		burst := s.config.GetInt("rateLimit.burst")
		engine.Use(middleware.RateLimitByIP(rate, burst))
	}

	s.engine = engine

	return nil
}

func (s *HTTPServer) registerRoutes() error {
	s.engine.GET("/healthz", s.healthcheck)
	s.engine.GET("/healthcheck", s.healthcheck)

	v1 := s.engine.Group("/v1")
	{
		// 公开接口（无需认证）
		v1.POST("/user", s.handler.Users().CreateUser)
		v1.POST("/user/login", s.handler.Users().UserLogin)

		// 需要认证的接口
		authGroup := v1.Group("")
		authGroup.Use(middleware.Auth(s.getJWT(), s.app.GetLogger()))
		{
			// 用户接口
			authGroup.GET("/user/me", s.handler.Users().GetUserInfo)
			authGroup.GET("/user/info", s.handler.Users().UserInfo)
			authGroup.POST("/user/logout", s.handler.Users().UserLogout)
			authGroup.GET("/user/list", s.handler.Users().ListUsers)
			authGroup.GET("/user/:id", s.handler.Users().GetUserByID)
			authGroup.PUT("/user/:username", s.handler.Users().UpdateUser)
			authGroup.DELETE("/user/:id", s.handler.Users().DeleteUser)
			authGroup.POST("/user/change-password", s.handler.Users().ChangePassword)
			authGroup.PUT("/user/profile", s.handler.Users().UpdateProfile)

			// 管理员接口（需要管理员角色）
			admin := authGroup.Group("")
			admin.Use(middleware.RBAC(2)) // 2-管理员角色
			{
				admin.POST("/admin/user/audit", s.handler.Users().AuditUser)
				admin.POST("/admin/user/ban", s.handler.Users().BanUser)
				// 分类管理
				admin.POST("/category", s.handler.Categories().CreateCategory)
				admin.PUT("/category", s.handler.Categories().UpdateCategory)
				admin.DELETE("/category/:id", s.handler.Categories().DeleteCategory)
				// 标签管理
				admin.POST("/tag", s.handler.Tags().CreateTag)
				admin.PUT("/tag", s.handler.Tags().UpdateTag)
				admin.DELETE("/tag/:id", s.handler.Tags().DeleteTag)
				// 文章管理
				admin.POST("/admin/post/:id/top", s.handler.Posts().TopPost)
				admin.POST("/admin/post/:id/untop", s.handler.Posts().UnTopPost)
			}

			// 公开读取接口（需要认证）
			// 分类接口
			authGroup.GET("/categories/tree", s.handler.Categories().GetCategoryTree)
			authGroup.GET("/categories", s.handler.Categories().ListCategories)
			authGroup.GET("/category/:id", s.handler.Categories().GetCategory)
			// 标签接口
			authGroup.GET("/tags", s.handler.Tags().ListTags)
			authGroup.GET("/tag/:id", s.handler.Tags().GetTag)

			// 文章接口
			authGroup.POST("/post", s.handler.Posts().CreatePost)
			authGroup.GET("/post/:id", s.handler.Posts().GetPostByID)
			authGroup.GET("/post/list", s.handler.Posts().ListPosts)
			authGroup.PUT("/post/:id", s.handler.Posts().UpdatePost)
			authGroup.DELETE("/post/:id", s.handler.Posts().DeletePost)
			authGroup.GET("/post/user/:userID", s.handler.Posts().GetPostsByUserID)
			// 文章搜索接口
			authGroup.GET("/posts/search", s.handler.Posts().SearchPosts)
			// 我的文章接口
			authGroup.GET("/posts/me", s.handler.Posts().GetMyPosts)

			// 评论接口
			authGroup.POST("/comment", s.handler.Comments().CreateComment)
			authGroup.GET("/comment/post/:postID", s.handler.Comments().ListComments)
			authGroup.GET("/comment/reply/:parentID", s.handler.Comments().ListReplies)
			authGroup.GET("/comment/:id", s.handler.Comments().GetCommentByID)
			authGroup.PUT("/comment/:id", s.handler.Comments().UpdateComment)
			authGroup.DELETE("/comment/:id", s.handler.Comments().DeleteComment)
			authGroup.POST("/comment/:id/like", s.handler.Comments().LikeComment)
			authGroup.POST("/comment/:id/unlike", s.handler.Comments().UnlikeComment)

			// 点赞接口
			authGroup.POST("/like/post/:postID", s.handler.Likes().LikePost)
			authGroup.POST("/like/post/:postID/unlike", s.handler.Likes().UnlikePost)
			authGroup.GET("/like/post/:postID/status", s.handler.Likes().GetLikeStatus)
		}

		// 图片上传接口（公开）
		v1.POST("/image/upload", s.handler.Images().UploadImage)
		v1.GET("/images", s.handler.Images().ListImages)
		v1.GET("/images/:filename", s.handler.Images().GetImage)
		v1.DELETE("/image/:filename", s.handler.Images().DeleteImage)
	}

	return nil
}

// getJWT 获取 JWT 实例
func (s *HTTPServer) getJWT() *auth.JWT {
	return auth.NewJWT(s.config)
}

func (s *HTTPServer) healthcheck(c *gin.Context) {
	core.WriteResponse(c, nil, gin.H{"status": "OK"})
}

func (s *HTTPServer) initHttpServer() error {
	s.http = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.GetString("server.http.host"), s.config.GetInt("server.http.port")),
		Handler:        s.engine,
		ReadTimeout:    time.Duration(s.config.GetInt("server.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(s.config.GetInt("server.http.write_timeout")) * time.Second,
		MaxHeaderBytes: s.config.GetInt("server.http.max_header_bytes"),
	}
	return nil
}

func (s *HTTPServer) Start() error {
	logger := s.app.GetLogger().Logger
	logger.Info("启动 HTTPServer", zap.String("addr", s.http.Addr), zap.String("mode", s.config.GetString("server.http.mode")))

	go func() {
		if err := s.http.ListenAndServe(); err != nil {
			s.app.GetLogger().Error("启动 HTTP 服务器失败", zap.Error(err))
		}
	}()
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.app.GetLogger().Logger.Info("正在停止 HTTP 服务器...")
	if err := s.http.Shutdown(ctx); err != nil {
		s.app.GetLogger().Logger.Error("HTTP 服务器停止失败", zap.Error(err))
	}
	return nil
}
