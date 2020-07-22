package service

import (
	"coding.jd.com/etcd-service/conf"
	"coding.jd.com/etcd-service/model"
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/namespace"
	"go.etcd.io/etcd/etcdserver/etcdserverpb"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"runtime/debug"
	"time"
)

/**
 * @File: etcdinst
 * @Author: Shijinyu
 * @Description:
 * @Date: 2020/7/15 12:01
 * @Version: 1.0.0
 */

//type EtcdClient *clientv3.Client
//上述模式也可以使用成员函数

type EtcdInst struct{
	RequestTimeout int
	Namespace      string
	etcdclient     *clientv3.Client

}


func (ec *EtcdInst) Init(mes *conf.EtcdConf) error{
	cli,err :=clientv3.New(clientv3.Config{
		Endpoints:mes.Endpoints,
		DialTimeout:time.Duration(mes.DialTimeout)*time.Millisecond,
		Username:mes.Username,
		Password:mes.PassWord,
	})
	if err!=nil{
		return errors.New("create etcd client error:" + err.Error())
	}
	ec.Namespace = mes.Namespace
	ec.RequestTimeout = mes.RequestTimeout
	ns := ec.Namespace + "/"
	cli.KV = namespace.NewKV(cli.KV, ns)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, ns)
	cli.Lease = namespace.NewLease(cli.Lease, ns)
	ec.etcdclient = cli
	return nil
}

func (ec *EtcdInst) Close() error{
	if err:=ec.etcdclient.Close();err!=nil{
		return errors.New("close etcd client error:" + err.Error())
	}
	return nil
}

func (ec *EtcdInst) SyncAll() (*model.KvMap,error){
	var kvmap = &model.KvMap{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(ec.RequestTimeout))
	resp, err := ec.etcdclient.Get(ctx, "\x00", clientv3.WithFromKey())
	cancel()
	if err != nil {
		return nil,errors.New("sync all from etcd error:" + err.Error())
	}
	for _, kv := range resp.Kvs {
		//fmt.Printf("syncall:key = %v	value = %v",string(kv.Key),string(kv.Value))
		kvmap.Put(string(kv.Key),string(kv.Value))
	}
	return kvmap,nil
}

func (ec *EtcdInst) Put(key,value string ) error{
	if putRes, err := ec.etcdclient.Put(context.TODO(), key, value,clientv3.WithPrevKV()); err != nil {
		return errors.New("insert the kvs into the etcd failed!"+err.Error())
	} else {
		//返回前值
		//fmt.Println("revision: ", putRes.Header.Revision)
		if putRes.PrevKv != nil{
			fmt.Println(string(putRes.PrevKv.Value))
		}
		return nil
	}
}




func (ec *EtcdInst) GetWithFromKey() ([]string,error){
	var values []string
	if getRes, err := ec.etcdclient.Get(context.TODO(), "\x00",clientv3.WithFromKey()); err != nil {
		return values,errors.New("get key failed"+err.Error())
	} else {
		//输出get的结果集合
		//fmt.Println(fmt.Println(getRes.Kvs))
		//输出value值，如果返回kv对只有一对，Kvs[1]将返回超出阈值错误。
		for _,kvs:= range getRes.Kvs{
			values = append(values, string(kvs.Value))
		}
		return values,nil
	}
}

func (ec *EtcdInst) GetWithPrefix(prefix string) ([]string,[]string,error){
	var values []string
	var keys []string
	if getRes, err := ec.etcdclient.Get(context.TODO(), prefix,clientv3.WithPrefix()); err != nil {
		return keys,values,errors.New("get key failed"+err.Error())
	} else {
		//输出get的结果集合
		//fmt.Println(fmt.Println(getRes.Kvs))
		//输出value值，如果返回kv对只有一对，Kvs[1]将返回超出阈值错误。
		for _,kvs:= range getRes.Kvs{
			values = append(values, string(kvs.Value))
			keys = append(keys,string(kvs.Key))
		}
		return keys,values,nil
	}
}

func (ec *EtcdInst) Get(key string) ([]*mvccpb.KeyValue,error){
	getRes, err := ec.etcdclient.Get(context.TODO(),key)
	if err != nil {
		return nil,errors.New("get key failed"+err.Error())
	}
	//输出get的结果集合
	//fmt.Println(fmt.Println(getRes.Kvs))
	//输出value值，如果返回kv对只有一对，Kvs[1]将返回超出阈值错误。
	if len(getRes.Kvs)<=0{
		return nil,errors.New(fmt.Sprintf("there is no key=%v in etcd",key))
	}
	return getRes.Kvs,nil
}

func (ec *EtcdInst) Delete(key string) error{
	if delRes, err := ec.etcdclient.Delete(context.TODO(), key, clientv3.WithPrevKV());err!=nil{
		return err
	} else {
		//what is the difference of the delRes.Delete and delRes.PrevKvs
		//deleted is the number of keys deleted by the delete range request.
		fmt.Println(delRes.Deleted)
		//if prev_kv is set in the request, the previous key-value pairs will be returned.
		if delRes.Deleted>0{
			for _,kv:=range delRes.PrevKvs{
				if kv!=nil{
					fmt.Println(string(kv.Value))
				}
			}
			fmt.Println(string(delRes.PrevKvs[0].Value))
		}
		return nil
	}
}

func (ec *EtcdInst) DeleteAll(prefix string) error{
	if delRes, err := ec.etcdclient.KV.Delete(context.TODO(), prefix, clientv3.WithPrefix());err!=nil{
		return err
	} else {
		fmt.Println(delRes.Deleted)
		return nil
	}

}

func (ec *EtcdInst) Listen(key string){
	etcdwatchs := ec.etcdclient.Watch(context.Background(),key,clientv3.WithFromKey())
	go func(){
		defer func(){
			if r:=recover();r!=nil{
				fmt.Println("unexpected panic when sync:", r, "\nstack=", string(debug.Stack()))
			}
		}()
		for etcdwatch := range etcdwatchs {
			for _,ee := range etcdwatch.Events{
				fmt.Println("changed key: "+string(ee.Kv.Key))
				fmt.Println("changed value: "+string(ee.Kv.Value))
			}
		}
	}()
}

func (ec *EtcdInst) ShowMembers()([]*etcdserverpb.Member,error){
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(ec.RequestTimeout))
	resp, err := ec.etcdclient.MemberList(ctx)
	cancel()
	if err != nil{
		return nil,err
	}
	return resp.Members,nil
}

 
 