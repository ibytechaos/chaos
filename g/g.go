/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: g.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-20
 * | Description: g.go
 * |-----------------------------------------------------------
 */

package g

import (
	"golang.org/x/sync/singleflight"
	"runtime"
)

var (
	CpuNumber      = runtime.NumCPU()
	ConfDir        = "./conf"
	SecurityDir    = ConfDir + "/security"
	DefaultBase    = 10
	DefaultBitSize = 64
	Call           = &singleflight.Group{}
)
