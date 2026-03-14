package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"jobrunner/pkg/config"
	"jobrunner/pkg/generator"
	"jobrunner/pkg/monitor"
	"jobrunner/pkg/service"
)

var version = "dev"

func main() {
	// 命令行参数
	var (
		configPath  = flag.String("config", "", "配置文件路径")
		showVersion = flag.Bool("version", false, "显示版本")
		command     = flag.String("command", "run", "命令：run/install/start/stop/restart/uninstall")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("jobrunner version: %s\n", version)
		os.Exit(0)
	}

	// 处理服务命令
	if *command != "run" {
		handleServiceCommand(*command, *configPath)
		return
	}

	// 运行服务
	run(*configPath)
}

func run(configPath string) {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("jobrunner 启动...")

	// 加载配置
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败：%v", err)
	}
	log.Printf("配置文件：%s", configPath)
	log.Printf("基础路径：%s", cfg.BasePath)
	log.Printf("检查间隔：%v", cfg.CheckInterval)

	// 确保基础路径存在
	gen := generator.NewGenerator(cfg.BasePath)
	if err := gen.EnsureBasePath(); err != nil {
		log.Fatalf("创建基础路径失败：%v", err)
	}

	// 创建监控器
	mon := monitor.NewMonitor(gen, cfg.CheckInterval)
	mon.SetOnCreate(func(path string) {
		log.Printf("已创建目录：%s", path)
	})

	// 初始延迟
	if cfg.InitialDelay > 0 {
		log.Printf("等待 %v 后创建初始目录...", cfg.InitialDelay)
		time.Sleep(cfg.InitialDelay)
	}

	// 启动监控
	if err := mon.Start(); err != nil {
		log.Fatalf("启动监控失败：%v", err)
	}

	// 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("收到信号：%v，正在关闭...", sig)

	if err := mon.Stop(); err != nil {
		log.Printf("停止监控失败：%v", err)
	}

	log.Println("jobrunner 已关闭")
}

func loadConfig(path string) (*config.Config, error) {
	// 如果指定了配置路径，使用指定路径
	if path != "" {
		return config.Load(path)
	}

	// 尝试默认路径
	defaultPaths := []string{
		"/etc/jobrunner/config.yaml",
		"./config/config.yaml",
		"config/config.yaml",
	}

	for _, p := range defaultPaths {
		if _, err := os.Stat(p); err == nil {
			return config.Load(p)
		}
	}

	// 都没有则返回默认配置
	log.Println("未找到配置文件，使用默认配置")
	return config.DefaultConfig(), nil
}

func handleServiceCommand(cmd, configPath string) {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("获取可执行文件路径失败：%v", err)
	}

	// 如果未指定配置文件，使用默认路径
	if configPath == "" {
		configPath = "/etc/jobrunner/config.yaml"
	}

	cfg := service.Config{
		Name:        "jobrunner",
		DisplayName: "Job Runner",
		Description: "定时目录生成器 - 按固定间隔检查并创建序号目录",
		ConfigPath:  configPath,
		Executable:  execPath,
		WorkingDir:  filepath.Dir(execPath),
	}

	svc, err := service.New(cfg)
	if err != nil {
		log.Fatalf("创建服务失败：%v", err)
	}

	switch cmd {
	case "install":
		if err := svc.Install(); err != nil {
			log.Fatalf("安装失败：%v", err)
		}
		log.Println("安装成功，使用 'systemctl start jobrunner' 启动服务")
	case "start":
		if err := svc.Start(); err != nil {
			log.Fatalf("启动失败：%v", err)
		}
		log.Println("服务已启动")
	case "stop":
		if err := svc.Stop(); err != nil {
			log.Fatalf("停止失败：%v", err)
		}
		log.Println("服务已停止")
	case "restart":
		if err := svc.Restart(); err != nil {
			log.Fatalf("重启失败：%v", err)
		}
		log.Println("服务已重启")
	case "uninstall":
		if err := svc.Uninstall(); err != nil {
			log.Fatalf("卸载失败：%v", err)
		}
		log.Println("服务已卸载")
	default:
		log.Fatalf("未知命令：%s", cmd)
	}
}
