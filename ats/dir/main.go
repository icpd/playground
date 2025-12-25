package main

import (
	"fmt"
)

// ATSDirInfo 包含计算出的目录结构信息
type ATSDirInfo struct {
	DiskSizeBytes     int64 `json:"disk_size_bytes"`
	AvgObjSize        int64 `json:"avg_obj_size"`
	Segments          int64 `json:"segments"`
	BucketsPerSegment int64 `json:"buckets_per_segment"`
	TotalDirEntries   int64 `json:"total_dir_entries"`
}

// CalculateATSDirEntries 模拟 Traffic Server 的目录数量计算逻辑
// diskSizeBytes: 分配给卷的磁盘大小（字节）
// avgObjSize: proxy.config.cache.min_average_object_size 的值 (默认 8000)
func CalculateATSDirEntries(diskSizeBytes int64, avgObjSize int64) ATSDirInfo {
	if avgObjSize <= 0 {
		avgObjSize = 8000 // 防止除零，默认使用 ATS 标准值
	}

	const (
		DirDepth             = 4
		MaxBucketsPerSegment = 16384 // (1 << 16) / 4
	)

	// 1. 估算理论条目总数
	totalEntriesEst := diskSizeBytes / avgObjSize

	// 2. 计算总桶数
	totalBuckets := totalEntriesEst / DirDepth
	if totalBuckets == 0 {
		return ATSDirInfo{DiskSizeBytes: diskSizeBytes, AvgObjSize: avgObjSize}
	}

	// 3. 计算段数 (Segments)
	// 对应源码：(total_buckets + (MAX-1)) / MAX
	segments := (totalBuckets + MaxBucketsPerSegment - 1) / MaxBucketsPerSegment

	// 4. 计算每个段的桶数 (Buckets per segment)
	// 对应源码：(total_buckets + segments - 1) / segments
	bucketsPerSegment := (totalBuckets + segments - 1) / segments

	// 5. 最终总条目数
	finalTotal := segments * bucketsPerSegment * int64(DirDepth)

	return ATSDirInfo{
		DiskSizeBytes:     diskSizeBytes,
		AvgObjSize:        avgObjSize,
		Segments:          segments,
		BucketsPerSegment: bucketsPerSegment,
		TotalDirEntries:   finalTotal,
	}
}

func main() {
	var diskSize int64 = 20 * 1024 * 1024 * 1024
	avgSize := int64(8000)

	result := CalculateATSDirEntries(diskSize, avgSize)

	fmt.Printf("--- ATS 目录预估结果 ---\n")
	fmt.Printf("磁盘总大小: %d Bytes (%.2f GB)\n", result.DiskSizeBytes, float64(result.DiskSizeBytes)/1024/1024/1024)
	fmt.Printf("平均对象大小: %d Bytes\n", result.AvgObjSize)
	fmt.Printf("段数 (Segments): %d\n", result.Segments)
	fmt.Printf("每段桶数 (Buckets): %d\n", result.BucketsPerSegment)
	fmt.Printf("目录条目总数: %d\n", result.TotalDirEntries)
}
