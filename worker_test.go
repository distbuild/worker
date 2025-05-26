package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"distbuild/boong/worker/proto"
)

const (
	testWorkerDir  = "test-worker"
	testWorkerFile = "test_worker.c"
	testWorkerData = "int main() {return 0;}"
)

func initWorkerTest() server {
	return server{}
}

func TestValidateBuildFile(t *testing.T) {
	s := initWorkerTest()

	workSpacePath, _ = os.MkdirTemp("", testWorkerDir)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(workSpacePath)

	base := "/remote"

	hash := sha256.New()
	r := strings.NewReader(testWorkerData)
	_, _ = io.Copy(hash, r)

	buildFiles := []*proto.BuildFile{
		{
			FilePath: filepath.Join(base, testWorkerFile),
			FileData: []byte(testWorkerData),
			CheckSum: hex.EncodeToString(hash.Sum(nil)),
		},
	}

	err := s.validateBuildFile(buildFiles)
	assert.Equal(t, nil, err)
}

func TestSetBuildPath(t *testing.T) {
	s := initWorkerTest()

	workSpacePath = os.TempDir()

	base := "/path"

	targets := []string{
		filepath.Join(base, testWorkerDir),
	}

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(targets[0])

	err := s.setBuildPath(targets)
	assert.Equal(t, nil, err)
}

func TestRunBuild(t *testing.T) {
	s := initWorkerTest()

	workSpacePath = os.TempDir()
	_ = os.WriteFile(filepath.Join(workSpacePath, testWorkerFile), []byte(testWorkerData), 0644)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(filepath.Join(workSpacePath, testWorkerFile))

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(filepath.Join(workSpacePath, "a.out"))

	dir := "/remote"
	rule := fmt.Sprintf("gcc %s", filepath.Join(dir, testWorkerFile))
	targets := []string{
		filepath.Join(dir, testWorkerFile),
	}
	id := "00:11:22:33:44:55-1737170443"

	_, err := s.runBuild(rule, targets, id)
	assert.Equal(t, nil, err)
}
