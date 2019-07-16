package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)
var(
	_name *string = flag.String("f","data.txt"," file name: 默认data.txt")
	_sep *string = flag.String("s"," ","分隔符：默认空格")

)
func main(){
	flag.Parse()

	p,err:=os.OpenFile(*_name,os.O_RDWR,os.ModePerm)
	if err != nil{
		fmt.Println(err)
		return
	}
	src,err1:=ioutil.ReadAll(p)
	if err1 != nil{
		fmt.Println(err)
		return
	}
	//实现16进制的数据 显示以string显示
	tmp:=strings.Split(string(src),*_sep)
	var source []byte
	for _,v:=range tmp{
		t:=""
		t+="0x"+v
		i,_:=strconv.ParseInt(v,16,32)
		source=append(source,byte(i))
	}
	fmt.Println(string(source))
}
