package utils

import (
	"fmt"
	"reflect"
)

type byteSize int

const (
	byteSizeB byteSize = 1 << (10 * iota)
	byteSizeKB
	byteSizeMB
	byteSizeGB
	byteSizeTB
	byteSizePB
)

func ByteSizeFormat(size float64) string {
	switch {
	case size >= float64(byteSizePB):
		return fmt.Sprintf("%.2fPB", size/float64(byteSizePB))
	case size >= float64(byteSizeTB):
		return fmt.Sprintf("%.2fTB", size/float64(byteSizeTB))
	case size >= float64(byteSizeGB):
		return fmt.Sprintf("%.2fGB", size/float64(byteSizeGB))
	case size >= float64(byteSizeMB):
		return fmt.Sprintf("%.2fMB", size/float64(byteSizeMB))
	case size >= float64(byteSizeKB):
		return fmt.Sprintf("%.2fKB", size/float64(byteSizeKB))
	case size >= float64(byteSizeB):
		return fmt.Sprintf("%.2fB", size/float64(byteSizeB))
	}
	return fmt.Sprintf("%.2fB", size/float64(byteSizeB))
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
