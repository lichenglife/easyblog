package core

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
)

// Response 定义了API响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误信息
	Data    interface{} `json:"data"`    // 响应数据
}

// ListResponse 定义了列表类API的响应结构
type ListResponse[T any] struct {
	TotalCount int64 `json:"totalCount"` // 总记录数
	HasMore    bool  `json:"hasMore"`    // 是否还有更多
	Items      []T   `json:"items"`      // 数据项
}

// WriteResponse 写入HTTP响应
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		// 解码错误信息
		e := errno.Decode(err)
		// 返回错误响应
		c.JSON(e.HTTP(), Response{
			Code:    e.Code(),
			Message: e.Message(),
			Data:    nil,
		})
		c.Abort()
	} else {
		// 返回成功响应
		c.JSON(http.StatusOK, Response{
			Code:    errno.OK.Code(),
			Message: errno.OK.Message(),
			Data:    data,
		})
	}
}

// WriteListResponse 写入列表类API的HTTP响应
func WriteListResponse[T any](c *gin.Context, total int64, page, pageSize int, items []T) {
	WriteResponse(c, nil, ListResponse[T]{
		TotalCount: total,
		HasMore:    total > int64(page*pageSize),
		Items:      items,
	})
}

// TODO 提取到utils中
// GetPageParam 获取分页参数
func GetPageParam(c *gin.Context) int {
	// 获取请求中的page参数，默认为1
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// GetLimitParam 获取每页条数参数
func GetLimitParam(c *gin.Context) int {
	// 获取请求中的limit参数，默认为10
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

// GetOffset 获取数据库查询的offset
func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
