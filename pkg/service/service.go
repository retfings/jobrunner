package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kardianos/service"
)

// Service systemd 服务管理
type Service struct {
	svc     service.Service
	name    string
	display string
	desc    string
}

// Config 服务配置
type Config struct {
	Name        string
	DisplayName string
	Description string
	ConfigPath  string
	Executable  string
	WorkingDir  string
}

// New 创建服务
func New(cfg Config) (*Service, error) {
	svcConfig := &service.Config{
		Name:             cfg.Name,
		DisplayName:      cfg.DisplayName,
		Description:      cfg.Description,
		Executable:       cfg.Executable,
		WorkingDirectory: cfg.WorkingDir,
		Arguments:        []string{"--config", cfg.ConfigPath},
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("创建服务失败：%w", err)
	}

	return &Service{
		svc:  svc,
		name: cfg.Name,
	}, nil
}

// program 实现 service.Interface
type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// 实际的运行逻辑在 main.go 中处理
	// 这里只是服务框架
	log.Println("服务运行中...")
}

func (p *program) Stop(s service.Service) error {
	log.Println("服务停止...")
	return nil
}

// Install 安装服务
func (s *Service) Install() error {
	if err := s.svc.Install(); err != nil {
		return fmt.Errorf("安装服务失败：%w", err)
	}
	log.Printf("服务 %s 安装成功", s.name)
	return nil
}

// Start 启动服务
func (s *Service) Start() error {
	if err := s.svc.Start(); err != nil {
		return fmt.Errorf("启动服务失败：%w", err)
	}
	log.Printf("服务 %s 已启动", s.name)
	return nil
}

// Stop 停止服务
func (s *Service) Stop() error {
	if err := s.svc.Stop(); err != nil {
		return fmt.Errorf("停止服务失败：%w", err)
	}
	log.Printf("服务 %s 已停止", s.name)
	return nil
}

// Restart 重启服务
func (s *Service) Restart() error {
	if err := s.svc.Restart(); err != nil {
		return fmt.Errorf("重启服务失败：%w", err)
	}
	log.Printf("服务 %s 已重启", s.name)
	return nil
}

// Uninstall 卸载服务
func (s *Service) Uninstall() error {
	if err := s.svc.Uninstall(); err != nil {
		return fmt.Errorf("卸载服务失败：%w", err)
	}
	log.Printf("服务 %s 已卸载", s.name)
	return nil
}

// Status 获取服务状态
func (s *Service) Status() (service.Status, error) {
	return s.svc.Status()
}

// GetSystemdUnitPath 获取 systemd unit 文件路径
func (s *Service) GetSystemdUnitPath() (string, error) {
	execPath, err := exec.LookPath("systemctl")
	if err != nil {
		return "", fmt.Errorf("找不到 systemctl: %w", err)
	}
	_ = execPath
	// systemd 服务文件通常安装在这里
	return filepath.Join("/etc", "systemd", "system", s.name+".service"), nil
}

// Run 运行服务（用于调试或前台运行）
func (s *Service) Run() error {
	return s.svc.Run()
}

// FindExecutable 查找可执行文件路径
func FindExecutable(name string) (string, error) {
	// 首先尝试当前目录
	execPath, err := os.Executable()
	if err == nil {
		return execPath, nil
	}

	// 尝试 PATH
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("找不到可执行文件：%w", err)
	}

	return path, nil
}
