package public

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed all:dist
var Public embed.FS

type ChunkInfo struct {
	Index    int    `json:"index"`
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
}

func HandleChunkUpload(chunk io.Reader, info ChunkInfo) error {
	// 创建临时目录存储分片
	tmpDir := filepath.Join("uploads", "tmp", info.Hash)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}

	// 保存分片
	chunkPath := filepath.Join(tmpDir, fmt.Sprintf("%d", info.Index))
	f, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, chunk)
	return err
}

func MergeChunks(filename, hash string, totalChunks int) error {
	tmpDir := filepath.Join("uploads", "tmp", hash)
	targetPath := filepath.Join("uploads", filename)

	// 创建目标文件
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	// 按顺序合并所有分片
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(tmpDir, fmt.Sprintf("%d", i))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(target, chunk)
		chunk.Close()
		if err != nil {
			return err
		}
	}

	// 清理临时文件
	return os.RemoveAll(tmpDir)
}
