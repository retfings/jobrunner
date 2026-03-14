package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.BasePath != "/root/code/jobs" {
		t.Errorf("期望 BasePath 为 /root/code/jobs, 得到 %s", cfg.BasePath)
	}
	if cfg.CheckInterval != 30*time.Second {
		t.Errorf("期望 CheckInterval 为 30s, 得到 %v", cfg.CheckInterval)
	}
	if cfg.InitialDelay != 5*time.Second {
		t.Errorf("期望 InitialDelay 为 5s, 得到 %v", cfg.InitialDelay)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "有效配置",
			cfg: &Config{
				BasePath:      "/root/code/jobs",
				CheckInterval: 30 * time.Second,
				InitialDelay:  5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "空 BasePath",
			cfg: &Config{
				BasePath:      "",
				CheckInterval: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "负数 CheckInterval",
			cfg: &Config{
				BasePath:      "/root/code/jobs",
				CheckInterval: -30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "负数 InitialDelay",
			cfg: &Config{
				BasePath:      "/root/code/jobs",
				CheckInterval: 30 * time.Second,
				InitialDelay:  -5 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() 错误 = %v, 期望错误 = %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndSave(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "jobrunner-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败：%v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	// 创建测试配置
	originalCfg := &Config{
		BasePath:      "/tmp/test/jobs",
		CheckInterval: 60 * time.Second,
		InitialDelay:  10 * time.Second,
	}

	// 保存配置
	if err := originalCfg.Save(configPath); err != nil {
		t.Fatalf("保存配置失败：%v", err)
	}

	// 加载配置
	loadedCfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("加载配置失败：%v", err)
	}

	// 验证
	if loadedCfg.BasePath != originalCfg.BasePath {
		t.Errorf("BasePath: 期望 %s, 得到 %s", originalCfg.BasePath, loadedCfg.BasePath)
	}
	if loadedCfg.CheckInterval != originalCfg.CheckInterval {
		t.Errorf("CheckInterval: 期望 %v, 得到 %v", originalCfg.CheckInterval, loadedCfg.CheckInterval)
	}
	if loadedCfg.InitialDelay != originalCfg.InitialDelay {
		t.Errorf("InitialDelay: 期望 %v, 得到 %v", originalCfg.InitialDelay, loadedCfg.InitialDelay)
	}
}

func TestLoadNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("期望加载不存在的文件返回错误")
	}
}
