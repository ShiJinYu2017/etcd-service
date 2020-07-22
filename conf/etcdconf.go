package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/**
 * @File: etcdconf
 * @Author: Shijinyu
 * @Description:
 * @Date: 2020/7/15 11:37
 * @Version: 1.0.0
 */

type EtcdConf struct {
	Endpoints      []string 	`json:"Endpoints"`
	DialTimeout    int			`json:"DialTimeout"`
	RequestTimeout int			`json:"RequestTimeout"`
	Namespace      string		`json:"Namespace"`
	Username	   string		`json:"Username"`
	PassWord       string       `json:"Password"`
}
var etcd EtcdConf
func InitConfig(configPath string) (EtcdConf,error){
	var err error
	configBytes, err := ioutil.ReadFile(configPath)
	fmt.Println(string(configBytes))
	if err != nil {
		return etcd,err
	}
	if err = json.Unmarshal(configBytes, &etcd); err != nil {
		fmt.Println(err.Error())
		return etcd,err
	}
	fmt.Println(etcd)
	return etcd,nil
}
 
 