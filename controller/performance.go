package controller

import (
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

// GetModelPerformance 获取模型性能数据
func GetModelPerformance(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	timeRange := c.Query("time_range")

	// If time_range is specified, use it to calculate timestamps
	if timeRange != "" {
		start, end := model.ParseTimeRangeFromString(timeRange)
		if startTimestamp == 0 {
			startTimestamp = start
		}
		if endTimestamp == 0 {
			endTimestamp = end
		}
	}

	// Use default time range if not specified
	if startTimestamp == 0 || endTimestamp == 0 {
		startTimestamp, endTimestamp = model.GetDefaultPerformanceTimeRange()
	}

	performance, err := model.GetModelPerformance(startTimestamp, endTimestamp)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	common.ApiSuccess(c, performance)
}
