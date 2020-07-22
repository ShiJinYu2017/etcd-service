package model

import (
	"errors"
	"fmt"
	"sync"
)

/**
 * @File: kvmodel
 * @Author: Shijinyu
 * @Description:
 * @Date: 2020/7/15 12:14
 * @Version: 1.0.0
 */
type KvMap struct {
	kvmap	sync.Map
}


func (km *KvMap) Put(key,value string) {
	km.kvmap.Store(key,value)
}

func (km *KvMap) Get(key string) (string,error){
	value,ok := km.kvmap.Load(key)
	if ok{
		return value.(string),nil
	} else {
		return "",errors.New("get value failed")
	}
}

func (km *KvMap) Remove(key string){
	km.kvmap.Delete(key)
}

func (km *KvMap) Length() int{
	len := 0
	km.kvmap.Range(func(key, value interface{}) bool{
		len++
		return true
	})
	return len
}

func (km *KvMap) ListAll() *map[string]string{
	var kvs map[string]string
	kvs = make(map[string]string)
	km.kvmap.Range(func(key, value interface{}) bool {
		kvs[key.(string)] = value.(string)
		fmt.Printf("key = %v	value = %v\n",key.(string),value.(string))
		return true
	})
	return &kvs
}

 