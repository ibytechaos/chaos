/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: zstd.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-20
 * | Description: zstd.go
 * |-----------------------------------------------------------
 */

package utils

import "github.com/klauspost/compress/zstd"

var (
	// ZstdCompress
	ZstdCompress, _ = zstd.NewWriter(nil)
	// ZstdDecompress
	ZstdDecompress, _ = zstd.NewReader(nil)
)

// Compress 压缩
func Compress(data []byte) []byte {
	return ZstdCompress.EncodeAll(data, make([]byte, 0, len(data)))
}

// Decompress 解压
func Decompress(data []byte) ([]byte, error) {
	return ZstdDecompress.DecodeAll(data, nil)
}
