package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"math/rand"
	"os"
	"time"
)

const (
	AccessKeyId="LTAI4G6k8fYijwXpH7ERX67j"
	AccessKeySecret="P16f4Y9BxGW5wuT6XmalQGrnKZspjY"
	EndPoint="oss-cn-shenzhen.aliyuncs.com"
	Bucket="fqhchattest"
)

func UpLoadOss(){
	client,err:=oss.New(EndPoint,AccessKeyId,AccessKeySecret)
	if err!=nil{
		fmt.Println("oss客户端对象错误")
		fmt.Println(err.Error())
		return
	}
	//todo 获得bucket
	bucket,err := client.Bucket(Bucket)
	if err!=nil{
		fmt.Println("bucket对象错误")
		fmt.Println(err.Error())
		return
	}
	//todo 设置文件名称
	//time.Now().Unix()
	//suffix := ".png"
	filename := fmt.Sprintf("mnt/%d%04d.png",
		time.Now().Unix(), rand.Int31())
	//打开文件
	f, err := os.OpenFile("./address.png", os.O_CREATE|os.O_APPEND, 6)

	if err!=nil{
		fmt.Println("文件io对象错误错误")
		print(err.Error())
		return
	}
	//todo 通过bucket上传
	//err=bucket.PutObject(filename,f)
	err=bucket.PutObject(filename,f)
	if err!=nil{
		fmt.Println("bucket上传错误错误")
		print(err.Error())
		return
	}
	//todo 获得url地址
	url := "http://"+Bucket+"."+EndPoint+"/"+filename
	print(url)
}
func main() {
	//UpLoadOss()
	fmt.Println(rand.Int31())
	fmt.Println(rand.Int31())
	fmt.Println(rand.Int31())
	fmt.Println(rand.Int31())

}
