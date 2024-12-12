package static

import (
	"net/http"

	"github.com/DevinMox/Mox-Cloud/public"
	"github.com/gin-gonic/gin"
)

type CompleteUploadRequest struct {
	Filename    string `json:"filename"`
	Hash        string `json:"hash"`
	TotalChunks int    `json:"totalChunks"`
}

func HandleChunkUpload(c *gin.Context) {
	file, err := c.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未收到分片文件"})
		return
	}

	index := c.PostForm("index")
	hash := c.PostForm("hash")
	filename := c.PostForm("filename")

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开上传的文件失败"})
		return
	}
	defer src.Close()

	info := public.ChunkInfo{
		Index:    index,
		Hash:     hash,
		Filename: filename,
	}

	if err := public.HandleChunkUpload(src, info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存分片文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分片上传成功"})
}

func HandleUploadComplete(c *gin.Context) {
	var req CompleteUploadRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求无效"})
		return
	}

	if err := public.MergeChunks(req.Filename, req.Hash, req.TotalChunks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "合并分片文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
}
