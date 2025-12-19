package main

import (
	"fmt"
)

func do(i interface{}){
	switch v:=i.(type){
	case int:
		fmt.Printf("二倍的%v 是%v\n",v,v*2)
	case string:
		fmt.Printf("%q长度为%v字节\n",v,len(v))
	case bool:
		fmt.Printf("我不知道类型%T\n",v)
	}
}

func main(){
	do(21)
	do("hello")
	do(false)
}