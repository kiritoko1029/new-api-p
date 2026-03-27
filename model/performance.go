package model

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// ModelPerformance 模型性能数据
type ModelPerformance struct {
	ModelName     string  `json:"model_name"`
	MedianTTFT    float64 `json:"median_ttft"`     // ms
	MedianTPS     float64 `json:"median_tps"`      // tokens/s
	CombinedScore float64 `json:"combined_score"`  // 综合评分
	SampleCount   int     `json:"sample_count"`    // 样本数量
}

// getJSONExtract returns the appropriate JSON extraction syntax for the current database
func getJSONExtract(field, key string) string {
	if common.UsingPostgreSQL {
		return field + "->>'" + key + "'"
	}
	if common.UsingMySQL {
		return "JSON_EXTRACT(" + field + ", '$." + key + "')"
	}
	// SQLite
	return "json_extract(" + field + ", '$." + key + "')"
}

// GetModelPerformance 查询各模型的性能数据
func GetModelPerformance(startTimestamp, endTimestamp int64) ([]ModelPerformance, error) {
	type logRow struct {
		ModelName        string `gorm:"column:model_name"`
		CompletionTokens int    `gorm:"column:completion_tokens"`
		Other            string `gorm:"column:other"`
	}

	var logs []logRow

	// Query logs with streaming (is_stream = 1) that have completion_tokens > 0
	query := LOG_DB.Table("logs").
		Select("model_name, completion_tokens, other").
		Where("type = ?", LogTypeConsume).
		Where("is_stream = ?", true).
		Where("completion_tokens > ?", 0).
		Where("created_at >= ?", startTimestamp).
		Where("created_at <= ?", endTimestamp)

	if err := query.Scan(&logs).Error; err != nil {
		return nil, err
	}

	// Group by model name
	modelData := make(map[string][]logRow)
	for _, log := range logs {
		modelData[log.ModelName] = append(modelData[log.ModelName], log)
	}

	// Calculate performance for each model
	var results []ModelPerformance
	for modelName, modelLogs := range modelData {
		if len(modelLogs) < 3 {
			// Skip models with too few samples
			continue
		}

		ttfts := make([]float64, 0, len(modelLogs))
		tpsList := make([]float64, 0, len(modelLogs))

		for _, log := range modelLogs {
			frt, tps := parsePerformanceFromOther(log.Other, log.CompletionTokens)
			if frt > 0 {
				ttfts = append(ttfts, frt)
			}
			if tps > 0 {
				tpsList = append(tpsList, tps)
			}
		}

		if len(ttfts) < 3 || len(tpsList) < 3 {
			continue
		}

		sort.Float64s(ttfts)
		sort.Float64s(tpsList)

		medianTTFT := calculateMedian(ttfts)
		medianTPS := calculateMedian(tpsList)

		// Calculate scores
		tpsScore := math.Min(100, medianTPS)
		ttftScore := math.Max(0, math.Min(100, 100-25*math.Log2(medianTTFT/300)))
		combinedScore := tpsScore*0.6 + ttftScore*0.4

		results = append(results, ModelPerformance{
			ModelName:     modelName,
			MedianTTFT:    medianTTFT,
			MedianTPS:     medianTPS,
			CombinedScore: combinedScore,
			SampleCount:   len(modelLogs),
		})
	}

	// Sort by combined score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].CombinedScore > results[j].CombinedScore
	})

	return results, nil
}

// parsePerformanceFromOther 从 other 字段解析性能和计算 TPS
func parsePerformanceFromOther(otherJSON string, completionTokens int) (frt float64, tps float64) {
	if otherJSON == "" {
		return 0, 0
	}

	// Parse the JSON to extract frt and output_duration_ms
	otherMap, err := common.StrToMap(otherJSON)
	if err != nil {
		return 0, 0
	}

	// Get frt (first response time in ms)
	if frtVal, ok := otherMap["frt"]; ok {
		switch v := frtVal.(type) {
		case float64:
			frt = v
		case int:
			frt = float64(v)
		}
	}

	// Get output_duration_ms and calculate TPS
	if outputDurationVal, ok := otherMap["output_duration_ms"]; ok {
		var outputDurationMs float64
		switch v := outputDurationVal.(type) {
		case float64:
			outputDurationMs = v
		case int:
			outputDurationMs = float64(v)
		}

		if outputDurationMs > 0 && completionTokens > 0 {
			// TPS = tokens / seconds = tokens * 1000 / ms
			tps = float64(completionTokens) * 1000 / outputDurationMs
		}
	}

	return frt, tps
}

// calculateMedian 计算中位数
func calculateMedian(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2
	}
	return values[n/2]
}

// GetDefaultPerformanceTimeRange 返回默认的性能查询时间范围（最近24小时）
func GetDefaultPerformanceTimeRange() (startTimestamp, endTimestamp int64) {
	now := time.Now()
	endTimestamp = now.Unix()
	startTimestamp = now.Add(-24 * time.Hour).Unix()
	return startTimestamp, endTimestamp
}

// parseTimeRangeFromString 解析时间范围字符串
func ParseTimeRangeFromString(timeStr string) (startTimestamp, endTimestamp int64) {
	switch strings.ToLower(timeStr) {
	case "today":
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return start.Unix(), now.Unix()
	case "week":
		return GetDefaultPerformanceTimeRange()
	default:
		return GetDefaultPerformanceTimeRange()
	}
}
