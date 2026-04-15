package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// ImageHandler 图片处理接口
type ImageHandler interface {
	// UploadImage 上传图片
	UploadImage(c *gin.Context)
	// GetImage 获取图片
	GetImage(c *gin.Context)
	// ListImages 获取图片列表
	ListImages(c *gin.Context)
	// DeleteImage 删除图片
	DeleteImage(c *gin.Context)
}

type imageHandler struct {
	logger    *log.Logger
	uploadDir string
	maxSize   int64 // 最大上传大小（字节）
}

// NewImageHandler 创建图片处理器的实例
func NewImageHandler(logger *log.Logger, uploadDir string, maxSize int64) ImageHandler {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("创建上传目录失败", zap.Error(err))
	}

	return &imageHandler{
		logger:    logger,
		uploadDir: uploadDir,
		maxSize:   maxSize,
	}
}

// UploadImage 上传图片
func (h *imageHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Warn("获取上传文件失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 检查文件大小
	if file.Size > h.maxSize {
		h.logger.Warn("文件超出最大限制", zap.Int64("size", file.Size), zap.Int64("max", h.maxSize))
		core.WriteResponse(c, errno.ErrFileTooLarge, nil)
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true,
	}
	if !allowedExts[ext] {
		h.logger.Warn("不支持的文件类型", zap.String("ext", ext))
		core.WriteResponse(c, errno.ErrInvalidFileType, nil)
		return
	}

	// 生成唯一的文件名
	timestamp := time.Now().UnixNano()
	random := time.Now().UnixNano() % 10000
	filename := fmt.Sprintf("%d_%d%s", timestamp, random, ext)
	fullPath := filepath.Join(h.uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		h.logger.Error("保存文件失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInternalServer, nil)
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		h.logger.Error("获取文件信息失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInternalServer, nil)
		return
	}

	// 构建访问 URL
	url := "/api/v1/images/" + filename

	h.logger.Info("图片上传成功",
		zap.String("filename", filename),
		zap.Int64("size", fileInfo.Size()),
		zap.String("url", url))

	core.WriteResponse(c, errno.OK, gin.H{
		"filename": filename,
		"url":      url,
		"size":     fileInfo.Size(),
	})
}

// GetImage 获取图片
func (h *imageHandler) GetImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		h.logger.Warn("文件名不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	fullPath := filepath.Join(h.uploadDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		h.logger.Warn("图片不存在", zap.String("filename", filename))
		core.WriteResponse(c, errno.ErrImageNotFound, nil)
		return
	}

	// 根据扩展名设置 Content-Type
	ext := strings.ToLower(filepath.Ext(filename))
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
	}
	contentType, ok := contentTypes[ext]
	if !ok {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存 1 年
	c.File(fullPath)
}

// ListImages 获取图片列表
func (h *imageHandler) ListImages(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 读取上传目录
	entries, err := os.ReadDir(h.uploadDir)
	if err != nil {
		h.logger.Error("读取上传目录失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInternalServer, nil)
		return
	}

	// 过滤图片文件
	var images []gin.H
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	count := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true,
		}[ext] {
			continue
		}

		// 分页
		if count >= startIndex && count < endIndex {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			images = append(images, gin.H{
				"filename": entry.Name(),
				"url":      "/api/v1/images/" + entry.Name(),
				"size":     info.Size(),
				"createdAt": info.ModTime().Format(time.RFC3339),
			})
		}
		count++
	}

	core.WriteResponse(c, errno.OK, gin.H{
		"total":  count,
		"page":   page,
		"size":   pageSize,
		"images": images,
	})
}

// DeleteImage 删除图片
func (h *imageHandler) DeleteImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		h.logger.Warn("文件名不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	fullPath := filepath.Join(h.uploadDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		h.logger.Warn("图片不存在", zap.String("filename", filename))
		core.WriteResponse(c, errno.ErrImageNotFound, nil)
		return
	}

	// 删除文件
	if err := os.Remove(fullPath); err != nil {
		h.logger.Error("删除文件失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInternalServer, nil)
		return
	}

	h.logger.Info("图片删除成功", zap.String("filename", filename))
	core.WriteResponse(c, errno.OK, nil)
}
