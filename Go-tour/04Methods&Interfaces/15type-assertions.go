package main

import "fmt"

func main(){
	var i interface{}="hello"

	s:=i.(string)
	fmt.Println(s)

	s,o:=i.(string)
	fmt.Println(s,o)

	f,ok:=i.(float64)
	fmt.Println(f,ok)

	f=i.(float64)
	fmt.Println(f)
}