package main

import (
	"fmt"
	"math"
)

func pow(x,n,lim float64)float64{
	//if后可以在条件表达式前执行一个语句，只在if语句(包括对应的else)内起作用
	if v:=math.Pow(x,n); v<lim{
		return v
	}
	return v
}
func main(){
	fmt.Println(
		pow(3,2,10),
		pow(3,3,20),
	)
}