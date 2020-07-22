package test

import (
	"coding.jd.com/etcd-service/conf"
	"coding.jd.com/etcd-service/service"
	"flag"
	"log"
	"testing"
)

/**
 * @File: etcdwatcher_test
 * @Author: Shijinyu
 * @Description:
 * @Date: 2020/7/17 15:02
 * @Version: 1.0.0
 */


func Test_Watcher(t *testing.T) {
	configPath := flag.String("f", "D:\\goproject\\src\\coding.jd.com\\etcd-service\\conf\\config.json", "config file")
	flag.Parse()
	var ec = &service.EtcdInst{}
	etcdmes,err :=conf.InitConfig(*configPath);
	if err != nil {
		log.Panic("parse config failed!")
	}
	ec.Init(&etcdmes)
	defer ec.Close()

	ec.Listen("test1")
	//阻塞主线程，禁止退出
	select {}


}
 