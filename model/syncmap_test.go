package model

import (
	"testing"
	"time"
)

/**
 * @File: syncmap_test
 * @Author: Shijinyu
 * @Description:
 * @Date: 2020/7/16 18:46
 * @Version: 1.0.0
 */

func TestKvMap_Put(t *testing.T) {
	var kv KvMap
	kv.Put("1","first")
	kv.ListAll()
	time.Sleep(5*time.Second)
	kv.Put("1","5")
	kv.ListAll()
}
 
 