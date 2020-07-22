package test

import (
	"coding.jd.com/etcd-service/conf"
	"coding.jd.com/etcd-service/service"
	"flag"
	"fmt"
	"log"
	"testing"
)

/**
 * @File: etcdinst_test
 * @Author: Shijinyu
 * @Description:
 * @Date: 2020/7/15 12:15
 * @Version: 1.0.0
 */

func TestEtcd(t *testing.T){
	configPath := flag.String("f", "D:\\goproject\\src\\coding.jd.com\\etcd-service\\conf\\config.json", "config file")
	flag.Parse()
	var ec = &service.EtcdInst{}
	etcdmes,err :=conf.InitConfig(*configPath);
	if err != nil {
		log.Panic("parse config failed!")
	}
	ec.Init(&etcdmes)
	defer ec.Close()

	//ec.Listen("test1")
	//time.Sleep(5+time.Second)
	//if err := ec.Put("test1/testlisten2","success3");err!=nil{
	//	fmt.Println("insert error"+err.Error())
	//} else{
	//	fmt.Println("insert success!")
	//}

	//test the get function
	if kvs,err := ec.Get("writer/1");err!=nil{
		fmt.Println("get value failed: "+err.Error())
	}else{
		fmt.Println(string(kvs[0].Key))
		fmt.Println(string(kvs[0].Value))
		fmt.Println(kvs[0].CreateRevision)
		fmt.Println(kvs[0].ModRevision)
		fmt.Println(kvs[0].Version)
		fmt.Println(kvs[0].Lease)
	}

	//test the getwithprefix function
	//if value,err := ec.GetWithPrefix("sys");err!=nil{
	//	fmt.Println("get value failed: "+err.Error())
	//}else{
	//	fmt.Println(value)
	//}

	//test the getwithfromkey function
	//if value,err := ec.GetWithFromKey();err!=nil{
	//	fmt.Println("get value failed: "+err.Error())
	//}else{
	//	fmt.Println(value)
	//}



	//test the deletd function
	//if err:= ec.Delete("test1");err!=nil{
	//	fmt.Println("delete failed"+err.Error())
	//} else {
	//	fmt.Println("delete successful!")
	//}

	//test the delete all function
	//if err:= ec.DeleteAll("test");err!=nil{
	//	fmt.Println("delete failed: "+err.Error())
	//}else{
	//	fmt.Println("delete successful!")
	//}

	//test sybcall function
	//if kvm,err:=ec.SyncAll();err!=nil{
	//	fmt.Println("sync failed: "+err.Error())
	//}else{
	//	mp := kvm.ListAll()
	//	fmt.Println(* mp)
	//}

}
 