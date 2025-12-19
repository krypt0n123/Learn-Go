package main

import "fmt"

func main(){
	names:=[4]string{
		"john",
		"Bob",
		"Alice",
		"Ronyu",
	}
	fmt.Println(names)

	a:=names[0:3]
	b:=names[1:3]
	fmt.Println(a,b)

	b[0]="XXX"
	fmt.Println(a,b)
	fmt.Println(names)
}