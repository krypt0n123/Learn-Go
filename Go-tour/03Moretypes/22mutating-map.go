package main

import "fmt"

func main(){
	m:=make(map[string]int)

	m["ans"]=42
	fmt.Println("num:",m["ans"])

	m["ans"]=48
	fmt.Println("num:",m["ans"])

	delete(m,"ans")
	fmt.Println("num:",m["ans"])

	v,ok:=m["ans"]
	fmt.Println("num:",v,"have?",ok)
}