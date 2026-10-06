// Package bootstrap 提供 Linux 信号处理与优雅关闭能力。
package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// DefaultShutdownTimeout 是所有清理操作的整体超时。
const DefaultShutdownTimeout = 15 * time.Second

// Cleanup 是一个命名清理单元。
type Cleanup struct {
	Name string
	Fn   func(ctx context.Context) error
}

// Shutdown 管理优雅关闭流程。
//
// 收到第一个信号 → 取消 Context，按注册逆序（LIFO）执行所有 Cleanup，整体受超时保护。
// 清理过程中收到第二个信号 → os.Exit(1) 强制退出。
type Shutdown struct {
	ctx      context.Context
	cancel   context.CancelFunc
	cleanups []Cleanup
	sigCh    chan os.Signal
	done     chan struct{}
	mu       sync.Mutex
	timeout  time.Duration
	logf     func(format string, v ...any)
}

// Option 配置 Shutdown。
type Option func(*Shutdown)

// WithTimeout 设置清理超时，默认 DefaultShutdownTimeout。
func WithTimeout(d time.Duration) Option {
	return func(s *Shutdown) {
		s.timeout = d
	}
}

// WithLogger 设置日志函数，默认 log.Printf。
func WithLogger(logf func(format string, v ...any)) Option {
	return func(s *Shutdown) {
		s.logf = logf
	}
}

// WithSignals 自定义监听信号。
func WithSignals(signals ...os.Signal) Option {
	return func(s *Shutdown) {
		signal.Notify(s.sigCh, signals...)
	}
}

// NewShutdown 创建 Shutdown，默认监听 SIGINT 与 SIGTERM。
func NewShutdown(ctx context.Context, opts ...Option) *Shutdown {
	ctx, cancel := context.WithCancel(ctx)
	s := &Shutdown{
		ctx:     ctx,
		cancel:  cancel,
		sigCh:   make(chan os.Signal, 1),
		done:    make(chan struct{}),
		timeout: DefaultShutdownTimeout,
		logf:    log.Printf,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Context 返回一个在收到关闭信号时取消的上下文，用于通知 goroutine 退出。
func (s *Shutdown) Context() context.Context {
	return s.ctx
}

// RegisterCleanup 注册一个 Cleanup。清理时按注册逆序（LIFO）执行。
// 可并发调用。
func (s *Shutdown) RegisterCleanup(c Cleanup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanups = append(s.cleanups, c)
}

// RegisterCleanupFn 包装为 Cleanup 后注册。
func (s *Shutdown) RegisterCleanupFn(name string, fn func(ctx context.Context) error) {
	s.RegisterCleanup(Cleanup{Name: name, Fn: fn})
}

// Wait 阻塞直至收到信号，依次：取消 Context、启动二次信号强制退出监听、LIFO 执行所有 Cleanup。
func (s *Shutdown) Wait() {
	defer close(s.done)

	// 注册默认信号
	signal.Notify(s.sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(s.sigCh)

	sig := <-s.sigCh
	s.logf("[bootstrap] Received signal %v, shutting down", sig)
	s.cancel()

	// 二次信号 → 强制退出（Linux 惯例）
	go func() {
		second := <-s.sigCh
		s.logf("[bootstrap] Received second signal %v, force exiting", second)
		os.Exit(1)
	}()

	s.mu.Lock()
	cleanups := make([]Cleanup, len(s.cleanups))
	copy(cleanups, s.cleanups)
	s.mu.Unlock()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), s.timeout)
	defer shutdownCancel()

	// LIFO 顺序清理
	for i := len(cleanups) - 1; i >= 0; i-- {
		c := cleanups[i]
		s.logf("[bootstrap] Cleanup: %s", c.Name)
		if err := c.Fn(shutdownCtx); err != nil {
			s.logf("[bootstrap] Cleanup %q failed: %v", c.Name, err)
		}
	}

	s.logf("[bootstrap] Shutdown completed")
}

// Done 返回一个在 Wait 返回后关闭的通道。
func (s *Shutdown) Done() <-chan struct{} {
	return s.done
}

// SignalChannel 返回信号通道，供外部直接注册。
func (s *Shutdown) SignalChannel() chan<- os.Signal {
	return s.sigCh
}