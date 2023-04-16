/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2023 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: ptr.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2023-04-16
 * | Description: ptr.go
 * |-----------------------------------------------------------
 */

package utils

// StringPtr 字符串指针
func StringPtr(s string) *string {
	return &s
}

// IntPtr 整型指针
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr 整型指针
func Int64Ptr(i int64) *int64 {
	return &i
}

// BoolPtr 布尔指针
func BoolPtr(b bool) *bool {
	return &b
}

// Float32Ptr 浮点型指针
func Float32Ptr(f float32) *float32 {
	return &f
}

// Float64Ptr 浮点型指针
func Float64Ptr(f float64) *float64 {
	return &f
}
