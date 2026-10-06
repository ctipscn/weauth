package server

import (
	"context"
	"net"
	"net/http"

	"github.com/ctipscn/weauth/backend/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

// GinServer 封装 Gin HTTP 引擎、http.Server 及其生命周期。
//
// 使用方式：
//
//	gs := server.New("127.0.0.1:8080")
//	route.RouteRegister(gs.Engine(), ...)
//	go gs.Start()
//	// shutdown:
//	gs.Shutdown(ctx)
type GinServer struct {
	engine *gin.Engine
	srv    *http.Server
	addr   string
}

// NewGinServer 创建 GinServer。
// addr 为监听地址（例如 ":8080" 或 "127.0.0.1:8080"）。
func NewGinServer(addr string) *GinServer {
	r := gin.New()
	return &GinServer{engine: r, addr: addr}
}

// Engine 返回内部的 Gin 引擎，供外部注册路由。
func (s *GinServer) Engine() *gin.Engine {
	return s.engine
}

// Addr 返回监听地址。
func (s *GinServer) Addr() string {
	return s.addr
}

// Start 创建 http.Server 并开始监听。通常应在 goroutine 中调用。
func (s *GinServer) Start() error {
	s.srv = &http.Server{
		Addr:    s.addr,
		Handler: s.engine,
	}
	return s.srv.ListenAndServe()
}

// Serve 使用外部创建的 net.Listener 启动 HTTP 服务。
// 调用者创建 Listener 确认端口可用，再启动服务。
func (s *GinServer) Serve(l net.Listener) error {
	s.srv = &http.Server{
		Handler: s.engine,
	}
	return s.srv.Serve(l)
}

// Shutdown 优雅关闭 HTTP 服务器，等待活跃连接处理完毕。
func (s *GinServer) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

// compile-time check: GinServer 实现了 Configurator 接口
var _ infrastructure.Configurator = (*GinServer)(nil)

// ApplyMode 根据运行模式切换 Gin 引擎行为和中间件。
func (s *GinServer) ApplyMode(m infrastructure.Mode) error {
	switch m {
	case infrastructure.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case infrastructure.TestMode:
		gin.SetMode(gin.TestMode)
	default: // DebugMode
		gin.SetMode(gin.DebugMode)
		s.engine.Use(gin.Logger(), gin.Recovery())
	}

	return nil
}
