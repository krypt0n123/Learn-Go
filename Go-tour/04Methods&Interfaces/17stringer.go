package main

import "fmt"

type Person struct{
	Name string
	Age int
}

func (p Person) String() string{
	return fmt.Sprintf("%v(%v years)",p.Name,p.Age)
}

func main(){
	a:=Person{"Kiko",20}
	b:=Person{"Lazy",18}
	fmt.Println(a,b)
}