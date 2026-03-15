package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 配置文件结构
type Config struct {
	BasePath      string        `yaml:"base_path"`
	CheckInterval time.Duration `yaml:"check_interval"`
	InitialDelay  time.Duration `yaml:"initial_delay"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BasePath:      "/root/code/jobs",
		CheckInterval: 30 * time.Second,
		InitialDelay:  5 * time.Second,
	}
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败：%w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 验证配置有效性
func (c *Config) Validate() error {
	if c.BasePath == "" {
		return fmt.Errorf("base_path 不能为空")
	}
	if c.CheckInterval <= 0 {
		return fmt.Errorf("check_interval 必须大于 0")
	}
	if c.InitialDelay < 0 {
		return fmt.Errorf("initial_delay 不能为负数")
	}
	return nil
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("序列化配置失败：%w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败：%w", err)
	}

	return nil
}
