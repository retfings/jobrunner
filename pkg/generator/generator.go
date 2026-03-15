package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Generator 目录生成器
type Generator struct {
	basePath string
	mu       sync.Mutex
}

// NewGenerator 创建新的目录生成器
func NewGenerator(basePath string) *Generator {
	return &Generator{
		basePath: basePath,
	}
}

// GetTodayPath 获取今天的日期路径 (不含序号)
func (g *Generator) GetTodayPath() string {
	now := time.Now()
	return filepath.Join(g.basePath,
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
}

// GetNextPath 获取下一个目录路径
func (g *Generator) GetNextPath() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	todayPath := g.GetTodayPath()
	nextSeq, err := g.getNextSequence(todayPath)
	if err != nil {
		return "", err
	}

	return filepath.Join(todayPath, fmt.Sprintf("%03d", nextSeq)), nil
}

// CreateNext 创建下一个目录
func (g *Generator) CreateNext() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	todayPath := g.GetTodayPath()
	nextSeq, err := g.getNextSequence(todayPath)
	if err != nil {
		return "", err
	}

	newPath := filepath.Join(todayPath, fmt.Sprintf("%03d", nextSeq))

	// 创建目录
	if err := os.MkdirAll(newPath, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败：%w", err)
	}

	return newPath, nil
}

// getNextSequence 获取下一个序号
func (g *Generator) getNextSequence(todayPath string) (int, error) {
	// 确保今天的基础路径存在
	if err := os.MkdirAll(todayPath, 0755); err != nil {
		return 0, fmt.Errorf("创建日期目录失败：%w", err)
	}

	// 读取今天的目录列表
	entries, err := os.ReadDir(todayPath)
	if err != nil {
		return 0, fmt.Errorf("读取目录列表失败：%w", err)
	}

	// 正则匹配 3 位序号目录
	seqPattern := regexp.MustCompile(`^\d{3}$`)
	maxSeq := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if seqPattern.MatchString(name) {
			seq, err := strconv.Atoi(name)
			if err != nil {
				continue
			}
			if seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	return maxSeq + 1, nil
}

// GetCurrentPath 获取当前序号目录（最大的序号目录）
func (g *Generator) GetCurrentPath() (string, error) {
	todayPath := g.GetTodayPath()

	entries, err := os.ReadDir(todayPath)
	if err != nil {
		return "", fmt.Errorf("读取目录列表失败：%w", err)
	}

	seqPattern := regexp.MustCompile(`^\d{3}$`)
	maxSeq := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if seqPattern.MatchString(name) {
			seq, err := strconv.Atoi(name)
			if err != nil {
				continue
			}
			if seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	if maxSeq == 0 {
		return "", nil
	}

	return filepath.Join(todayPath, fmt.Sprintf("%03d", maxSeq)), nil
}

// IsEmpty 检查目录是否为空
func (g *Generator) IsEmpty(path string) (bool, error) {
	if path == "" {
		return true, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("读取目录失败：%w", err)
	}

	return len(entries) == 0, nil
}

// EnsureBasePath 确保基础路径存在
func (g *Generator) EnsureBasePath() error {
	return os.MkdirAll(g.basePath, 0755)
}
