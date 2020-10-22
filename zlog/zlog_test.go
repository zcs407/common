package zlog

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	err := InitLog("aaa", "aaa", 0, 0)
	if err != nil {
		fmt.Println("test log init failed", err)
	}
	Err(errors.New("测试"), "test")
}
