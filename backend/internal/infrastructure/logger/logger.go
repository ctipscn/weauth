package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// LogService 基于文件的日志服务
type LogService struct {
	logger *log.Logger
	file   *os.File
	mu     sync.Mutex
}

// NewLogService 创建文件日志服务
func NewLogService(filePath string) (*LogService, error) {
	// 确保日志目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("[LOGGER] Failed to create log directory %s: %w", dir, err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("[LOGGER] Failed to open log file %s: %w", filePath, err)
	}

	// 输出到文件和控制台
	multiWriter := io.MultiWriter(file)
	logger := log.New(multiWriter, "", log.Ldate|log.Ltime)

	return &LogService{
		logger: logger,
		file:   file,
	}, nil
}

// Info 记录 INFO 级别日志
func (l *LogService) Info(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[INFO] %s", msg)
}

// Infof 格式化记录 INFO 级别日志
func (l *LogService) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn 记录 WARN 级别日志
func (l *LogService) Warn(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[WARN] %s", msg)
}

// Warnf 格式化记录 WARN 级别日志
func (l *LogService) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Error 记录 ERROR 级别日志
func (l *LogService) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[ERROR] %s", msg)
}

// Errorf 格式化记录 ERROR 级别日志
func (l *LogService) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

// Close 关闭日志文件
func (l *LogService) Close() error {
	return l.file.Close()
}