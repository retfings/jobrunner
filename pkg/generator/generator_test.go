package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetTodayPath(t *testing.T) {
	gen := NewGenerator("/root/code/jobs")

	path := gen.GetTodayPath()

	// 验证路径格式
	if path == "" {
		t.Error("GetTodayPath() 返回空路径")
	}

	// 验证路径包含基础路径
	if !filepath.HasPrefix(path, gen.basePath) {
		t.Errorf("GetTodayPath() 路径不以 basePath 开头：%s", path)
	}
}

func TestEnsureBasePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jobrunner-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败：%v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := NewGenerator(filepath.Join(tmpDir, "jobs"))
	if err := gen.EnsureBasePath(); err != nil {
		t.Fatalf("EnsureBasePath() 失败：%v", err)
	}

	if _, err := os.Stat(gen.basePath); os.IsNotExist(err) {
		t.Error("基础路径未创建")
	}
}

func TestCreateNext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jobrunner-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败：%v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := NewGenerator(filepath.Join(tmpDir, "jobs"))

	// 创建第一个目录
	path1, err := gen.CreateNext()
	if err != nil {
		t.Fatalf("CreateNext() 失败：%v", err)
	}

	// 验证目录存在
	if _, err := os.Stat(path1); os.IsNotExist(err) {
		t.Error("目录未创建")
	}

	// 创建第二个目录
	path2, err := gen.CreateNext()
	if err != nil {
		t.Fatalf("CreateNext() 失败：%v", err)
	}

	if path1 >= path2 {
		t.Errorf("第二个目录应该大于第一个：path1=%s, path2=%s", path1, path2)
	}
}

func TestIsEmpty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jobrunner-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败：%v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := NewGenerator(tmpDir)

	// 测试空目录
	emptyDir := filepath.Join(tmpDir, "empty")
	os.MkdirAll(emptyDir, 0755)

	isEmpty, err := gen.IsEmpty(emptyDir)
	if err != nil {
		t.Fatalf("IsEmpty() 失败：%v", err)
	}
	if !isEmpty {
		t.Error("空目录应该返回 true")
	}

	// 测试非空目录
	nonEmptyDir := filepath.Join(tmpDir, "nonempty")
	os.MkdirAll(nonEmptyDir, 0755)
	os.WriteFile(filepath.Join(nonEmptyDir, "file.txt"), []byte("test"), 0644)

	isEmpty, err = gen.IsEmpty(nonEmptyDir)
	if err != nil {
		t.Fatalf("IsEmpty() 失败：%v", err)
	}
	if isEmpty {
		t.Error("非空目录应该返回 false")
	}

	// 测试不存在的目录
	isEmpty, err = gen.IsEmpty(filepath.Join(tmpDir, "notexist"))
	if err != nil {
		t.Fatalf("IsEmpty() 失败：%v", err)
	}
	if !isEmpty {
		t.Error("不存在的目录应该返回 true")
	}
}

func TestGetNextSequence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "jobrunner-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败：%v", err)
	}
	defer os.RemoveAll(tmpDir)

	gen := NewGenerator(filepath.Join(tmpDir, "jobs"))

	// 创建一些测试目录
	todayPath := gen.GetTodayPath()
	os.MkdirAll(filepath.Join(todayPath, "001"), 0755)
	os.MkdirAll(filepath.Join(todayPath, "002"), 0755)
	os.MkdirAll(filepath.Join(todayPath, "005"), 0755)

	// 获取下一个序号
	nextSeq, err := gen.getNextSequence(todayPath)
	if err != nil {
		t.Fatalf("getNextSequence() 失败：%v", err)
	}

	if nextSeq != 6 {
		t.Errorf("getNextSequence() = %d, 期望 6", nextSeq)
	}
}
