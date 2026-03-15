package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"jobrunner/pkg/generator"
)

// Monitor 目录监控器
type Monitor struct {
	gen      *generator.Generator
	interval time.Duration
	onCreate func(path string)
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewMonitor 创建新的监控器
func NewMonitor(gen *generator.Generator, interval time.Duration) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		gen:      gen,
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// SetOnCreate 设置目录创建回调
func (m *Monitor) SetOnCreate(fn func(path string)) {
	m.onCreate = fn
}

// Start 启动监控
func (m *Monitor) Start() error {
	log.Println("启动目录监控器...")

	// 创建初始目录
	initialPath, err := m.gen.CreateNext()
	if err != nil {
		return fmt.Errorf("创建初始目录失败：%w", err)
	}
	log.Printf("创建初始目录：%s", initialPath)

	if m.onCreate != nil {
		m.onCreate(initialPath)
	}

	// 启动监控循环
	go m.run()

	return nil
}

// Stop 停止监控
func (m *Monitor) Stop() error {
	log.Println("停止目录监控器...")
	m.cancel()
	return nil
}

// run 监控循环
func (m *Monitor) run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.check()
		}
	}
}

// check 检查当前目录状态
func (m *Monitor) check() {
	currentPath, err := m.gen.GetCurrentPath()
	if err != nil {
		log.Printf("获取当前目录失败：%v", err)
		return
	}

	if currentPath == "" {
		log.Println("今天还没有创建目录")
		return
	}

	isEmpty, err := m.gen.IsEmpty(currentPath)
	if err != nil {
		log.Printf("检查目录状态失败：%v", err)
		return
	}

	log.Printf("检查目录：%s, 为空：%v", currentPath, isEmpty)

	if !isEmpty {
		// 目录不为空，创建下一个
		newPath, err := m.gen.CreateNext()
		if err != nil {
			log.Printf("创建新目录失败：%v", err)
			return
		}
		log.Printf("目录不为空，创建新目录：%s", newPath)

		if m.onCreate != nil {
			m.onCreate(newPath)
		}
	}
}
