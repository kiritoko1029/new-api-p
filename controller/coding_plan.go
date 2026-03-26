package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service/coding_plan"
	"github.com/gin-gonic/gin"
)

// GetChannelCodingPlanUsage 获取单个渠道的编程套餐用量
func GetChannelCodingPlanUsage(c *gin.Context) {
	channelId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的渠道 ID",
		})
		return
	}

	channel, err := model.GetChannelById(channelId, false)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "渠道不存在",
		})
		return
	}

	if !coding_plan.IsCodingPlanChannel(channel.Type) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "该渠道不是编程套餐类型",
		})
		return
	}

	// 如果有缓存的用量数据，直接返回
	if channel.CodingPlanUsage != "" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data":    channel.CodingPlanUsage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    "{}",
	})
}

// RefreshChannelCodingPlanUsage 刷新单个渠道的编程套餐用量
func RefreshChannelCodingPlanUsage(c *gin.Context) {
	channelId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的渠道 ID",
		})
		return
	}

	usage, err := coding_plan.RefreshChannelUsage(channelId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "刷新失败: " + err.Error(),
		})
		return
	}

	data, _ := coding_plan.MarshalUsage(usage)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
}

// BatchRefreshCodingPlanUsage 批量刷新所有编程套餐渠道用量
func BatchRefreshCodingPlanUsage(c *gin.Context) {
	go coding_plan.RefreshAllUsage()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已开始批量刷新",
	})
}
