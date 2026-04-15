package core

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
)

// Response 定义了 API 响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误信息
	Data    interface{} `json:"data"`    // 响应数据
}

// ListResponse 定义了列表类 API 的响应结构
type ListResponse[T any] struct {
	TotalCount int64 `json:"totalCount"` // 总记录数
	HasMore    bool  `json:"hasMore"`    // 是否还有更多
	Items      []T   `json:"items"`      // 数据项
}

// WriteResponse 写入 HTTP 响应
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		e := errno.Decode(err)
		c.JSON(e.HTTP(), Response{
			Code:    e.Code(),
			Message: e.Message(),
			Data:    nil,
		})
		return
	}
	// 使用 json.Marshal 确保正确序列化
	resp := Response{
		Code:    errno.OK.Code(),
		Message: errno.OK.Message(),
		Data:    data,
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonBytes)
}

// WriteListResponse 写入列表类 API 的 HTTP 响应
func WriteListResponse[T any](c *gin.Context, total int64, page, pageSize int, items []T) {
	WriteResponse(c, nil, ListResponse[T]{
		TotalCount: total,
		HasMore:    total > int64(page*pageSize),
		Items:      items,
	})
}

// TODO 提取到 utils 中
// GetPageParam 获取分页参数
func GetPageParam(c *gin.Context) int {
	// 获取请求中的 page 参数，默认为 1
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// GetLimitParam 获取每页条数参数
func GetLimitParam(c *gin.Context) int {
	// 获取请求中的 limit 参数，默认为 10
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		return 10
	}

	// 限制最大值
	if limit > 100 {
		return 100
	}

	return limit
}

// GetPaginationParams 获取分页参数
func GetPaginationParams(c *gin.Context) (page, pageSize int) {
	page = GetPageParam(c)
	pageSize = GetLimitParam(c)
	return
}

// GetOffset 获取数据库查询的 offset
func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
